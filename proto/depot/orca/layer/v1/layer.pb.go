// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: depot/orca/layer/v1/layer.proto

package layerv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Digest_Algorithm int32

const (
	Digest_ALGORITHM_UNSPECIFIED Digest_Algorithm = 0
	Digest_ALGORITHM_XXH64       Digest_Algorithm = 1
)

// Enum value maps for Digest_Algorithm.
var (
	Digest_Algorithm_name = map[int32]string{
		0: "ALGORITHM_UNSPECIFIED",
		1: "ALGORITHM_XXH64",
	}
	Digest_Algorithm_value = map[string]int32{
		"ALGORITHM_UNSPECIFIED": 0,
		"ALGORITHM_XXH64":       1,
	}
)

func (x Digest_Algorithm) Enum() *Digest_Algorithm {
	p := new(Digest_Algorithm)
	*p = x
	return p
}

func (x Digest_Algorithm) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Digest_Algorithm) Descriptor() protoreflect.EnumDescriptor {
	return file_depot_orca_layer_v1_layer_proto_enumTypes[0].Descriptor()
}

func (Digest_Algorithm) Type() protoreflect.EnumType {
	return &file_depot_orca_layer_v1_layer_proto_enumTypes[0]
}

func (x Digest_Algorithm) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Digest_Algorithm.Descriptor instead.
func (Digest_Algorithm) EnumDescriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{3, 0}
}

// LayerEntries
type LayerEntries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*LayerEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *LayerEntries) Reset() {
	*x = LayerEntries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LayerEntries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LayerEntries) ProtoMessage() {}

func (x *LayerEntries) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LayerEntries.ProtoReflect.Descriptor instead.
func (*LayerEntries) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{0}
}

func (x *LayerEntries) GetEntries() []*LayerEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type LayerEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Path specifies the path from the bundle root. If more than one
	// path is present, the entry may represent a hardlink, rather than using
	// a link target. The path format is operating system specific.
	Path []string `protobuf:"bytes,1,rep,name=path,proto3" json:"path,omitempty"`
	// Size specifies the size in bytes.
	SizeBytes uint64 `protobuf:"varint,2,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	// Ordered set of blocks that make up the content of the entry.
	Blocks []*Block `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty"`
	// Uid specifies the user id for the entry.
	Uid int64 `protobuf:"varint,4,opt,name=uid,proto3" json:"uid,omitempty"`
	// Gid specifies the group id for the entry.
	Gid int64 `protobuf:"varint,5,opt,name=gid,proto3" json:"gid,omitempty"`
	// Modified time of the entry.
	Mtime *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=mtime,proto3" json:"mtime,omitempty"`
	// Mode defines the file mode and permissions. We've used the same
	// bit-packing from Go's os package,
	// http://golang.org/pkg/os/#FileMode, since they've done the work of
	// creating a cross-platform layout.
	Mode uint32 `protobuf:"varint,7,opt,name=mode,proto3" json:"mode,omitempty"`
	// Target defines the target of a hard or soft link. Absolute links start
	// with a slash and specify the entry relative to the bundle root.
	// Relative links do not start with a slash and are relative to the
	// entry path.
	Target string `protobuf:"bytes,8,opt,name=target,proto3" json:"target,omitempty"`
	// Major specifies the major device number for character and block devices.
	Major uint64 `protobuf:"varint,9,opt,name=major,proto3" json:"major,omitempty"`
	// Minor specifies the minor device number for character and block devices.
	Minor uint64 `protobuf:"varint,10,opt,name=minor,proto3" json:"minor,omitempty"`
	// Xattr provides storage for extended attributes for the target entry.
	Xattr []*XAttr `protobuf:"bytes,11,rep,name=xattr,proto3" json:"xattr,omitempty"`
	// Ads stores one or more alternate data streams for the target entry.
	Ads []*ADSEntry `protobuf:"bytes,12,rep,name=ads,proto3" json:"ads,omitempty"`
}

func (x *LayerEntry) Reset() {
	*x = LayerEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LayerEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LayerEntry) ProtoMessage() {}

func (x *LayerEntry) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LayerEntry.ProtoReflect.Descriptor instead.
func (*LayerEntry) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{1}
}

func (x *LayerEntry) GetPath() []string {
	if x != nil {
		return x.Path
	}
	return nil
}

func (x *LayerEntry) GetSizeBytes() uint64 {
	if x != nil {
		return x.SizeBytes
	}
	return 0
}

func (x *LayerEntry) GetBlocks() []*Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

func (x *LayerEntry) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *LayerEntry) GetGid() int64 {
	if x != nil {
		return x.Gid
	}
	return 0
}

func (x *LayerEntry) GetMtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Mtime
	}
	return nil
}

func (x *LayerEntry) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

func (x *LayerEntry) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *LayerEntry) GetMajor() uint64 {
	if x != nil {
		return x.Major
	}
	return 0
}

