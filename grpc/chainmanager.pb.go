// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chainmanager.proto

package grpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SetMonitorsRequest struct {
	// Types that are valid to be assigned to MsgType:
	//	*SetMonitorsRequest_Preflight_
	//	*SetMonitorsRequest_Monitor
	MsgType              isSetMonitorsRequest_MsgType `protobuf_oneof:"msg_type"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SetMonitorsRequest) Reset()         { *m = SetMonitorsRequest{} }
func (m *SetMonitorsRequest) String() string { return proto.CompactTextString(m) }
func (*SetMonitorsRequest) ProtoMessage()    {}
func (*SetMonitorsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{0}
}

func (m *SetMonitorsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetMonitorsRequest.Unmarshal(m, b)
}
func (m *SetMonitorsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetMonitorsRequest.Marshal(b, m, deterministic)
}
func (m *SetMonitorsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetMonitorsRequest.Merge(m, src)
}
func (m *SetMonitorsRequest) XXX_Size() int {
	return xxx_messageInfo_SetMonitorsRequest.Size(m)
}
func (m *SetMonitorsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetMonitorsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetMonitorsRequest proto.InternalMessageInfo

type isSetMonitorsRequest_MsgType interface {
	isSetMonitorsRequest_MsgType()
}

type SetMonitorsRequest_Preflight_ struct {
	Preflight *SetMonitorsRequest_Preflight `protobuf:"bytes,1,opt,name=preflight,proto3,oneof"`
}

type SetMonitorsRequest_Monitor struct {
	Monitor *Monitor `protobuf:"bytes,2,opt,name=monitor,proto3,oneof"`
}

func (*SetMonitorsRequest_Preflight_) isSetMonitorsRequest_MsgType() {}

func (*SetMonitorsRequest_Monitor) isSetMonitorsRequest_MsgType() {}

func (m *SetMonitorsRequest) GetMsgType() isSetMonitorsRequest_MsgType {
	if m != nil {
		return m.MsgType
	}
	return nil
}

func (m *SetMonitorsRequest) GetPreflight() *SetMonitorsRequest_Preflight {
	if x, ok := m.GetMsgType().(*SetMonitorsRequest_Preflight_); ok {
		return x.Preflight
	}
	return nil
}

func (m *SetMonitorsRequest) GetMonitor() *Monitor {
	if x, ok := m.GetMsgType().(*SetMonitorsRequest_Monitor); ok {
		return x.Monitor
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SetMonitorsRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SetMonitorsRequest_OneofMarshaler, _SetMonitorsRequest_OneofUnmarshaler, _SetMonitorsRequest_OneofSizer, []interface{}{
		(*SetMonitorsRequest_Preflight_)(nil),
		(*SetMonitorsRequest_Monitor)(nil),
	}
}

func _SetMonitorsRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SetMonitorsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *SetMonitorsRequest_Preflight_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Preflight); err != nil {
			return err
		}
	case *SetMonitorsRequest_Monitor:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Monitor); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SetMonitorsRequest.MsgType has unexpected type %T", x)
	}
	return nil
}

func _SetMonitorsRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SetMonitorsRequest)
	switch tag {
	case 1: // msg_type.preflight
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SetMonitorsRequest_Preflight)
		err := b.DecodeMessage(msg)
		m.MsgType = &SetMonitorsRequest_Preflight_{msg}
		return true, err
	case 2: // msg_type.monitor
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Monitor)
		err := b.DecodeMessage(msg)
		m.MsgType = &SetMonitorsRequest_Monitor{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SetMonitorsRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SetMonitorsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *SetMonitorsRequest_Preflight_:
		s := proto.Size(x.Preflight)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *SetMonitorsRequest_Monitor:
		s := proto.Size(x.Monitor)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type SetMonitorsRequest_Preflight struct {
	// SessionId must match the one specified in RegisterChainManagerRequest.
	SessionId []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// MonitorSetVersion is the version of the new MonitorSet.
	// It must be greater than ChainManager's local version.
	MonitorSetVersion uint32 `protobuf:"varint,2,opt,name=monitor_set_version,json=monitorSetVersion,proto3" json:"monitor_set_version,omitempty"`
	// ResumeAfterBlockHash identifies the block after which the new MonitorSet should be applied.
	// If this is a known block (in the past), ChainManager should do a rewind and start reporting again.
	// Otherwise, ChainManager should apply the MonitorSet to the next incoming block.
	ResumeAfterBlockHash []byte   `protobuf:"bytes,3,opt,name=resume_after_block_hash,json=resumeAfterBlockHash,proto3" json:"resume_after_block_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetMonitorsRequest_Preflight) Reset()         { *m = SetMonitorsRequest_Preflight{} }
func (m *SetMonitorsRequest_Preflight) String() string { return proto.CompactTextString(m) }
func (*SetMonitorsRequest_Preflight) ProtoMessage()    {}
func (*SetMonitorsRequest_Preflight) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{0, 0}
}

