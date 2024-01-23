package buildspec

import (
	"context"
	"strings"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/contrib/seccomp"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// TODO:
// get the uid/gid
// WithUIDGID
// make the cwd
// mkdirall/ chown the cwd
// setup optional networking
// GenerateSpec
// Adds another host to the /etc/hosts file.
//  WithTTY

type Options struct {
	Platform *platforms.Platform

	RootlessRunC bool

	Args []string
	Cwd  string

	RootFSPath     string
	ReadonlyRootFS bool
	Userstr        string // NOTE: this option requires a parent container as it reads the parent's /etc/passwd and /etc/group.

	EtcHostsPath    string
	AdditionalHosts []string // NOTE: this option will create a new snapshot and appends the hosts to the containerd host's /etc/hosts. The default is to use the containerd host's /etc/hosts.
	ResolvConfPath  string

	ProxyEnv     *ProxyEnv           // Adds to the env vars of the container. Fascinating history about 'em: https://about.gitlab.com/blog/2021/01/27/we-need-to-talk-no-proxy/
	Secrets      map[string]string   // Adds to the env vars of the container.
	EnvVars      map[string]string   // Adds to the env vars of the container.
	NoNetworking bool                // If true, then the container will not have any networking.
	RLimits      []specs.POSIXRlimit // If empty then the container will have no limits.
}

type ProxyEnv struct {
	HTTPProxy  string
	HTTPSProxy string
	FTPProxy   string
	NoProxy    string
	AllProxy   string
}

func WithBuildSpec(opts *Options) oci.SpecOpts {
	if opts.Platform == nil {
		platform := platforms.DefaultSpec()
		opts.Platform = &platform
	}

	var specOpts []oci.SpecOpts

	specOpts = append(specOpts, oci.WithRootFSPath(opts.RootFSPath))

	specOpts = append(specOpts, oci.WithProcessArgs(opts.Args...))
	specOpts = append(specOpts, oci.WithProcessCwd(opts.Cwd)) // TODO: make the cwd and chown it.
	specOpts = append(specOpts, oci.WithDefaultPathEnv)

	if opts.Userstr != "" {
		specOpts = append(specOpts, oci.WithUser(opts.Userstr))
		specOpts = append(specOpts, oci.WithAdditionalGIDs(opts.Userstr))
		// TODO: Also, I need to make the cwd and chown it.
	}
	specOpts = append(specOpts, EnsureAdditionalGids())

	if opts.ReadonlyRootFS {
		specOpts = append(specOpts, oci.WithRootFSReadonly())
	}

	if len(opts.AdditionalHosts) == 0 {
		specOpts = append(specOpts, oci.WithHostHostsFile)
	} else {
		// TODO: make a snapshot and append the hosts to the snapshot's /etc/hosts.
	}

	if opts.EtcHostsPath != "" {
		specOpts = append(specOpts, EtcHosts(opts.EtcHostsPath))
	}

	if opts.ResolvConfPath != "" {
		specOpts = append(specOpts, EtcResolvConf(opts.ResolvConfPath))
	}

	if opts.ProxyEnv != nil {
		specOpts = append(specOpts, SetProxyEnv(opts.ProxyEnv))
	}

	specOpts = append(specOpts, SetEnvVars(opts.Secrets))
	specOpts = append(specOpts, SetEnvVars(opts.EnvVars))

	if opts.Platform != nil && opts.Platform.OS == "linux" {
		specOpts = append(specOpts, LinuxHostNetworking)
	}

	if opts.RootlessRunC {
		specOpts = append(specOpts, WithoutRoot)
	}

	specOpts = append(specOpts, Limit(opts.RLimits))
	specOpts = append(specOpts, seccomp.WithDefaultProfile())

	return oci.Compose(specOpts...)
}

func SetProxyEnv(env *ProxyEnv) oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		vars := []string{}
		if env.HTTPProxy != "" {
			vars = append(vars, "HTTP_PROXY="+env.HTTPProxy, "http_proxy="+env.HTTPProxy)
		}
		if env.HTTPSProxy != "" {
			vars = append(vars, "HTTPS_PROXY="+env.HTTPSProxy, "https_proxy="+env.HTTPSProxy)
		}
		if env.FTPProxy != "" {
			vars = append(vars, "FTP_PROXY="+env.FTPProxy, "ftp_proxy="+env.FTPProxy)
		}
		if env.NoProxy != "" {
			vars = append(vars, "NO_PROXY="+env.NoProxy, "no_proxy="+env.NoProxy)
		}
		if env.AllProxy != "" {
			vars = append(vars, "ALL_PROXY="+env.AllProxy, "all_proxy="+env.AllProxy)
		}

		return oci.WithEnv(vars)(ctx, client, container, spec)
	}
}

func SetEnvVars(envvars map[string]string) oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		vars := []string{}
		for key, value := range envvars {
			vars = append(vars, key+"="+value)
		}
		return oci.WithEnv(vars)(ctx, client, container, spec)
	}
}

// WithRootless converts spec to be compatible with "rootless" runc.
// * Remove /sys mount
// * Remove cgroups
func WithoutRoot(ctx context.Context, client oci.Client, c *containers.Container, spec *specs.Spec) error {
	// Remove /sys mount because we can't mount /sys when the daemon netns
	// is not unshared from the host.
	//
	// Instead, we could bind-mount /sys from the host, however, `rbind, ro`
	// does not make /sys/fs/cgroup read-only (and we can't bind-mount /sys
	// without rbind)
	//
	// PR for making /sys/fs/cgroup read-only is proposed, but it is very
	// complicated: https://github.com/opencontainers/runc/pull/1869
	//
	// For buildkit usecase, we suppose we don't need to provide /sys to
	// containers and remove /sys mount as a workaround.
	var mounts []specs.Mount
	for _, mount := range spec.Mounts {
		if strings.HasPrefix(mount.Destination, "/sys") {
			continue
		}
		mounts = append(mounts, mount)
	}
	spec.Mounts = mounts

	if spec.Linux != nil {
		spec.Linux.Resources = nil
		spec.Linux.CgroupsPath = ""
	}
	return nil
}

func LinuxHostNetworking(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
	return oci.WithHostNamespace(specs.NetworkNamespace)(ctx, client, container, spec)
}

func Limit(limits []specs.POSIXRlimit) oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		if len(limits) == 0 {
			spec.Process.Rlimits = nil
		} else {
			spec.Process.Rlimits = limits
		}
		return nil
	}
}

// EnsureAdditionalGids ensures that the primary GID is also included in the additional GID list.
// From https://github.com/containerd/containerd/blob/v1.7.0-beta.4/oci/spec_opts.go#L124-L133
func EnsureAdditionalGids() oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		if spec.Process == nil {
			spec.Process = &specs.Process{}
		}
		for _, f := range spec.Process.User.AdditionalGids {
			if f == spec.Process.User.GID {
				return nil
			}
		}
		spec.Process.User.AdditionalGids = append([]uint32{spec.Process.User.GID}, spec.Process.User.AdditionalGids...)
		return nil
	}
}

func EtcHosts(path string) oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/etc/hosts",
			Source:      path,
			Type:        "bind",
			Options:     []string{"nosuid", "noexec", "nodev", "rbind", "ro", "noatime"},
		})
		return nil
	}
}

func EtcResolvConf(path string) oci.SpecOpts {
	return func(ctx context.Context, client oci.Client, container *containers.Container, spec *oci.Spec) error {
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/etc/resolv.conf",
			Source:      path,
			Type:        "bind",
			Options:     []string{"nosuid", "noexec", "nodev", "rbind", "ro", "noatime"},
		})
		return nil
	}
}