func (x *LayerEntry) GetMinor() uint64 {
	if x != nil {
		return x.Minor
	}
	return 0
}

func (x *LayerEntry) GetXattr() []*XAttr {
	if x != nil {
		return x.Xattr
	}
	return nil
}

func (x *LayerEntry) GetAds() []*ADSEntry {
	if x != nil {
		return x.Ads
	}
	return nil
}

// Block is a content-addressable variable-sized block of file content.
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SizeBytes uint64  `protobuf:"varint,1,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	Digest    *Digest `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{2}
}

func (x *Block) GetSizeBytes() uint64 {
	if x != nil {
		return x.SizeBytes
	}
	return 0
}

func (x *Block) GetDigest() *Digest {
	if x != nil {
		return x.Digest
	}
	return nil
}

type Digest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm Digest_Algorithm `protobuf:"varint,1,opt,name=algorithm,proto3,enum=depot.orca.layer.v1.Digest_Algorithm" json:"algorithm,omitempty"`
	Sum       uint64           `protobuf:"varint,2,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *Digest) Reset() {
	*x = Digest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Digest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Digest) ProtoMessage() {}

func (x *Digest) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Digest.ProtoReflect.Descriptor instead.
func (*Digest) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{3}
}

func (x *Digest) GetAlgorithm() Digest_Algorithm {
	if x != nil {
		return x.Algorithm
	}
	return Digest_ALGORITHM_UNSPECIFIED
}

func (x *Digest) GetSum() uint64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

// XAttr encodes extended attributes for a entry.
type XAttr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name specifies the attribute name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Data specifies the associated data for the attribute.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *XAttr) Reset() {
	*x = XAttr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XAttr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XAttr) ProtoMessage() {}

func (x *XAttr) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XAttr.ProtoReflect.Descriptor instead.
func (*XAttr) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{4}
}

func (x *XAttr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *XAttr) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// ADSEntry encodes information for a Windows Alternate Data Stream.
type ADSEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name specifices the stream name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Data specifies the stream data.
	// See also the description about the digest below.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Digest is a CAS representation of the stream data.
	//
	// At least one of data or digest MUST be specified, and either one of them
	// SHOULD be specified.
	//
	// How to access the actual data using the digest is implementation-specific,
	// and implementations can choose not to implement digest.
	// So, digest SHOULD be used only when the stream data is large.
	Digest string `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *ADSEntry) Reset() {
	*x = ADSEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ADSEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ADSEntry) ProtoMessage() {}

func (x *ADSEntry) ProtoReflect() protoreflect.Message {
	mi := &file_depot_orca_layer_v1_layer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ADSEntry.ProtoReflect.Descriptor instead.
func (*ADSEntry) Descriptor() ([]byte, []int) {
	return file_depot_orca_layer_v1_layer_proto_rawDescGZIP(), []int{5}
}

func (x *ADSEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ADSEntry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ADSEntry) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

var File_depot_orca_layer_v1_layer_proto protoreflect.FileDescriptor

var file_depot_orca_layer_v1_layer_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2f, 0x6f, 0x72, 0x63, 0x61, 0x2f, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72, 0x63, 0x61, 0x2e, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x0c, 0x4c, 0x61, 0x79, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x39, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74,
	0x2e, 0x6f, 0x72, 0x63, 0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c,
	0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x22, 0x84, 0x03, 0x0a, 0x0a, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73, 0x69, 0x7a, 0x65, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72, 0x63,
	0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x67, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x05,
	0x6d, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x6d, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x61,
	0x6a, 0x6f, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x12, 0x30, 0x0a, 0x05, 0x78, 0x61, 0x74, 0x74, 0x72, 0x18,
	0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72,
	0x63, 0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x58, 0x41, 0x74, 0x74,
	0x72, 0x52, 0x05, 0x78, 0x61, 0x74, 0x74, 0x72, 0x12, 0x2f, 0x0a, 0x03, 0x61, 0x64, 0x73, 0x18,
	0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72,
	0x63, 0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x44, 0x53, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x61, 0x64, 0x73, 0x22, 0x5b, 0x0a, 0x05, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73, 0x69, 0x7a, 0x65, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x12, 0x33, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72, 0x63, 0x61, 0x2e, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x52, 0x06,
	0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x9c, 0x01, 0x0a, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x12, 0x43, 0x0a, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2e, 0x6f, 0x72, 0x63,
	0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x09, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x22, 0x3b, 0x0a, 0x09, 0x41, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x19, 0x0a, 0x15, 0x41, 0x4c, 0x47, 0x4f, 0x52, 0x49, 0x54,
	0x48, 0x4d, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x13, 0x0a, 0x0f, 0x41, 0x4c, 0x47, 0x4f, 0x52, 0x49, 0x54, 0x48, 0x4d, 0x5f, 0x58, 0x58,
	0x48, 0x36, 0x34, 0x10, 0x01, 0x22, 0x2f, 0x0a, 0x05, 0x58, 0x41, 0x74, 0x74, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x4a, 0x0a, 0x08, 0x41, 0x44, 0x53, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69,
	0x67, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65,
	0x73, 0x74, 0x42, 0xcd, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74,
	0x2e, 0x6f, 0x72, 0x63, 0x61, 0x2e, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x0a,
	0x4c, 0x61, 0x79, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2f, 0x6f,
	0x72, 0x63, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2f,
	0x6f, 0x72, 0x63, 0x61, 0x2f, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x44, 0x4f, 0x4c, 0xaa, 0x02, 0x13, 0x44, 0x65,
	0x70, 0x6f, 0x74, 0x2e, 0x4f, 0x72, 0x63, 0x61, 0x2e, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x13, 0x44, 0x65, 0x70, 0x6f, 0x74, 0x5c, 0x4f, 0x72, 0x63, 0x61, 0x5c, 0x4c,
	0x61, 0x79, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1f, 0x44, 0x65, 0x70, 0x6f, 0x74, 0x5c,
	0x4f, 0x72, 0x63, 0x61, 0x5c, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x16, 0x44, 0x65, 0x70, 0x6f,
	0x74, 0x3a, 0x3a, 0x4f, 0x72, 0x63, 0x61, 0x3a, 0x3a, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_depot_orca_layer_v1_layer_proto_rawDescOnce sync.Once
	file_depot_orca_layer_v1_layer_proto_rawDescData = file_depot_orca_layer_v1_layer_proto_rawDesc
)