func (m *SetMonitorsRequest_Preflight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetMonitorsRequest_Preflight.Unmarshal(m, b)
}
func (m *SetMonitorsRequest_Preflight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetMonitorsRequest_Preflight.Marshal(b, m, deterministic)
}
func (m *SetMonitorsRequest_Preflight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetMonitorsRequest_Preflight.Merge(m, src)
}
func (m *SetMonitorsRequest_Preflight) XXX_Size() int {
	return xxx_messageInfo_SetMonitorsRequest_Preflight.Size(m)
}
func (m *SetMonitorsRequest_Preflight) XXX_DiscardUnknown() {
	xxx_messageInfo_SetMonitorsRequest_Preflight.DiscardUnknown(m)
}

var xxx_messageInfo_SetMonitorsRequest_Preflight proto.InternalMessageInfo

func (m *SetMonitorsRequest_Preflight) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

func (m *SetMonitorsRequest_Preflight) GetMonitorSetVersion() uint32 {
	if m != nil {
		return m.MonitorSetVersion
	}
	return 0
}

func (m *SetMonitorsRequest_Preflight) GetResumeAfterBlockHash() []byte {
	if m != nil {
		return m.ResumeAfterBlockHash
	}
	return nil
}

type UpdateMonitorsRequest struct {
	// Types that are valid to be assigned to MsgType:
	//	*UpdateMonitorsRequest_Preflight_
	//	*UpdateMonitorsRequest_AddMonitor_
	//	*UpdateMonitorsRequest_RemoveMonitor_
	MsgType              isUpdateMonitorsRequest_MsgType `protobuf_oneof:"msg_type"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *UpdateMonitorsRequest) Reset()         { *m = UpdateMonitorsRequest{} }
func (m *UpdateMonitorsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateMonitorsRequest) ProtoMessage()    {}
func (*UpdateMonitorsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{1}
}

func (m *UpdateMonitorsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateMonitorsRequest.Unmarshal(m, b)
}
func (m *UpdateMonitorsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateMonitorsRequest.Marshal(b, m, deterministic)
}
func (m *UpdateMonitorsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateMonitorsRequest.Merge(m, src)
}
func (m *UpdateMonitorsRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateMonitorsRequest.Size(m)
}
func (m *UpdateMonitorsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateMonitorsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateMonitorsRequest proto.InternalMessageInfo

type isUpdateMonitorsRequest_MsgType interface {
	isUpdateMonitorsRequest_MsgType()
}

type UpdateMonitorsRequest_Preflight_ struct {
	Preflight *UpdateMonitorsRequest_Preflight `protobuf:"bytes,1,opt,name=preflight,proto3,oneof"`
}

type UpdateMonitorsRequest_AddMonitor_ struct {
	AddMonitor *UpdateMonitorsRequest_AddMonitor `protobuf:"bytes,2,opt,name=add_monitor,json=addMonitor,proto3,oneof"`
}

type UpdateMonitorsRequest_RemoveMonitor_ struct {
	RemoveMonitor *UpdateMonitorsRequest_RemoveMonitor `protobuf:"bytes,3,opt,name=remove_monitor,json=removeMonitor,proto3,oneof"`
}

func (*UpdateMonitorsRequest_Preflight_) isUpdateMonitorsRequest_MsgType() {}

func (*UpdateMonitorsRequest_AddMonitor_) isUpdateMonitorsRequest_MsgType() {}

func (*UpdateMonitorsRequest_RemoveMonitor_) isUpdateMonitorsRequest_MsgType() {}

func (m *UpdateMonitorsRequest) GetMsgType() isUpdateMonitorsRequest_MsgType {
	if m != nil {
		return m.MsgType
	}
	return nil
}

func (m *UpdateMonitorsRequest) GetPreflight() *UpdateMonitorsRequest_Preflight {
	if x, ok := m.GetMsgType().(*UpdateMonitorsRequest_Preflight_); ok {
		return x.Preflight
	}
	return nil
}

func (m *UpdateMonitorsRequest) GetAddMonitor() *UpdateMonitorsRequest_AddMonitor {
	if x, ok := m.GetMsgType().(*UpdateMonitorsRequest_AddMonitor_); ok {
		return x.AddMonitor
	}
	return nil
}

func (m *UpdateMonitorsRequest) GetRemoveMonitor() *UpdateMonitorsRequest_RemoveMonitor {
	if x, ok := m.GetMsgType().(*UpdateMonitorsRequest_RemoveMonitor_); ok {
		return x.RemoveMonitor
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*UpdateMonitorsRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _UpdateMonitorsRequest_OneofMarshaler, _UpdateMonitorsRequest_OneofUnmarshaler, _UpdateMonitorsRequest_OneofSizer, []interface{}{
		(*UpdateMonitorsRequest_Preflight_)(nil),
		(*UpdateMonitorsRequest_AddMonitor_)(nil),
		(*UpdateMonitorsRequest_RemoveMonitor_)(nil),
	}
}

func _UpdateMonitorsRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*UpdateMonitorsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *UpdateMonitorsRequest_Preflight_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Preflight); err != nil {
			return err
		}
	case *UpdateMonitorsRequest_AddMonitor_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AddMonitor); err != nil {
			return err
		}
	case *UpdateMonitorsRequest_RemoveMonitor_:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.RemoveMonitor); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("UpdateMonitorsRequest.MsgType has unexpected type %T", x)
	}
	return nil
}

func _UpdateMonitorsRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*UpdateMonitorsRequest)
	switch tag {
	case 1: // msg_type.preflight
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpdateMonitorsRequest_Preflight)
		err := b.DecodeMessage(msg)
		m.MsgType = &UpdateMonitorsRequest_Preflight_{msg}
		return true, err
	case 2: // msg_type.add_monitor
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpdateMonitorsRequest_AddMonitor)
		err := b.DecodeMessage(msg)
		m.MsgType = &UpdateMonitorsRequest_AddMonitor_{msg}
		return true, err
	case 3: // msg_type.remove_monitor
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpdateMonitorsRequest_RemoveMonitor)
		err := b.DecodeMessage(msg)
		m.MsgType = &UpdateMonitorsRequest_RemoveMonitor_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _UpdateMonitorsRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*UpdateMonitorsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *UpdateMonitorsRequest_Preflight_:
		s := proto.Size(x.Preflight)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *UpdateMonitorsRequest_AddMonitor_:
		s := proto.Size(x.AddMonitor)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *UpdateMonitorsRequest_RemoveMonitor_:
		s := proto.Size(x.RemoveMonitor)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type UpdateMonitorsRequest_Preflight struct {
	// SessionId must match the one specified in RegisterChainManagerRequest.
	SessionId []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// PreviousMonitorSetVersion is the version of the previous MonitorSet.
	// It must match ChainManager's local version.
	PreviousMonitorSetVersion uint32 `protobuf:"varint,2,opt,name=previous_monitor_set_version,json=previousMonitorSetVersion,proto3" json:"previous_monitor_set_version,omitempty"`
	// MonitorSetVersion is the version of the new MonitorSet.
	// It must be greater than PreviousMonitorSetVersion.
	MonitorSetVersion uint32 `protobuf:"varint,3,opt,name=monitor_set_version,json=monitorSetVersion,proto3" json:"monitor_set_version,omitempty"`
	// ResumeAfterBlockHash identifies the block after which the new MonitorSet should be applied.
	// If this is a known block (in the past), ChainManager should do a rewind and start reporting again.
	// Otherwise, ChainManager should apply the MonitorSet to the next incoming block.
	ResumeAfterBlockHash []byte   `protobuf:"bytes,4,opt,name=resume_after_block_hash,json=resumeAfterBlockHash,proto3" json:"resume_after_block_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateMonitorsRequest_Preflight) Reset()         { *m = UpdateMonitorsRequest_Preflight{} }
func (m *UpdateMonitorsRequest_Preflight) String() string { return proto.CompactTextString(m) }
func (*UpdateMonitorsRequest_Preflight) ProtoMessage()    {}
func (*UpdateMonitorsRequest_Preflight) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{1, 0}
}

