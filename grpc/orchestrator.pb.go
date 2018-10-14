// Code generated by protoc-gen-go. DO NOT EDIT.
// source: orchestrator.proto

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
	return fileDescriptor_96b6e6782baaa298, []int{0}
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
	// ResumeAfterBlockHash indicates the last processed block by Orchestrator.
	// If not nil, ChainManager is expected to start reporting from right after this block.
	// Otherwise, ChainManager can start reporting from the next new block.
	ResumeAfterBlockHash []byte `protobuf:"bytes,2,opt,name=resume_after_block_hash,json=resumeAfterBlockHash,proto3" json:"resume_after_block_hash,omitempty"`
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
	return fileDescriptor_96b6e6782baaa298, []int{1}
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

func (m *RegisterChainManagerResponse) GetResumeAfterBlockHash() []byte {
	if m != nil {
		return m.ResumeAfterBlockHash
	}
	return nil
}

func (m *RegisterChainManagerResponse) GetMonitors() *MonitorSet {
	if m != nil {
		return m.Monitors
	}
	return nil
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
	return fileDescriptor_96b6e6782baaa298, []int{2}
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

type MonitorSet struct {
	// Monitors contain all things to be monitored by the requesting ChainManager.
	Monitors []*Monitor `protobuf:"bytes,1,rep,name=monitors,proto3" json:"monitors,omitempty"`
	// Version is used for Orchestrator and ChainManager to agree on the set of monitors evaluated on each block.
	Version              uint32   `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MonitorSet) Reset()         { *m = MonitorSet{} }
func (m *MonitorSet) String() string { return proto.CompactTextString(m) }
func (*MonitorSet) ProtoMessage()    {}
func (*MonitorSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b6e6782baaa298, []int{3}
}

func (m *MonitorSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorSet.Unmarshal(m, b)
}
func (m *MonitorSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorSet.Marshal(b, m, deterministic)
}
func (m *MonitorSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorSet.Merge(m, src)
}
func (m *MonitorSet) XXX_Size() int {
	return xxx_messageInfo_MonitorSet.Size(m)
}
func (m *MonitorSet) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorSet.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorSet proto.InternalMessageInfo

func (m *MonitorSet) GetMonitors() []*Monitor {
	if m != nil {
		return m.Monitors
	}
	return nil
}

func (m *MonitorSet) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

// Monitor contains a list of encoded workflow expressions and a triggering condition
type Monitor struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// evaluations is a list of encoded workflow expressions, which will be evaluated if a monitor is triggered
	// the result will be passed back to flow manager as part of event. These expression results are intended for
	// evaluating actions
	Monitor              []byte   `protobuf:"bytes,2,opt,name=monitor,proto3" json:"monitor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Monitor) Reset()         { *m = Monitor{} }
func (m *Monitor) String() string { return proto.CompactTextString(m) }
func (*Monitor) ProtoMessage()    {}
func (*Monitor) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b6e6782baaa298, []int{4}
}

func (m *Monitor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Monitor.Unmarshal(m, b)
}
func (m *Monitor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Monitor.Marshal(b, m, deterministic)
}
func (m *Monitor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Monitor.Merge(m, src)
}
func (m *Monitor) XXX_Size() int {
	return xxx_messageInfo_Monitor.Size(m)
}
func (m *Monitor) XXX_DiscardUnknown() {
	xxx_messageInfo_Monitor.DiscardUnknown(m)
}

var xxx_messageInfo_Monitor proto.InternalMessageInfo

func (m *Monitor) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Monitor) GetMonitor() []byte {
	if m != nil {
		return m.Monitor
	}
	return nil
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
	return fileDescriptor_96b6e6782baaa298, []int{5}
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
	return fileDescriptor_96b6e6782baaa298, []int{6}
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
	return fileDescriptor_96b6e6782baaa298, []int{7}
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
	Block *Block `protobuf:"bytes,2,opt,name=block,proto3,oneof"`
}

type ReportBlockEventsRequest_Event struct {
	Event *Event `protobuf:"bytes,3,opt,name=event,proto3,oneof"`
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

func (m *ReportBlockEventsRequest) GetBlock() *Block {
	if x, ok := m.GetMsgType().(*ReportBlockEventsRequest_Block); ok {
		return x.Block
	}
	return nil
}

func (m *ReportBlockEventsRequest) GetEvent() *Event {
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
		msg := new(Block)
		err := b.DecodeMessage(msg)
		m.MsgType = &ReportBlockEventsRequest_Block{msg}
		return true, err
	case 3: // msg_type.event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Event)
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
	return fileDescriptor_96b6e6782baaa298, []int{7, 0}
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

// TODO - TBD
type Block struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	ParentHash           []byte   `protobuf:"bytes,2,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	Height               *BigInt  `protobuf:"bytes,3,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b6e6782baaa298, []int{8}
}

func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Block) GetParentHash() []byte {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *Block) GetHeight() *BigInt {
	if m != nil {
		return m.Height
	}
	return nil
}