func file_depot_orca_layer_v1_layer_proto_rawDescGZIP() []byte {
	file_depot_orca_layer_v1_layer_proto_rawDescOnce.Do(func() {
		file_depot_orca_layer_v1_layer_proto_rawDescData = protoimpl.X.CompressGZIP(file_depot_orca_layer_v1_layer_proto_rawDescData)
	})
	return file_depot_orca_layer_v1_layer_proto_rawDescData
}

var file_depot_orca_layer_v1_layer_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_depot_orca_layer_v1_layer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_depot_orca_layer_v1_layer_proto_goTypes = []interface{}{
	(Digest_Algorithm)(0),         // 0: depot.orca.layer.v1.Digest.Algorithm
	(*LayerEntries)(nil),          // 1: depot.orca.layer.v1.LayerEntries
	(*LayerEntry)(nil),            // 2: depot.orca.layer.v1.LayerEntry
	(*Block)(nil),                 // 3: depot.orca.layer.v1.Block
	(*Digest)(nil),                // 4: depot.orca.layer.v1.Digest
	(*XAttr)(nil),                 // 5: depot.orca.layer.v1.XAttr
	(*ADSEntry)(nil),              // 6: depot.orca.layer.v1.ADSEntry
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_depot_orca_layer_v1_layer_proto_depIdxs = []int32{
	2, // 0: depot.orca.layer.v1.LayerEntries.entries:type_name -> depot.orca.layer.v1.LayerEntry
	3, // 1: depot.orca.layer.v1.LayerEntry.blocks:type_name -> depot.orca.layer.v1.Block
	7, // 2: depot.orca.layer.v1.LayerEntry.mtime:type_name -> google.protobuf.Timestamp
	5, // 3: depot.orca.layer.v1.LayerEntry.xattr:type_name -> depot.orca.layer.v1.XAttr
	6, // 4: depot.orca.layer.v1.LayerEntry.ads:type_name -> depot.orca.layer.v1.ADSEntry
	4, // 5: depot.orca.layer.v1.Block.digest:type_name -> depot.orca.layer.v1.Digest
	0, // 6: depot.orca.layer.v1.Digest.algorithm:type_name -> depot.orca.layer.v1.Digest.Algorithm
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_depot_orca_layer_v1_layer_proto_init() }
func file_depot_orca_layer_v1_layer_proto_init() {
	if File_depot_orca_layer_v1_layer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_depot_orca_layer_v1_layer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LayerEntries); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depot_orca_layer_v1_layer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LayerEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depot_orca_layer_v1_layer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depot_orca_layer_v1_layer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Digest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depot_orca_layer_v1_layer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*XAttr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depot_orca_layer_v1_layer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ADSEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_depot_orca_layer_v1_layer_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_depot_orca_layer_v1_layer_proto_goTypes,
		DependencyIndexes: file_depot_orca_layer_v1_layer_proto_depIdxs,
		EnumInfos:         file_depot_orca_layer_v1_layer_proto_enumTypes,
		MessageInfos:      file_depot_orca_layer_v1_layer_proto_msgTypes,
	}.Build()
	File_depot_orca_layer_v1_layer_proto = out.File
	file_depot_orca_layer_v1_layer_proto_rawDesc = nil
	file_depot_orca_layer_v1_layer_proto_goTypes = nil
	file_depot_orca_layer_v1_layer_proto_depIdxs = nil
}
