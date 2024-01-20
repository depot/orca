package layers

// This is a copy of the continuity BuildManifest function.
// This copy creates layer entries rather than building a manifest as we
// need the timestamps for each file.
//
// Additionally, we made this copy to copy the hardlink handling of continuity.

import (
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/containerd/continuity"
	layerv1 "github.com/depot/orca/proto/depot/orca/layer/v1"
	"github.com/depot/orca/util/chunk"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Walk creates layer entries for the given context.
func Walk(ctx continuity.Context) (*layerv1.LayerEntries, error) {
	resourcesByPath := map[string]*modTimeResource{}
	hardLinks := newHardlinkManager()

	if err := ctx.Walk(func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking %s: %w", p, err)
		}

		if p == string(os.PathSeparator) {
			// skip root
			return nil
		}

		inner, err := ctx.Resource(p, fi)
		if err != nil {
			if err == continuity.ErrNotFound {
				return nil
			}
			return fmt.Errorf("failed to get resource %q: %w", p, err)
		}
		resource := &modTimeResource{inner: inner, modTime: fi.ModTime()}

		if _, ok := resource.inner.(continuity.Hardlinkable); ok {
			err := hardLinks.Add(fi, resource)
			if err == nil {
				// Resource has been accepted by hardlink manager so we don't add
				// it to the resourcesByPath until we merge at the end.
				return nil
			} else if err != errNotAHardLink {
				// handle any other case where we have a proper error.
				return fmt.Errorf("adding hardlink %s: %w", p, err)
			}
		}

		resourcesByPath[p] = resource

		return nil
	}); err != nil {
		return nil, err
	}

	// merge and post-process the hardlinks.
	hardLinked, err := hardLinks.Merge()
	if err != nil {
		return nil, err
	}

	for _, resource := range hardLinked {
		resourcesByPath[resource.inner.Path()] = resource
	}

	var resources []*modTimeResource
	for _, resource := range resourcesByPath {
		resources = append(resources, resource)
	}

	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].inner.Path() < resources[j].inner.Path()
	})

	entries := make([]*layerv1.LayerEntry, 0, len(resources))
	for _, resource := range resources {
		entry, err := newLayerEntry(resource)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &layerv1.LayerEntries{Entries: entries}, nil
}

// Returns an error when the digest is not made up of chunks.
func newLayerEntry(resource *modTimeResource) (*layerv1.LayerEntry, error) {
	entry := &layerv1.LayerEntry{
		Path:  []string{resource.inner.Path()},
		Mode:  uint32(resource.inner.Mode()),
		Uid:   resource.inner.UID(),
		Gid:   resource.inner.GID(),
		Mtime: timestamppb.New(resource.modTime),
	}

	if xattrer, ok := resource.inner.(continuity.XAttrer); ok {
		// Sorts the XAttrs by name for consistent ordering.
		keys := []string{}
		xattrs := xattrer.XAttrs()
		for k := range xattrs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			entry.Xattr = append(entry.Xattr, &layerv1.XAttr{Name: k, Data: xattrs[k]})
		}
	}

	switch r := resource.inner.(type) {
	case continuity.RegularFile:
		entry.Path = r.Paths()
		entry.SizeBytes = uint64(r.Size())

		digests := r.Digests()
		if len(digests) > 0 {
			chunks, err := chunk.NewChunksFromString(string(digests[0]))
			if err != nil {
				return nil, err
			}
			entry.Blocks = make([]*layerv1.Block, len(chunks))
			for i, c := range chunks {
				dgst := &layerv1.Digest{
					Algorithm: layerv1.Digest_ALGORITHM_XXH64,
					Sum:       c.Hash,
				}
				entry.Blocks[i] = &layerv1.Block{
					Digest:    dgst,
					SizeBytes: uint64(c.Size),
				}
			}
		}
	case continuity.SymLink:
		entry.Target = r.Target()
	case continuity.Device:
		entry.Major, entry.Minor = r.Major(), r.Minor()
		entry.Path = r.Paths()
	case continuity.NamedPipe:
		entry.Path = r.Paths()
	}

	// enforce a few stability guarantees that may not be provided by the
	// resource implementation.
	sort.Strings(entry.Path)

	return entry, nil
}

// Continuity does not provide a way to get the mod time of a resource.
// We are adding this extra interface to allow us to get the mod time of a resource.
type modTimeResource struct {
	inner   continuity.Resource
	modTime time.Time
}

var errNotAHardLink = fmt.Errorf("invalid hardlink")

type hardlinkManager struct {
	hardlinks map[hardlinkKey][]*modTimeResource
}

func newHardlinkManager() *hardlinkManager {
	return &hardlinkManager{
		hardlinks: map[hardlinkKey][]*modTimeResource{},
	}
}

// Add attempts to add the resource to the hardlink manager. If the resource
// cannot be considered as a hardlink candidate, errNotAHardLink is returned.
func (hlm *hardlinkManager) Add(fi os.FileInfo, resource *modTimeResource) error {
	if _, ok := resource.inner.(continuity.Hardlinkable); !ok {
		return errNotAHardLink
	}

	key, err := newHardlinkKey(fi)
	if err != nil {
		return err
	}

	hlm.hardlinks[key] = append(hlm.hardlinks[key], resource)

	return nil
}

// Merge processes the current state of the hardlink manager and merges any
// shared nodes into hard linked resources.
func (hlm *hardlinkManager) Merge() ([]*modTimeResource, error) {
	var resources []*modTimeResource
	for key, linked := range hlm.hardlinks {
		if len(linked) < 1 {
			return nil, fmt.Errorf("no hardlink entries for dev, inode pair: %#v", key)
		}

		innerResources := make([]continuity.Resource, len(linked))
		for i, linkedResource := range linked {
			innerResources[i] = linkedResource.inner
		}
		merged, err := continuity.Merge(innerResources...)
		if err != nil {
			return nil, fmt.Errorf("error merging hardlink: %w", err)
		}

		// If the merge is successful, then the mod time of the first resource
		// should be the mod time of all the resources.
		first := linked[0]
		mergedModTime := &modTimeResource{
			inner:   merged,
			modTime: first.modTime,
		}

		resources = append(resources, mergedModTime)
	}

	return resources, nil
}

// hardlinkKey provides a tuple-key for managing hardlinks. This is system-
// specific.
type hardlinkKey struct {
	dev   uint64
	inode uint64
}

// newHardlinkKey returns a hardlink key for the provided file info. If the
// resource does not represent a possible hardlink, errNotAHardLink will be
// returned.
func newHardlinkKey(fi os.FileInfo) (hardlinkKey, error) {
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return hardlinkKey{}, fmt.Errorf("cannot resolve (*syscall.Stat_t) from os.FileInfo")
	}

	if sys.Nlink < 2 {
		// NOTE(stevvooe): This is not always true for all filesystems. We
		// should somehow detect this and provided a slow "polyfill" that
		// leverages os.SameFile if we detect a filesystem where link counts
		// is not really supported.
		return hardlinkKey{}, errNotAHardLink
	}

	//nolint:unconvert
	return hardlinkKey{dev: uint64(sys.Dev), inode: uint64(sys.Ino)}, nil
}