func (m *UpdateMonitorsRequest_Preflight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateMonitorsRequest_Preflight.Unmarshal(m, b)
}
func (m *UpdateMonitorsRequest_Preflight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateMonitorsRequest_Preflight.Marshal(b, m, deterministic)
}
func (m *UpdateMonitorsRequest_Preflight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateMonitorsRequest_Preflight.Merge(m, src)
}
func (m *UpdateMonitorsRequest_Preflight) XXX_Size() int {
	return xxx_messageInfo_UpdateMonitorsRequest_Preflight.Size(m)
}
func (m *UpdateMonitorsRequest_Preflight) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateMonitorsRequest_Preflight.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateMonitorsRequest_Preflight proto.InternalMessageInfo

func (m *UpdateMonitorsRequest_Preflight) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

func (m *UpdateMonitorsRequest_Preflight) GetPreviousMonitorSetVersion() uint32 {
	if m != nil {
		return m.PreviousMonitorSetVersion
	}
	return 0
}

func (m *UpdateMonitorsRequest_Preflight) GetMonitorSetVersion() uint32 {
	if m != nil {
		return m.MonitorSetVersion
	}
	return 0
}

func (m *UpdateMonitorsRequest_Preflight) GetResumeAfterBlockHash() []byte {
	if m != nil {
		return m.ResumeAfterBlockHash
	}
	return nil
}

