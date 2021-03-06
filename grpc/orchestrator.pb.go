// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/orchestrator.proto

package grpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	protos "github.com/megaspacelab/megaconnect/protos"
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

type RegisterChainManagerRequest struct {
	// ChainManagerId identifies the instance of the requesting ChainManager.
	ChainManagerId *InstanceId `protobuf:"bytes,1,opt,name=chain_manager_id,json=chainManagerId,proto3" json:"chain_manager_id,omitempty"`
	// ChainId identifies the chain being managed.
	ChainId string `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// ListenPort allows Orchestrator to reach out to this ChainManager.
	ListenPort int32 `protobuf:"varint,3,opt,name=listen_port,json=listenPort,proto3" json:"listen_port,omitempty"`
	// SessionId should be included in all outbound communications from Orchestrator to ChainManager as a validity proof.
	SessionId            []byte   `protobuf:"bytes,4,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterChainManagerRequest) Reset()         { *m = RegisterChainManagerRequest{} }
func (m *RegisterChainManagerRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterChainManagerRequest) ProtoMessage()    {}
func (*RegisterChainManagerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{0}
}

func (m *RegisterChainManagerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterChainManagerRequest.Unmarshal(m, b)
}
func (m *RegisterChainManagerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterChainManagerRequest.Marshal(b, m, deterministic)
}
func (m *RegisterChainManagerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterChainManagerRequest.Merge(m, src)
}
func (m *RegisterChainManagerRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterChainManagerRequest.Size(m)
}
func (m *RegisterChainManagerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterChainManagerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterChainManagerRequest proto.InternalMessageInfo

func (m *RegisterChainManagerRequest) GetChainManagerId() *InstanceId {
	if m != nil {
		return m.ChainManagerId
	}
	return nil
}

func (m *RegisterChainManagerRequest) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *RegisterChainManagerRequest) GetListenPort() int32 {
	if m != nil {
		return m.ListenPort
	}
	return 0
}

func (m *RegisterChainManagerRequest) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

type RegisterChainManagerResponse struct {
	// Lease indicates the validity period allocated to the requesting ChainManager.
	Lease *Lease `protobuf:"bytes,1,opt,name=lease,proto3" json:"lease,omitempty"`
	// ResumeAfter indicates the last processed block by Orchestrator.
	// If not nil, ChainManager is expected to start reporting from right after this block.
	// Otherwise, ChainManager can start reporting from the next new block.
	ResumeAfter *BlockSpec `protobuf:"bytes,2,opt,name=resume_after,json=resumeAfter,proto3" json:"resume_after,omitempty"`
	// Monitors contain all things to be monitored by the requesting ChainManager.
	Monitors             *MonitorSet `protobuf:"bytes,3,opt,name=monitors,proto3" json:"monitors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RegisterChainManagerResponse) Reset()         { *m = RegisterChainManagerResponse{} }
func (m *RegisterChainManagerResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterChainManagerResponse) ProtoMessage()    {}
func (*RegisterChainManagerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{1}
}

func (m *RegisterChainManagerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterChainManagerResponse.Unmarshal(m, b)
}
func (m *RegisterChainManagerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterChainManagerResponse.Marshal(b, m, deterministic)
}
func (m *RegisterChainManagerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterChainManagerResponse.Merge(m, src)
}
func (m *RegisterChainManagerResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterChainManagerResponse.Size(m)
}
func (m *RegisterChainManagerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterChainManagerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterChainManagerResponse proto.InternalMessageInfo

func (m *RegisterChainManagerResponse) GetLease() *Lease {
	if m != nil {
		return m.Lease
	}
	return nil
}

func (m *RegisterChainManagerResponse) GetResumeAfter() *BlockSpec {
	if m != nil {
		return m.ResumeAfter
	}
	return nil
}

func (m *RegisterChainManagerResponse) GetMonitors() *MonitorSet {
	if m != nil {
		return m.Monitors
	}
	return nil
}

type UnregisterChainManagerRequest struct {
	// ChainId identifies the chain being managed.
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// LeaseId must be the id of a still valid Lease of to be unregistered chain manager.
	LeaseId []byte `protobuf:"bytes,2,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	// Message indicates the reason of the unregistration.
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnregisterChainManagerRequest) Reset()         { *m = UnregisterChainManagerRequest{} }
func (m *UnregisterChainManagerRequest) String() string { return proto.CompactTextString(m) }
func (*UnregisterChainManagerRequest) ProtoMessage()    {}
func (*UnregisterChainManagerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{2}
}

func (m *UnregisterChainManagerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnregisterChainManagerRequest.Unmarshal(m, b)
}
func (m *UnregisterChainManagerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnregisterChainManagerRequest.Marshal(b, m, deterministic)
}
func (m *UnregisterChainManagerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnregisterChainManagerRequest.Merge(m, src)
}
func (m *UnregisterChainManagerRequest) XXX_Size() int {
	return xxx_messageInfo_UnregisterChainManagerRequest.Size(m)
}
func (m *UnregisterChainManagerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UnregisterChainManagerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UnregisterChainManagerRequest proto.InternalMessageInfo

func (m *UnregisterChainManagerRequest) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *UnregisterChainManagerRequest) GetLeaseId() []byte {
	if m != nil {
		return m.LeaseId
	}
	return nil
}

func (m *UnregisterChainManagerRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// InstanceId identifies a specific instance of an entity.
type InstanceId struct {
	// Id is the entity id.
	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Instance is a monotonically increasing sequence number, namespaced by Id.
	Instance             uint32   `protobuf:"varint,2,opt,name=instance,proto3" json:"instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstanceId) Reset()         { *m = InstanceId{} }
func (m *InstanceId) String() string { return proto.CompactTextString(m) }
func (*InstanceId) ProtoMessage()    {}
func (*InstanceId) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{3}
}

func (m *InstanceId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstanceId.Unmarshal(m, b)
}
func (m *InstanceId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstanceId.Marshal(b, m, deterministic)
}
func (m *InstanceId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstanceId.Merge(m, src)
}
func (m *InstanceId) XXX_Size() int {
	return xxx_messageInfo_InstanceId.Size(m)
}
func (m *InstanceId) XXX_DiscardUnknown() {
	xxx_messageInfo_InstanceId.DiscardUnknown(m)
}

var xxx_messageInfo_InstanceId proto.InternalMessageInfo

func (m *InstanceId) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *InstanceId) GetInstance() uint32 {
	if m != nil {
		return m.Instance
	}
	return 0
}

type RenewLeaseRequest struct {
	// LeaseId must be the id of a still valid Lease.
	LeaseId              []byte   `protobuf:"bytes,1,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RenewLeaseRequest) Reset()         { *m = RenewLeaseRequest{} }
func (m *RenewLeaseRequest) String() string { return proto.CompactTextString(m) }
func (*RenewLeaseRequest) ProtoMessage()    {}
func (*RenewLeaseRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{4}
}

func (m *RenewLeaseRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RenewLeaseRequest.Unmarshal(m, b)
}
func (m *RenewLeaseRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RenewLeaseRequest.Marshal(b, m, deterministic)
}
func (m *RenewLeaseRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RenewLeaseRequest.Merge(m, src)
}
func (m *RenewLeaseRequest) XXX_Size() int {
	return xxx_messageInfo_RenewLeaseRequest.Size(m)
}
func (m *RenewLeaseRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RenewLeaseRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RenewLeaseRequest proto.InternalMessageInfo

func (m *RenewLeaseRequest) GetLeaseId() []byte {
	if m != nil {
		return m.LeaseId
	}
	return nil
}

type Lease struct {
	// Id uniquely identifies this Lease.
	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// RemainingSeconds indicates the length of the validity period in seconds.
	RemainingSeconds     uint32   `protobuf:"varint,2,opt,name=remaining_seconds,json=remainingSeconds,proto3" json:"remaining_seconds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Lease) Reset()         { *m = Lease{} }
func (m *Lease) String() string { return proto.CompactTextString(m) }
func (*Lease) ProtoMessage()    {}
func (*Lease) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{5}
}

func (m *Lease) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lease.Unmarshal(m, b)
}
func (m *Lease) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lease.Marshal(b, m, deterministic)
}
func (m *Lease) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lease.Merge(m, src)
}
func (m *Lease) XXX_Size() int {
	return xxx_messageInfo_Lease.Size(m)
}
func (m *Lease) XXX_DiscardUnknown() {
	xxx_messageInfo_Lease.DiscardUnknown(m)
}

var xxx_messageInfo_Lease proto.InternalMessageInfo

func (m *Lease) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Lease) GetRemainingSeconds() uint32 {
	if m != nil {
		return m.RemainingSeconds
	}
	return 0
}

type ReportBlockEventsRequest struct {
	// Types that are valid to be assigned to MsgType:
	//	*ReportBlockEventsRequest_Preflight_
	//	*ReportBlockEventsRequest_Block
	//	*ReportBlockEventsRequest_Event
	MsgType              isReportBlockEventsRequest_MsgType `protobuf_oneof:"msg_type"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *ReportBlockEventsRequest) Reset()         { *m = ReportBlockEventsRequest{} }
func (m *ReportBlockEventsRequest) String() string { return proto.CompactTextString(m) }
func (*ReportBlockEventsRequest) ProtoMessage()    {}
func (*ReportBlockEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{6}
}

func (m *ReportBlockEventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportBlockEventsRequest.Unmarshal(m, b)
}
func (m *ReportBlockEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportBlockEventsRequest.Marshal(b, m, deterministic)
}
func (m *ReportBlockEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportBlockEventsRequest.Merge(m, src)
}
func (m *ReportBlockEventsRequest) XXX_Size() int {
	return xxx_messageInfo_ReportBlockEventsRequest.Size(m)
}
func (m *ReportBlockEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportBlockEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportBlockEventsRequest proto.InternalMessageInfo

type isReportBlockEventsRequest_MsgType interface {
	isReportBlockEventsRequest_MsgType()
}

type ReportBlockEventsRequest_Preflight_ struct {
	Preflight *ReportBlockEventsRequest_Preflight `protobuf:"bytes,1,opt,name=preflight,proto3,oneof"`
}

type ReportBlockEventsRequest_Block struct {
	Block *protos.Block `protobuf:"bytes,2,opt,name=block,proto3,oneof"`
}

type ReportBlockEventsRequest_Event struct {
	Event *protos.MonitorEvent `protobuf:"bytes,3,opt,name=event,proto3,oneof"`
}

func (*ReportBlockEventsRequest_Preflight_) isReportBlockEventsRequest_MsgType() {}

func (*ReportBlockEventsRequest_Block) isReportBlockEventsRequest_MsgType() {}

func (*ReportBlockEventsRequest_Event) isReportBlockEventsRequest_MsgType() {}

func (m *ReportBlockEventsRequest) GetMsgType() isReportBlockEventsRequest_MsgType {
	if m != nil {
		return m.MsgType
	}
	return nil
}

func (m *ReportBlockEventsRequest) GetPreflight() *ReportBlockEventsRequest_Preflight {
	if x, ok := m.GetMsgType().(*ReportBlockEventsRequest_Preflight_); ok {
		return x.Preflight
	}
	return nil
}

func (m *ReportBlockEventsRequest) GetBlock() *protos.Block {
	if x, ok := m.GetMsgType().(*ReportBlockEventsRequest_Block); ok {
		return x.Block
	}
	return nil
}

func (m *ReportBlockEventsRequest) GetEvent() *protos.MonitorEvent {
	if x, ok := m.GetMsgType().(*ReportBlockEventsRequest_Event); ok {
		return x.Event
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ReportBlockEventsRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ReportBlockEventsRequest_OneofMarshaler, _ReportBlockEventsRequest_OneofUnmarshaler, _ReportBlockEventsRequest_OneofSizer, []interface{}{
		(*ReportBlockEventsRequest_Preflight_)(nil),
		(*ReportBlockEventsRequest_Block)(nil),
		(*ReportBlockEventsRequest_Event)(nil),
	}
}

func _ReportBlockEventsRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ReportBlockEventsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *ReportBlockEventsRequest_Preflight_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Preflight); err != nil {
			return err
		}
	case *ReportBlockEventsRequest_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *ReportBlockEventsRequest_Event:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Event); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ReportBlockEventsRequest.MsgType has unexpected type %T", x)
	}
	return nil
}

func _ReportBlockEventsRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ReportBlockEventsRequest)
	switch tag {
	case 1: // msg_type.preflight
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ReportBlockEventsRequest_Preflight)
		err := b.DecodeMessage(msg)
		m.MsgType = &ReportBlockEventsRequest_Preflight_{msg}
		return true, err
	case 2: // msg_type.block
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(protos.Block)
		err := b.DecodeMessage(msg)
		m.MsgType = &ReportBlockEventsRequest_Block{msg}
		return true, err
	case 3: // msg_type.event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(protos.MonitorEvent)
		err := b.DecodeMessage(msg)
		m.MsgType = &ReportBlockEventsRequest_Event{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ReportBlockEventsRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ReportBlockEventsRequest)
	// msg_type
	switch x := m.MsgType.(type) {
	case *ReportBlockEventsRequest_Preflight_:
		s := proto.Size(x.Preflight)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ReportBlockEventsRequest_Block:
		s := proto.Size(x.Block)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ReportBlockEventsRequest_Event:
		s := proto.Size(x.Event)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ReportBlockEventsRequest_Preflight struct {
	// LeaseId must be the id of a still valid Lease.
	LeaseId []byte `protobuf:"bytes,1,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	// MonitorSetVersion is the version of the MonitorSet this block was processed with.
	MonitorSetVersion    uint32   `protobuf:"varint,2,opt,name=monitor_set_version,json=monitorSetVersion,proto3" json:"monitor_set_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportBlockEventsRequest_Preflight) Reset()         { *m = ReportBlockEventsRequest_Preflight{} }
func (m *ReportBlockEventsRequest_Preflight) String() string { return proto.CompactTextString(m) }
func (*ReportBlockEventsRequest_Preflight) ProtoMessage()    {}
func (*ReportBlockEventsRequest_Preflight) Descriptor() ([]byte, []int) {
	return fileDescriptor_303630d4a056d2d3, []int{6, 0}
}

func (m *ReportBlockEventsRequest_Preflight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportBlockEventsRequest_Preflight.Unmarshal(m, b)
}
func (m *ReportBlockEventsRequest_Preflight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportBlockEventsRequest_Preflight.Marshal(b, m, deterministic)
}
func (m *ReportBlockEventsRequest_Preflight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportBlockEventsRequest_Preflight.Merge(m, src)
}
func (m *ReportBlockEventsRequest_Preflight) XXX_Size() int {
	return xxx_messageInfo_ReportBlockEventsRequest_Preflight.Size(m)
}
func (m *ReportBlockEventsRequest_Preflight) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportBlockEventsRequest_Preflight.DiscardUnknown(m)
}

var xxx_messageInfo_ReportBlockEventsRequest_Preflight proto.InternalMessageInfo

func (m *ReportBlockEventsRequest_Preflight) GetLeaseId() []byte {
	if m != nil {
		return m.LeaseId
	}
	return nil
}

func (m *ReportBlockEventsRequest_Preflight) GetMonitorSetVersion() uint32 {
	if m != nil {
		return m.MonitorSetVersion
	}
	return 0
}

func init() {
	proto.RegisterType((*RegisterChainManagerRequest)(nil), "grpc.RegisterChainManagerRequest")
	proto.RegisterType((*RegisterChainManagerResponse)(nil), "grpc.RegisterChainManagerResponse")
	proto.RegisterType((*UnregisterChainManagerRequest)(nil), "grpc.UnregisterChainManagerRequest")
	proto.RegisterType((*InstanceId)(nil), "grpc.InstanceId")
	proto.RegisterType((*RenewLeaseRequest)(nil), "grpc.RenewLeaseRequest")
	proto.RegisterType((*Lease)(nil), "grpc.Lease")
	proto.RegisterType((*ReportBlockEventsRequest)(nil), "grpc.ReportBlockEventsRequest")
	proto.RegisterType((*ReportBlockEventsRequest_Preflight)(nil), "grpc.ReportBlockEventsRequest.Preflight")
}

func init() { proto.RegisterFile("grpc/orchestrator.proto", fileDescriptor_303630d4a056d2d3) }

var fileDescriptor_303630d4a056d2d3 = []byte{
	// 671 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xd1, 0x4e, 0xdb, 0x4a,
	0x10, 0xc5, 0xb9, 0xe4, 0x92, 0x4c, 0x02, 0x97, 0xec, 0x45, 0x10, 0x42, 0x69, 0x83, 0xab, 0x4a,
	0x56, 0x8b, 0x1c, 0x29, 0x7d, 0xa9, 0xfa, 0x56, 0x5a, 0x24, 0x22, 0x15, 0x15, 0x2d, 0x82, 0x87,
	0x4a, 0x95, 0xe5, 0xd8, 0x83, 0x63, 0x35, 0xde, 0x75, 0x77, 0x37, 0x54, 0x7c, 0x4c, 0x9f, 0xfa,
	0x0b, 0xfd, 0x87, 0xfe, 0x56, 0xb5, 0xbb, 0xb6, 0x13, 0x52, 0xc8, 0xd3, 0x6a, 0x66, 0xce, 0xec,
	0x1c, 0x9d, 0x39, 0xbb, 0xb0, 0x97, 0x88, 0x3c, 0x1a, 0x70, 0x11, 0x4d, 0x50, 0x2a, 0x11, 0x2a,
	0x2e, 0xfc, 0x5c, 0x70, 0xc5, 0xc9, 0xba, 0x2e, 0xf4, 0x3a, 0xa6, 0x1c, 0xf1, 0x2c, 0xe3, 0xcc,
	0x16, 0x7a, 0xc4, 0x1c, 0x72, 0x30, 0x9e, 0xf2, 0xe8, 0xeb, 0x52, 0x0e, 0x6f, 0x91, 0xa9, 0x22,
	0x77, 0x90, 0x70, 0x9e, 0x4c, 0x71, 0x60, 0xa2, 0xf1, 0xec, 0x66, 0x80, 0x59, 0xae, 0xee, 0x6c,
	0xd1, 0xfd, 0xe5, 0xc0, 0x01, 0xc5, 0x24, 0x95, 0x0a, 0xc5, 0xfb, 0x49, 0x98, 0xb2, 0xf3, 0x90,
	0x85, 0x09, 0x0a, 0x8a, 0xdf, 0x66, 0x28, 0x15, 0x79, 0x0b, 0xdb, 0x91, 0x4e, 0x07, 0x99, 0xcd,
	0x07, 0x69, 0xdc, 0x75, 0xfa, 0x8e, 0xd7, 0x1a, 0x6e, 0xfb, 0x9a, 0x92, 0x3f, 0x62, 0x52, 0x85,
	0x2c, 0xc2, 0x51, 0x4c, 0xb7, 0xa2, 0x85, 0x0b, 0x46, 0x31, 0xd9, 0x87, 0x86, 0xed, 0x4d, 0xe3,
	0x6e, 0xad, 0xef, 0x78, 0x4d, 0xba, 0x61, 0xe2, 0x51, 0x4c, 0x9e, 0x41, 0x6b, 0xaa, 0x67, 0xb2,
	0x20, 0xe7, 0x42, 0x75, 0xff, 0xe9, 0x3b, 0x5e, 0x9d, 0x82, 0x4d, 0x5d, 0x70, 0xa1, 0xc8, 0x21,
	0x80, 0x44, 0x29, 0x53, 0x6e, 0xba, 0xd7, 0xfb, 0x8e, 0xd7, 0xa6, 0xcd, 0x22, 0x33, 0x8a, 0xdd,
	0x9f, 0x0e, 0x3c, 0x79, 0x98, 0xb6, 0xcc, 0x39, 0x93, 0x48, 0x8e, 0xa0, 0x3e, 0xc5, 0x50, 0x62,
	0x41, 0xb6, 0x65, 0xc9, 0x7e, 0xd4, 0x29, 0x6a, 0x2b, 0x64, 0x08, 0x6d, 0x81, 0x72, 0x96, 0x61,
	0x10, 0xde, 0x28, 0x14, 0x86, 0x62, 0x6b, 0xf8, 0x9f, 0x45, 0x9e, 0x68, 0x51, 0x2f, 0x73, 0x8c,
	0x68, 0xcb, 0x82, 0xde, 0x69, 0x0c, 0x39, 0x86, 0x46, 0xc6, 0x59, 0xaa, 0xb8, 0x90, 0x86, 0x74,
	0x25, 0xc3, 0xb9, 0xcd, 0x5e, 0xa2, 0xa2, 0x15, 0xc2, 0xe5, 0x70, 0x78, 0xc5, 0xc4, 0x0a, 0x75,
	0x17, 0x15, 0x72, 0xee, 0x2b, 0xb4, 0x0f, 0x0d, 0x43, 0xb3, 0x14, 0xaf, 0x4d, 0x37, 0x4c, 0x3c,
	0x8a, 0x49, 0x17, 0x36, 0x32, 0x94, 0x32, 0x4c, 0xd0, 0x70, 0x68, 0xd2, 0x32, 0x74, 0xdf, 0x00,
	0xcc, 0xf7, 0x41, 0xb6, 0xa0, 0x56, 0xdc, 0xdb, 0xa6, 0xb5, 0x34, 0x26, 0x3d, 0x68, 0xa4, 0x45,
	0xd5, 0x5c, 0xb9, 0x49, 0xab, 0xd8, 0xf5, 0xa1, 0x43, 0x91, 0xe1, 0x77, 0xab, 0xd0, 0x9c, 0x5e,
	0xc5, 0xc1, 0xb9, 0xc7, 0xc1, 0xfd, 0x00, 0x75, 0x03, 0xfd, 0x6b, 0xc8, 0x2b, 0xe8, 0x08, 0xcc,
	0xc2, 0x94, 0xa5, 0x2c, 0x09, 0x24, 0x46, 0x9c, 0xc5, 0xb2, 0x98, 0xb6, 0x5d, 0x15, 0x2e, 0x6d,
	0xde, 0xfd, 0x51, 0x83, 0x2e, 0x45, 0x6d, 0x01, 0xa3, 0xf7, 0xa9, 0x76, 0xad, 0x2c, 0xa7, 0x9f,
	0x41, 0x33, 0x17, 0x78, 0x33, 0x4d, 0x93, 0x89, 0x2a, 0xd6, 0xe8, 0x59, 0xb1, 0x1f, 0x6b, 0xf1,
	0x2f, 0x4a, 0xfc, 0xd9, 0x1a, 0x9d, 0x37, 0x93, 0x17, 0x50, 0x37, 0x8f, 0xa4, 0x58, 0xf1, 0xa6,
	0xf5, 0xbe, 0xb4, 0x4b, 0x3e, 0x5b, 0xa3, 0xb6, 0x4a, 0x8e, 0xa1, 0x6e, 0xde, 0x4d, 0xb1, 0xd9,
	0x9d, 0x12, 0x56, 0xec, 0xd6, 0x8c, 0xd2, 0x68, 0x03, 0xea, 0x5d, 0x43, 0xb3, 0x1a, 0xb7, 0x42,
	0x29, 0xe2, 0xc3, 0xff, 0x85, 0x21, 0x02, 0x89, 0x2a, 0xb8, 0x45, 0xa1, 0x3d, 0x5c, 0x48, 0xd2,
	0xc9, 0x2a, 0xdf, 0x5c, 0xdb, 0xc2, 0x09, 0x40, 0x23, 0x93, 0x49, 0xa0, 0xee, 0x72, 0x1c, 0xfe,
	0xae, 0x41, 0xfb, 0xd3, 0xc2, 0x97, 0x40, 0xbe, 0xc0, 0xce, 0x43, 0xb6, 0x27, 0x47, 0xa5, 0x30,
	0x8f, 0x7a, 0xad, 0xe7, 0xae, 0x82, 0x14, 0xaf, 0xe6, 0x0a, 0x76, 0x8d, 0x61, 0x65, 0xba, 0x3c,
	0xe0, 0xb9, 0xed, 0x5e, 0x69, 0xe7, 0xde, 0xae, 0x6f, 0xbf, 0x1a, 0xbf, 0xfc, 0x6a, 0xfc, 0x53,
	0xfd, 0xd5, 0x90, 0x21, 0xc0, 0xdc, 0x5c, 0x64, 0xaf, 0x24, 0xb2, 0x64, 0xb7, 0xde, 0xe2, 0x23,
	0x25, 0xe7, 0xda, 0x90, 0x4b, 0x6b, 0x26, 0x4f, 0x57, 0xef, 0xff, 0x31, 0x02, 0x9e, 0x73, 0xf2,
	0xf2, 0xb3, 0x97, 0xa4, 0x6a, 0x32, 0x1b, 0xfb, 0x11, 0xcf, 0x06, 0x19, 0x26, 0xa1, 0xcc, 0xc3,
	0x08, 0xa7, 0xe1, 0xd8, 0x04, 0x11, 0x67, 0x0c, 0x23, 0x35, 0xd0, 0xd7, 0x8f, 0xff, 0x35, 0xdd,
	0xaf, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x86, 0x1b, 0xd0, 0x00, 0x93, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrchestratorClient is the client API for Orchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrchestratorClient interface {
	// RegisterChainManager is invoked by ChainManagers to register themselves with this Orchestrator.
	RegisterChainManager(ctx context.Context, in *RegisterChainManagerRequest, opts ...grpc.CallOption) (*RegisterChainManagerResponse, error)
	// UnregisterChainManager is invoked by ChainManager to gracefully unregister themselves with this Orchestrator.
	// This usually happens when problems are noticed on its connector and ChainManager needs to terminate itself.
	UnregsiterChainManager(ctx context.Context, in *UnregisterChainManagerRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// RenewLease renews the lease between a ChainManager and this Orchestrator.
	// Upon lease expiration, Orchestrator will consider the ChainManager dead and be free to accept new registrations.
	RenewLease(ctx context.Context, in *RenewLeaseRequest, opts ...grpc.CallOption) (*Lease, error)
	// ReportBlockEvents reports a new block on the monitored chain, as well as all fired events from this block.
	// The request stream must begin with a valid Preflight, followed by a Block, followed by zero or more Events.
	ReportBlockEvents(ctx context.Context, opts ...grpc.CallOption) (Orchestrator_ReportBlockEventsClient, error)
}

type orchestratorClient struct {
	cc *grpc.ClientConn
}

func NewOrchestratorClient(cc *grpc.ClientConn) OrchestratorClient {
	return &orchestratorClient{cc}
}

func (c *orchestratorClient) RegisterChainManager(ctx context.Context, in *RegisterChainManagerRequest, opts ...grpc.CallOption) (*RegisterChainManagerResponse, error) {
	out := new(RegisterChainManagerResponse)
	err := c.cc.Invoke(ctx, "/grpc.Orchestrator/RegisterChainManager", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) UnregsiterChainManager(ctx context.Context, in *UnregisterChainManagerRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/grpc.Orchestrator/UnregsiterChainManager", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) RenewLease(ctx context.Context, in *RenewLeaseRequest, opts ...grpc.CallOption) (*Lease, error) {
	out := new(Lease)
	err := c.cc.Invoke(ctx, "/grpc.Orchestrator/RenewLease", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) ReportBlockEvents(ctx context.Context, opts ...grpc.CallOption) (Orchestrator_ReportBlockEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Orchestrator_serviceDesc.Streams[0], "/grpc.Orchestrator/ReportBlockEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &orchestratorReportBlockEventsClient{stream}
	return x, nil
}

type Orchestrator_ReportBlockEventsClient interface {
	Send(*ReportBlockEventsRequest) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type orchestratorReportBlockEventsClient struct {
	grpc.ClientStream
}

func (x *orchestratorReportBlockEventsClient) Send(m *ReportBlockEventsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *orchestratorReportBlockEventsClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrchestratorServer is the server API for Orchestrator service.
type OrchestratorServer interface {
	// RegisterChainManager is invoked by ChainManagers to register themselves with this Orchestrator.
	RegisterChainManager(context.Context, *RegisterChainManagerRequest) (*RegisterChainManagerResponse, error)
	// UnregisterChainManager is invoked by ChainManager to gracefully unregister themselves with this Orchestrator.
	// This usually happens when problems are noticed on its connector and ChainManager needs to terminate itself.
	UnregsiterChainManager(context.Context, *UnregisterChainManagerRequest) (*empty.Empty, error)
	// RenewLease renews the lease between a ChainManager and this Orchestrator.
	// Upon lease expiration, Orchestrator will consider the ChainManager dead and be free to accept new registrations.
	RenewLease(context.Context, *RenewLeaseRequest) (*Lease, error)
	// ReportBlockEvents reports a new block on the monitored chain, as well as all fired events from this block.
	// The request stream must begin with a valid Preflight, followed by a Block, followed by zero or more Events.
	ReportBlockEvents(Orchestrator_ReportBlockEventsServer) error
}

func RegisterOrchestratorServer(s *grpc.Server, srv OrchestratorServer) {
	s.RegisterService(&_Orchestrator_serviceDesc, srv)
}

func _Orchestrator_RegisterChainManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterChainManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).RegisterChainManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Orchestrator/RegisterChainManager",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).RegisterChainManager(ctx, req.(*RegisterChainManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_UnregsiterChainManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnregisterChainManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).UnregsiterChainManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Orchestrator/UnregsiterChainManager",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).UnregsiterChainManager(ctx, req.(*UnregisterChainManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_RenewLease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewLeaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).RenewLease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Orchestrator/RenewLease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).RenewLease(ctx, req.(*RenewLeaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_ReportBlockEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrchestratorServer).ReportBlockEvents(&orchestratorReportBlockEventsServer{stream})
}

type Orchestrator_ReportBlockEventsServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*ReportBlockEventsRequest, error)
	grpc.ServerStream
}

type orchestratorReportBlockEventsServer struct {
	grpc.ServerStream
}

func (x *orchestratorReportBlockEventsServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *orchestratorReportBlockEventsServer) Recv() (*ReportBlockEventsRequest, error) {
	m := new(ReportBlockEventsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Orchestrator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Orchestrator",
	HandlerType: (*OrchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterChainManager",
			Handler:    _Orchestrator_RegisterChainManager_Handler,
		},
		{
			MethodName: "UnregsiterChainManager",
			Handler:    _Orchestrator_UnregsiterChainManager_Handler,
		},
		{
			MethodName: "RenewLease",
			Handler:    _Orchestrator_RenewLease_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReportBlockEvents",
			Handler:       _Orchestrator_ReportBlockEvents_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "grpc/orchestrator.proto",
}