type BigInt struct {
	// Bytes are interpreted as a big-endian unsigned integer.
	Bytes                []byte   `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
	Negative             bool     `protobuf:"varint,2,opt,name=negative,proto3" json:"negative,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BigInt) Reset()         { *m = BigInt{} }
func (m *BigInt) String() string { return proto.CompactTextString(m) }
func (*BigInt) ProtoMessage()    {}
func (*BigInt) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b6e6782baaa298, []int{9}
}

func (m *BigInt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BigInt.Unmarshal(m, b)
}
func (m *BigInt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BigInt.Marshal(b, m, deterministic)
}
func (m *BigInt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BigInt.Merge(m, src)
}
func (m *BigInt) XXX_Size() int {
	return xxx_messageInfo_BigInt.Size(m)
}
func (m *BigInt) XXX_DiscardUnknown() {
	xxx_messageInfo_BigInt.DiscardUnknown(m)
}

var xxx_messageInfo_BigInt proto.InternalMessageInfo

func (m *BigInt) GetBytes() []byte {
	if m != nil {
		return m.Bytes
	}
	return nil
}

func (m *BigInt) GetNegative() bool {
	if m != nil {
		return m.Negative
	}
	return false
}

// Event describes a fired Monitor.
type Event struct {
	// MonitorId identifies the Monitor being fired.
	MonitorId int64 `protobuf:"varint,1,opt,name=monitor_id,json=monitorId,proto3" json:"monitor_id,omitempty"`
	// Context contains a list of evalutaion result of expressions in Monitor, the order of list is
	// the same the order of list in Monitor
	EvaluationsResults   [][]byte `protobuf:"bytes,2,rep,name=evaluations_results,json=evaluationsResults,proto3" json:"evaluations_results,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b6e6782baaa298, []int{10}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetMonitorId() int64 {
	if m != nil {
		return m.MonitorId
	}
	return 0
}

func (m *Event) GetEvaluationsResults() [][]byte {
	if m != nil {
		return m.EvaluationsResults
	}
	return nil
}

func init() {
	proto.RegisterType((*RegisterChainManagerRequest)(nil), "grpc.RegisterChainManagerRequest")
	proto.RegisterType((*RegisterChainManagerResponse)(nil), "grpc.RegisterChainManagerResponse")
	proto.RegisterType((*InstanceId)(nil), "grpc.InstanceId")
	proto.RegisterType((*MonitorSet)(nil), "grpc.MonitorSet")
	proto.RegisterType((*Monitor)(nil), "grpc.Monitor")
	proto.RegisterType((*RenewLeaseRequest)(nil), "grpc.RenewLeaseRequest")
	proto.RegisterType((*Lease)(nil), "grpc.Lease")
	proto.RegisterType((*ReportBlockEventsRequest)(nil), "grpc.ReportBlockEventsRequest")
	proto.RegisterType((*ReportBlockEventsRequest_Preflight)(nil), "grpc.ReportBlockEventsRequest.Preflight")
	proto.RegisterType((*Block)(nil), "grpc.Block")
	proto.RegisterType((*BigInt)(nil), "grpc.BigInt")
	proto.RegisterType((*Event)(nil), "grpc.Event")
}

func init() { proto.RegisterFile("orchestrator.proto", fileDescriptor_96b6e6782baaa298) }

var fileDescriptor_96b6e6782baaa298 = []byte{
	// 721 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xdd, 0x4e, 0xdb, 0x4a,
	0x10, 0xc6, 0x09, 0x21, 0xc9, 0x24, 0x20, 0x58, 0xd0, 0x21, 0x84, 0xf3, 0x13, 0x7c, 0xce, 0x85,
	0x8f, 0xce, 0x91, 0x91, 0x82, 0x2a, 0x55, 0xdc, 0x95, 0x16, 0x09, 0x4b, 0x45, 0xa5, 0x8b, 0x44,
	0xaf, 0x2a, 0xcb, 0x89, 0x07, 0xc7, 0x6a, 0xb2, 0xeb, 0xee, 0x6e, 0x52, 0xf1, 0x1a, 0x7d, 0x88,
	0x3e, 0x41, 0x1f, 0xac, 0x8f, 0x50, 0xed, 0x8f, 0x9d, 0x40, 0x81, 0x3b, 0xcf, 0xf7, 0x8d, 0x67,
	0x66, 0xbf, 0xfd, 0x76, 0x80, 0x70, 0x31, 0x9e, 0xa0, 0x54, 0x22, 0x51, 0x5c, 0x84, 0x85, 0xe0,
	0x8a, 0x93, 0xf5, 0x4c, 0x14, 0xe3, 0xfe, 0x61, 0xc6, 0x79, 0x36, 0xc5, 0x63, 0x83, 0x8d, 0xe6,
	0xb7, 0xc7, 0x38, 0x2b, 0xd4, 0x9d, 0x4d, 0xf1, 0xbf, 0x7b, 0x70, 0x48, 0x31, 0xcb, 0xa5, 0x42,
	0xf1, 0x7a, 0x92, 0xe4, 0xec, 0x32, 0x61, 0x49, 0x86, 0x82, 0xe2, 0xe7, 0x39, 0x4a, 0x45, 0x4e,
	0x61, 0x7b, 0xac, 0xe1, 0x78, 0x66, 0xf1, 0x38, 0x4f, 0x7b, 0xde, 0xc0, 0x0b, 0x3a, 0xc3, 0xed,
	0x50, 0x57, 0x0f, 0x23, 0x26, 0x55, 0xc2, 0xc6, 0x18, 0xa5, 0x74, 0x6b, 0xbc, 0x52, 0x20, 0x4a,
	0xc9, 0x01, 0xb4, 0xec, 0xbf, 0x79, 0xda, 0xab, 0x0d, 0xbc, 0xa0, 0x4d, 0x9b, 0x26, 0x8e, 0x52,
	0xf2, 0x17, 0x74, 0xa6, 0xba, 0x27, 0x8b, 0x0b, 0x2e, 0x54, 0xaf, 0x3e, 0xf0, 0x82, 0x06, 0x05,
	0x0b, 0x5d, 0x71, 0xa1, 0xc8, 0x1f, 0x00, 0x12, 0xa5, 0xcc, 0xb9, 0xf9, 0x7b, 0x7d, 0xe0, 0x05,
	0x5d, 0xda, 0x76, 0x48, 0x94, 0xfa, 0xdf, 0x3c, 0xf8, 0xfd, 0xf1, 0xb1, 0x65, 0xc1, 0x99, 0x44,
	0x72, 0x04, 0x8d, 0x29, 0x26, 0x12, 0xdd, 0xb0, 0x1d, 0x3b, 0xec, 0x5b, 0x0d, 0x51, 0xcb, 0x90,
	0x17, 0xb0, 0x2f, 0x50, 0xce, 0x67, 0x18, 0x27, 0xb7, 0x0a, 0x45, 0x3c, 0x9a, 0xf2, 0xf1, 0xa7,
	0x78, 0x92, 0xc8, 0x89, 0x99, 0xb6, 0x4b, 0xf7, 0x2c, 0xfd, 0x4a, 0xb3, 0x67, 0x9a, 0xbc, 0x48,
	0xe4, 0x84, 0xfc, 0x0f, 0xad, 0x19, 0x67, 0xb9, 0xe2, 0x42, 0x9a, 0xb9, 0x2b, 0x25, 0x2e, 0x2d,
	0x7a, 0x8d, 0x8a, 0x56, 0x19, 0xfe, 0x4b, 0x80, 0xa5, 0x42, 0x64, 0x0b, 0x6a, 0x4e, 0xbf, 0x2e,
	0xad, 0xe5, 0x29, 0xe9, 0x43, 0x2b, 0x77, 0xac, 0xe9, 0xb9, 0x49, 0xab, 0xd8, 0x7f, 0x0f, 0xb0,
	0xac, 0x48, 0xfe, 0x5d, 0xe9, 0xea, 0x0d, 0xea, 0x41, 0x67, 0xb8, 0x79, 0xaf, 0xeb, 0xb2, 0x25,
	0xe9, 0x41, 0x73, 0x81, 0x42, 0x0b, 0xe5, 0x6a, 0x96, 0xa1, 0x7f, 0x02, 0x4d, 0x97, 0xbe, 0x32,
	0x49, 0xdd, 0x4c, 0xd2, 0x83, 0xa6, 0x2b, 0xe0, 0x0e, 0x5f, 0x86, 0x7e, 0x08, 0x3b, 0x14, 0x19,
	0x7e, 0xb1, 0xda, 0x39, 0x5b, 0x1c, 0x40, 0xcb, 0x88, 0x18, 0x57, 0xc7, 0x69, 0x9a, 0x38, 0x4a,
	0xfd, 0x37, 0xd0, 0x30, 0xa9, 0xbf, 0x1c, 0xf6, 0x3f, 0xd8, 0x11, 0x38, 0x4b, 0x72, 0x96, 0xb3,
	0x2c, 0x96, 0x38, 0xe6, 0x2c, 0x95, 0x6e, 0xc2, 0xed, 0x8a, 0xb8, 0xb6, 0xb8, 0xff, 0xb5, 0x06,
	0x3d, 0x8a, 0xda, 0x1c, 0x46, 0xf9, 0xf3, 0x05, 0x32, 0x25, 0xcb, 0xee, 0x17, 0xd0, 0x2e, 0x04,
	0xde, 0x4e, 0xf3, 0x6c, 0xa2, 0xdc, 0x05, 0x07, 0x56, 0x8d, 0xa7, 0x7e, 0x09, 0xaf, 0xca, 0xfc,
	0x8b, 0x35, 0xba, 0xfc, 0x99, 0xfc, 0x0d, 0x0d, 0x73, 0xed, 0x66, 0x8e, 0xca, 0x26, 0xf6, 0xb2,
	0xd7, 0xa8, 0xe5, 0x74, 0x12, 0xea, 0x62, 0xee, 0xba, 0x5d, 0x92, 0xa9, 0xaf, 0x93, 0x0c, 0xd7,
	0xbf, 0x81, 0x76, 0xd5, 0xe3, 0x19, 0x79, 0x48, 0x08, 0xbb, 0x4e, 0xd9, 0x58, 0xa2, 0x8a, 0xef,
	0xdf, 0xd4, 0xce, 0xac, 0xba, 0xf1, 0x1b, 0x4b, 0x9c, 0x01, 0xb4, 0x66, 0x32, 0x8b, 0xd5, 0x5d,
	0x81, 0xfe, 0x08, 0x1a, 0x66, 0x34, 0x42, 0x60, 0xdd, 0xf8, 0xd4, 0xd6, 0x36, 0xdf, 0xfa, 0x49,
	0x15, 0x89, 0x40, 0xa6, 0x56, 0x2d, 0x0c, 0x16, 0x32, 0xc6, 0xfd, 0x07, 0x36, 0x26, 0x68, 0x24,
	0xb3, 0xe7, 0xe8, 0xba, 0xc3, 0xe6, 0x59, 0xc4, 0x14, 0x75, 0x9c, 0x7f, 0x0a, 0x1b, 0x16, 0x21,
	0x7b, 0xd0, 0x18, 0xdd, 0x29, 0x94, 0xae, 0x8b, 0x0d, 0xb4, 0x65, 0x19, 0x66, 0x89, 0xca, 0x17,
	0xd6, 0xb2, 0x2d, 0x5a, 0xc5, 0xfe, 0x07, 0x68, 0x18, 0x55, 0xf4, 0xeb, 0x2d, 0x0f, 0x59, 0xb9,
	0xac, 0xed, 0x90, 0x28, 0x25, 0xc7, 0xb0, 0x8b, 0x8b, 0x64, 0x3a, 0x4f, 0x54, 0xce, 0x99, 0x8c,
	0xf5, 0x33, 0x9b, 0x2a, 0xed, 0x85, 0x7a, 0xd0, 0xa5, 0x64, 0x85, 0xa2, 0x96, 0x19, 0xfe, 0xf0,
	0xa0, 0xfb, 0x6e, 0x65, 0xbf, 0x91, 0x8f, 0xb0, 0xf7, 0xd8, 0xf3, 0x27, 0x47, 0xa5, 0x0d, 0x9e,
	0xdc, 0x68, 0x7d, 0xff, 0xb9, 0x14, 0xb7, 0x3d, 0x86, 0x00, 0x4b, 0xcf, 0x93, 0xfd, 0xf2, 0x8f,
	0x07, 0xaf, 0xa0, 0xbf, 0xba, 0x55, 0xc8, 0xa5, 0x7e, 0x27, 0x0f, 0xdc, 0x47, 0xfe, 0x7c, 0xde,
	0x96, 0xfd, 0xdf, 0x42, 0xbb, 0x9c, 0xc3, 0x72, 0x39, 0x87, 0xe7, 0x7a, 0x39, 0x07, 0xde, 0x68,
	0xc3, 0x20, 0x27, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x1e, 0xba, 0x00, 0xd8, 0x05, 0x00,
	0x00,
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
	Metadata: "orchestrator.proto",
}