type UpdateMonitorsRequest_AddMonitor struct {
	Monitor              *Monitor `protobuf:"bytes,1,opt,name=monitor,proto3" json:"monitor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateMonitorsRequest_AddMonitor) Reset()         { *m = UpdateMonitorsRequest_AddMonitor{} }
func (m *UpdateMonitorsRequest_AddMonitor) String() string { return proto.CompactTextString(m) }
func (*UpdateMonitorsRequest_AddMonitor) ProtoMessage()    {}
func (*UpdateMonitorsRequest_AddMonitor) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{1, 1}
}

func (m *UpdateMonitorsRequest_AddMonitor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateMonitorsRequest_AddMonitor.Unmarshal(m, b)
}
func (m *UpdateMonitorsRequest_AddMonitor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateMonitorsRequest_AddMonitor.Marshal(b, m, deterministic)
}
func (m *UpdateMonitorsRequest_AddMonitor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateMonitorsRequest_AddMonitor.Merge(m, src)
}
func (m *UpdateMonitorsRequest_AddMonitor) XXX_Size() int {
	return xxx_messageInfo_UpdateMonitorsRequest_AddMonitor.Size(m)
}
func (m *UpdateMonitorsRequest_AddMonitor) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateMonitorsRequest_AddMonitor.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateMonitorsRequest_AddMonitor proto.InternalMessageInfo

func (m *UpdateMonitorsRequest_AddMonitor) GetMonitor() *Monitor {
	if m != nil {
		return m.Monitor
	}
	return nil
}

type UpdateMonitorsRequest_RemoveMonitor struct {
	MonitorId            int64    `protobuf:"varint,1,opt,name=monitor_id,json=monitorId,proto3" json:"monitor_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateMonitorsRequest_RemoveMonitor) Reset()         { *m = UpdateMonitorsRequest_RemoveMonitor{} }
func (m *UpdateMonitorsRequest_RemoveMonitor) String() string { return proto.CompactTextString(m) }
func (*UpdateMonitorsRequest_RemoveMonitor) ProtoMessage()    {}
func (*UpdateMonitorsRequest_RemoveMonitor) Descriptor() ([]byte, []int) {
	return fileDescriptor_c87c4370f7161d54, []int{1, 2}
}

func (m *UpdateMonitorsRequest_RemoveMonitor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor.Unmarshal(m, b)
}
func (m *UpdateMonitorsRequest_RemoveMonitor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor.Marshal(b, m, deterministic)
}
func (m *UpdateMonitorsRequest_RemoveMonitor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor.Merge(m, src)
}
func (m *UpdateMonitorsRequest_RemoveMonitor) XXX_Size() int {
	return xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor.Size(m)
}
func (m *UpdateMonitorsRequest_RemoveMonitor) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateMonitorsRequest_RemoveMonitor proto.InternalMessageInfo

func (m *UpdateMonitorsRequest_RemoveMonitor) GetMonitorId() int64 {
	if m != nil {
		return m.MonitorId
	}
	return 0
}

func init() {
	proto.RegisterType((*SetMonitorsRequest)(nil), "grpc.SetMonitorsRequest")
	proto.RegisterType((*SetMonitorsRequest_Preflight)(nil), "grpc.SetMonitorsRequest.Preflight")
	proto.RegisterType((*UpdateMonitorsRequest)(nil), "grpc.UpdateMonitorsRequest")
	proto.RegisterType((*UpdateMonitorsRequest_Preflight)(nil), "grpc.UpdateMonitorsRequest.Preflight")
	proto.RegisterType((*UpdateMonitorsRequest_AddMonitor)(nil), "grpc.UpdateMonitorsRequest.AddMonitor")
	proto.RegisterType((*UpdateMonitorsRequest_RemoveMonitor)(nil), "grpc.UpdateMonitorsRequest.RemoveMonitor")
}

func init() { proto.RegisterFile("chainmanager.proto", fileDescriptor_c87c4370f7161d54) }

var fileDescriptor_c87c4370f7161d54 = []byte{
	// 462 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x5d, 0x6f, 0x12, 0x41,
	0x14, 0xed, 0x16, 0xa2, 0x72, 0x29, 0x4d, 0x1c, 0xbf, 0x70, 0xab, 0x89, 0x21, 0x51, 0xdb, 0x97,
	0x69, 0x52, 0xd3, 0x67, 0x03, 0xa6, 0x11, 0x1e, 0x48, 0xcc, 0x34, 0xfa, 0x3a, 0x19, 0x98, 0xcb,
	0xee, 0x46, 0x96, 0x59, 0x67, 0x06, 0x92, 0xfe, 0x0c, 0x5f, 0x7d, 0xf5, 0xb7, 0xf8, 0x07, 0xfc,
	0x45, 0x66, 0x67, 0x76, 0x59, 0x10, 0x96, 0x44, 0x1f, 0xb9, 0xe7, 0xdc, 0x33, 0xec, 0x39, 0xf7,
	0x00, 0x99, 0xc6, 0x22, 0x59, 0xa4, 0x62, 0x21, 0x22, 0xd4, 0x34, 0xd3, 0xca, 0x2a, 0xd2, 0x8c,
	0x74, 0x36, 0x0d, 0xcf, 0x22, 0xa5, 0xa2, 0x39, 0x5e, 0xba, 0xd9, 0x64, 0x39, 0xbb, 0xc4, 0x34,
	0xb3, 0x77, 0x9e, 0x12, 0x12, 0xa5, 0xa7, 0x31, 0x1a, 0xab, 0x85, 0x55, 0xc5, 0x5a, 0xef, 0xe7,
	0x31, 0x90, 0x5b, 0xb4, 0x63, 0xb5, 0x48, 0xac, 0xd2, 0x86, 0xe1, 0xb7, 0x25, 0x1a, 0x4b, 0x06,
	0xd0, 0xca, 0x34, 0xce, 0xe6, 0x49, 0x14, 0xdb, 0x6e, 0xf0, 0x2a, 0x38, 0x6f, 0x5f, 0xf5, 0x68,
	0xfe, 0x02, 0xdd, 0x25, 0xd3, 0x4f, 0x25, 0x73, 0x78, 0xc4, 0xaa, 0x35, 0x72, 0x01, 0xf7, 0x53,
	0xcf, 0xec, 0x1e, 0x3b, 0x85, 0x8e, 0x57, 0x28, 0xd6, 0x87, 0x47, 0xac, 0xc4, 0xc3, 0xef, 0x01,
	0xb4, 0xd6, 0x2a, 0xe4, 0x25, 0x80, 0x41, 0x63, 0x12, 0xb5, 0xe0, 0x89, 0x74, 0xaf, 0x9f, 0xb0,
	0x56, 0x31, 0x19, 0x49, 0x42, 0xe1, 0x51, 0xb1, 0xc7, 0x0d, 0x5a, 0xbe, 0x42, 0x9d, 0x03, 0xee,
	0x8d, 0x0e, 0x7b, 0x58, 0x40, 0xb7, 0x68, 0xbf, 0x78, 0x80, 0x5c, 0xc3, 0x33, 0x8d, 0x66, 0x99,
	0x22, 0x17, 0x33, 0x8b, 0x9a, 0x4f, 0xe6, 0x6a, 0xfa, 0x95, 0xc7, 0xc2, 0xc4, 0xdd, 0x86, 0xd3,
	0x7e, 0xec, 0xe1, 0x7e, 0x8e, 0x0e, 0x72, 0x70, 0x28, 0x4c, 0x3c, 0x00, 0x78, 0x90, 0x9a, 0x88,
	0xdb, 0xbb, 0x0c, 0x7b, 0xbf, 0x9a, 0xf0, 0xe4, 0x73, 0x26, 0x85, 0xc5, 0xbf, 0x8d, 0xba, 0xd9,
	0x35, 0xea, 0xb5, 0xff, 0xcc, 0xbd, 0xfc, 0x3a, 0xaf, 0x46, 0xd0, 0x16, 0x52, 0xf2, 0x6d, 0xbf,
	0xde, 0x1c, 0x12, 0xea, 0x4b, 0x59, 0x19, 0x09, 0x62, 0xfd, 0x8b, 0x30, 0x38, 0xd5, 0x98, 0xaa,
	0x15, 0xae, 0xd5, 0x1a, 0x4e, 0xed, 0xe2, 0x90, 0x1a, 0x73, 0x1b, 0x95, 0x60, 0x47, 0x6f, 0x0e,
	0xc2, 0xdf, 0xff, 0x92, 0xcf, 0x7b, 0x78, 0x91, 0x69, 0x5c, 0x25, 0x6a, 0x69, 0x78, 0x7d, 0x50,
	0xcf, 0x4b, 0xce, 0x78, 0x27, 0xb0, 0x9a, 0x80, 0x1b, 0xff, 0x11, 0x70, 0xb3, 0x3e, 0xe0, 0xf0,
	0x1a, 0xa0, 0x32, 0x91, 0xbc, 0xad, 0xae, 0x35, 0xd8, 0x73, 0xad, 0xd5, 0xad, 0x52, 0xe8, 0x6c,
	0xb9, 0x95, 0xdb, 0x51, 0xfe, 0xdd, 0xc2, 0x8e, 0x06, 0x6b, 0x15, 0x93, 0x91, 0xdc, 0xbc, 0xa3,
	0xab, 0x1f, 0x01, 0x9c, 0x7c, 0xc8, 0xbb, 0x3b, 0xf6, 0xdd, 0x25, 0x7d, 0x68, 0x6f, 0x14, 0x8a,
	0x74, 0xeb, 0x3a, 0x16, 0x3e, 0xa5, 0xbe, 0xd9, 0xb4, 0x6c, 0x36, 0xbd, 0xc9, 0x9b, 0x7d, 0x1e,
	0x90, 0x8f, 0x70, 0xba, 0x9d, 0x29, 0x39, 0x3b, 0x90, 0x74, 0xbd, 0xd0, 0xe4, 0x9e, 0x9b, 0xbc,
	0xfb, 0x13, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x17, 0x23, 0x02, 0x5e, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChainManagerClient is the client API for ChainManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChainManagerClient interface {
	// SetMonitors resets ChainManager's current MonitorSet.
	// The request stream must begin with a valid Preflight, followed by zero or more Monitors.
	SetMonitors(ctx context.Context, opts ...grpc.CallOption) (ChainManager_SetMonitorsClient, error)
	// UpdateMonitors patches ChainManager's current MonitorSet with additions and removals.
	// The request stream must begin with a valid Preflight, followed by zero or more AddMonitor or RemoveMonitor requests.
	// The MonitorSetVersion in Preflight must be exactly the local version plus one.
	UpdateMonitors(ctx context.Context, opts ...grpc.CallOption) (ChainManager_UpdateMonitorsClient, error)
}

type chainManagerClient struct {
	cc *grpc.ClientConn
}

func NewChainManagerClient(cc *grpc.ClientConn) ChainManagerClient {
	return &chainManagerClient{cc}
}

func (c *chainManagerClient) SetMonitors(ctx context.Context, opts ...grpc.CallOption) (ChainManager_SetMonitorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChainManager_serviceDesc.Streams[0], "/grpc.ChainManager/SetMonitors", opts...)
	if err != nil {
		return nil, err
	}
	x := &chainManagerSetMonitorsClient{stream}
	return x, nil
}

type ChainManager_SetMonitorsClient interface {
	Send(*SetMonitorsRequest) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type chainManagerSetMonitorsClient struct {
	grpc.ClientStream
}

func (x *chainManagerSetMonitorsClient) Send(m *SetMonitorsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chainManagerSetMonitorsClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chainManagerClient) UpdateMonitors(ctx context.Context, opts ...grpc.CallOption) (ChainManager_UpdateMonitorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChainManager_serviceDesc.Streams[1], "/grpc.ChainManager/UpdateMonitors", opts...)
	if err != nil {
		return nil, err
	}
	x := &chainManagerUpdateMonitorsClient{stream}
	return x, nil
}

type ChainManager_UpdateMonitorsClient interface {
	Send(*UpdateMonitorsRequest) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type chainManagerUpdateMonitorsClient struct {
	grpc.ClientStream
}

func (x *chainManagerUpdateMonitorsClient) Send(m *UpdateMonitorsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chainManagerUpdateMonitorsClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChainManagerServer is the server API for ChainManager service.
type ChainManagerServer interface {
	// SetMonitors resets ChainManager's current MonitorSet.
	// The request stream must begin with a valid Preflight, followed by zero or more Monitors.
	SetMonitors(ChainManager_SetMonitorsServer) error
	// UpdateMonitors patches ChainManager's current MonitorSet with additions and removals.
	// The request stream must begin with a valid Preflight, followed by zero or more AddMonitor or RemoveMonitor requests.
	// The MonitorSetVersion in Preflight must be exactly the local version plus one.
	UpdateMonitors(ChainManager_UpdateMonitorsServer) error
}

func RegisterChainManagerServer(s *grpc.Server, srv ChainManagerServer) {
	s.RegisterService(&_ChainManager_serviceDesc, srv)
}

func _ChainManager_SetMonitors_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChainManagerServer).SetMonitors(&chainManagerSetMonitorsServer{stream})
}

type ChainManager_SetMonitorsServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*SetMonitorsRequest, error)
	grpc.ServerStream
}

type chainManagerSetMonitorsServer struct {
	grpc.ServerStream
}

func (x *chainManagerSetMonitorsServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chainManagerSetMonitorsServer) Recv() (*SetMonitorsRequest, error) {
	m := new(SetMonitorsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChainManager_UpdateMonitors_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChainManagerServer).UpdateMonitors(&chainManagerUpdateMonitorsServer{stream})
}

type ChainManager_UpdateMonitorsServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*UpdateMonitorsRequest, error)
	grpc.ServerStream
}

type chainManagerUpdateMonitorsServer struct {
	grpc.ServerStream
}

func (x *chainManagerUpdateMonitorsServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chainManagerUpdateMonitorsServer) Recv() (*UpdateMonitorsRequest, error) {
	m := new(UpdateMonitorsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ChainManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.ChainManager",
	HandlerType: (*ChainManagerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetMonitors",
			Handler:       _ChainManager_SetMonitors_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UpdateMonitors",
			Handler:       _ChainManager_UpdateMonitors_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "chainmanager.proto",
}