// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

// 指定等会文件生成出来的package

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//VoteRequest 投票的请求
type VoteRequest struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId          string   `protobuf:"bytes,2,opt,name=candidate_id,json=candidateId,proto3" json:"candidate_id,omitempty"`
	LastLogIndex         int64    `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	LastLogTerm          int64    `protobuf:"varint,4,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteRequest) Reset()         { *m = VoteRequest{} }
func (m *VoteRequest) String() string { return proto.CompactTextString(m) }
func (*VoteRequest) ProtoMessage()    {}
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *VoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteRequest.Unmarshal(m, b)
}
func (m *VoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteRequest.Marshal(b, m, deterministic)
}
func (m *VoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteRequest.Merge(m, src)
}
func (m *VoteRequest) XXX_Size() int {
	return xxx_messageInfo_VoteRequest.Size(m)
}
func (m *VoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VoteRequest proto.InternalMessageInfo

func (m *VoteRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *VoteRequest) GetCandidateId() string {
	if m != nil {
		return m.CandidateId
	}
	return ""
}

func (m *VoteRequest) GetLastLogIndex() int64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *VoteRequest) GetLastLogTerm() int64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

//VoteResponse 投票的响应
type VoteResponse struct {
	Term                 int64    `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted          bool     `protobuf:"varint,3,opt,name=vote_granted,json=voteGranted,proto3" json:"vote_granted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteResponse) Reset()         { *m = VoteResponse{} }
func (m *VoteResponse) String() string { return proto.CompactTextString(m) }
func (*VoteResponse) ProtoMessage()    {}
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *VoteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteResponse.Unmarshal(m, b)
}
func (m *VoteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteResponse.Marshal(b, m, deterministic)
}
func (m *VoteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteResponse.Merge(m, src)
}
func (m *VoteResponse) XXX_Size() int {
	return xxx_messageInfo_VoteResponse.Size(m)
}
func (m *VoteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VoteResponse proto.InternalMessageInfo

func (m *VoteResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *VoteResponse) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

type AppendEntriesRequest struct {
	Term                 int64       `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId             string      `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	PreLogIndex          int64       `protobuf:"varint,3,opt,name=pre_log_index,json=preLogIndex,proto3" json:"pre_log_index,omitempty"`
	PreLogTerm           int64       `protobuf:"varint,4,opt,name=pre_log_term,json=preLogTerm,proto3" json:"pre_log_term,omitempty"`
	LeaderCommit         int64       `protobuf:"varint,5,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
	Entry                []*LogEntry `protobuf:"bytes,6,rep,name=entry,proto3" json:"entry,omitempty"`
	IsApply              bool        `protobuf:"varint,7,opt,name=is_apply,json=isApply,proto3" json:"is_apply,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *AppendEntriesRequest) Reset()         { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()    {}
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *AppendEntriesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesRequest.Unmarshal(m, b)
}
func (m *AppendEntriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesRequest.Marshal(b, m, deterministic)
}
func (m *AppendEntriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesRequest.Merge(m, src)
}
func (m *AppendEntriesRequest) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesRequest.Size(m)
}
func (m *AppendEntriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesRequest proto.InternalMessageInfo

func (m *AppendEntriesRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}

func (m *AppendEntriesRequest) GetPreLogIndex() int64 {
	if m != nil {
		return m.PreLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetPreLogTerm() int64 {
	if m != nil {
		return m.PreLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderCommit() int64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

func (m *AppendEntriesRequest) GetEntry() []*LogEntry {
	if m != nil {
		return m.Entry
	}
	return nil
}

func (m *AppendEntriesRequest) GetIsApply() bool {
	if m != nil {
		return m.IsApply
	}
	return false
}

type HeartbeatRequest struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId             string   `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	LeaderCommit         int64    `protobuf:"varint,3,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatRequest) Reset()         { *m = HeartbeatRequest{} }
func (m *HeartbeatRequest) String() string { return proto.CompactTextString(m) }
func (*HeartbeatRequest) ProtoMessage()    {}
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{3}
}

func (m *HeartbeatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatRequest.Unmarshal(m, b)
}
func (m *HeartbeatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatRequest.Marshal(b, m, deterministic)
}
func (m *HeartbeatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatRequest.Merge(m, src)
}
func (m *HeartbeatRequest) XXX_Size() int {
	return xxx_messageInfo_HeartbeatRequest.Size(m)
}
func (m *HeartbeatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatRequest proto.InternalMessageInfo

func (m *HeartbeatRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *HeartbeatRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}

func (m *HeartbeatRequest) GetLeaderCommit() int64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type HeartbeatResponse struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success              bool     `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatResponse) Reset()         { *m = HeartbeatResponse{} }
func (m *HeartbeatResponse) String() string { return proto.CompactTextString(m) }
func (*HeartbeatResponse) ProtoMessage()    {}
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{4}
}

func (m *HeartbeatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatResponse.Unmarshal(m, b)
}
func (m *HeartbeatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatResponse.Marshal(b, m, deterministic)
}
func (m *HeartbeatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatResponse.Merge(m, src)
}
func (m *HeartbeatResponse) XXX_Size() int {
	return xxx_messageInfo_HeartbeatResponse.Size(m)
}
func (m *HeartbeatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatResponse proto.InternalMessageInfo

func (m *HeartbeatResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *HeartbeatResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type InstallSnapshotRequest struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId             string   `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	LastIncludedIndex    int64    `protobuf:"varint,3,opt,name=last_included_index,json=lastIncludedIndex,proto3" json:"last_included_index,omitempty"`
	LastIncludedTerm     int64    `protobuf:"varint,4,opt,name=last_included_term,json=lastIncludedTerm,proto3" json:"last_included_term,omitempty"`
	Offset               int64    `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
	Data                 []byte   `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	Done                 bool     `protobuf:"varint,7,opt,name=done,proto3" json:"done,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotRequest) Reset()         { *m = InstallSnapshotRequest{} }
func (m *InstallSnapshotRequest) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotRequest) ProtoMessage()    {}
func (*InstallSnapshotRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{5}
}

func (m *InstallSnapshotRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotRequest.Unmarshal(m, b)
}
func (m *InstallSnapshotRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotRequest.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotRequest.Merge(m, src)
}
func (m *InstallSnapshotRequest) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotRequest.Size(m)
}
func (m *InstallSnapshotRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotRequest proto.InternalMessageInfo

func (m *InstallSnapshotRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}

func (m *InstallSnapshotRequest) GetLastIncludedIndex() int64 {
	if m != nil {
		return m.LastIncludedIndex
	}
	return 0
}

func (m *InstallSnapshotRequest) GetLastIncludedTerm() int64 {
	if m != nil {
		return m.LastIncludedTerm
	}
	return 0
}

func (m *InstallSnapshotRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *InstallSnapshotRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *InstallSnapshotRequest) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

type InstallSnapshotResponse struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotResponse) Reset()         { *m = InstallSnapshotResponse{} }
func (m *InstallSnapshotResponse) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotResponse) ProtoMessage()    {}
func (*InstallSnapshotResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{6}
}

func (m *InstallSnapshotResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotResponse.Unmarshal(m, b)
}
func (m *InstallSnapshotResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotResponse.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotResponse.Merge(m, src)
}
func (m *InstallSnapshotResponse) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotResponse.Size(m)
}
func (m *InstallSnapshotResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotResponse proto.InternalMessageInfo

func (m *InstallSnapshotResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

type AppendEntriesResponse struct {
	Term                 int64    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success              bool     `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppendEntriesResponse) Reset()         { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()    {}
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{7}
}

func (m *AppendEntriesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesResponse.Unmarshal(m, b)
}
func (m *AppendEntriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesResponse.Marshal(b, m, deterministic)
}
func (m *AppendEntriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesResponse.Merge(m, src)
}
func (m *AppendEntriesResponse) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesResponse.Size(m)
}
func (m *AppendEntriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesResponse proto.InternalMessageInfo

func (m *AppendEntriesResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

// LogEntry 日志的对象
type LogEntry struct {
	CurrentTerm          int64    `protobuf:"varint,1,opt,name=currentTerm,proto3" json:"currentTerm,omitempty"`
	Index                int64    `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Position             int64    `protobuf:"varint,4,opt,name=position,proto3" json:"position,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (*LogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{8}
}

func (m *LogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntry.Unmarshal(m, b)
}
func (m *LogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntry.Marshal(b, m, deterministic)
}
func (m *LogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntry.Merge(m, src)
}
func (m *LogEntry) XXX_Size() int {
	return xxx_messageInfo_LogEntry.Size(m)
}
func (m *LogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntry proto.InternalMessageInfo

func (m *LogEntry) GetCurrentTerm() int64 {
	if m != nil {
		return m.CurrentTerm
	}
	return 0
}

func (m *LogEntry) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *LogEntry) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *LogEntry) GetPosition() int64 {
	if m != nil {
		return m.Position
	}
	return 0
}

func init() {
	proto.RegisterType((*VoteRequest)(nil), "proto.VoteRequest")
	proto.RegisterType((*VoteResponse)(nil), "proto.VoteResponse")
	proto.RegisterType((*AppendEntriesRequest)(nil), "proto.AppendEntriesRequest")
	proto.RegisterType((*HeartbeatRequest)(nil), "proto.HeartbeatRequest")
	proto.RegisterType((*HeartbeatResponse)(nil), "proto.HeartbeatResponse")
	proto.RegisterType((*InstallSnapshotRequest)(nil), "proto.InstallSnapshotRequest")
	proto.RegisterType((*InstallSnapshotResponse)(nil), "proto.InstallSnapshotResponse")
	proto.RegisterType((*AppendEntriesResponse)(nil), "proto.AppendEntriesResponse")
	proto.RegisterType((*LogEntry)(nil), "proto.LogEntry")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 580 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdf, 0x6e, 0xd3, 0x3e,
	0x14, 0x56, 0xd6, 0xb5, 0xcb, 0x4e, 0xda, 0xdf, 0x36, 0x6f, 0xbf, 0x2d, 0x74, 0x80, 0x42, 0x00,
	0xa9, 0x17, 0x30, 0xa4, 0x71, 0x8f, 0x54, 0xa1, 0x0a, 0x8a, 0x76, 0x81, 0xb2, 0x89, 0xdb, 0xca,
	0x8b, 0x4f, 0xbb, 0x48, 0xa9, 0x1d, 0x6c, 0x77, 0xa2, 0x2f, 0xc1, 0x13, 0xf0, 0x7e, 0xbc, 0x00,
	0x0f, 0x80, 0x62, 0xbb, 0x59, 0xfa, 0x67, 0x5c, 0xf4, 0x2a, 0x3e, 0xdf, 0xf9, 0xe2, 0x73, 0xbe,
	0xf3, 0xc7, 0xd0, 0x51, 0x28, 0xef, 0xb3, 0x14, 0x2f, 0x0a, 0x29, 0xb4, 0x20, 0x4d, 0xf3, 0x89,
	0x7f, 0x7a, 0x10, 0x7c, 0x13, 0x1a, 0x13, 0xfc, 0x3e, 0x43, 0xa5, 0x09, 0x81, 0x5d, 0x8d, 0x72,
	0x1a, 0x7a, 0x91, 0xd7, 0x6b, 0x24, 0xe6, 0x4c, 0x5e, 0x40, 0x3b, 0xa5, 0x9c, 0x65, 0x8c, 0x6a,
	0x1c, 0x65, 0x2c, 0xdc, 0x89, 0xbc, 0xde, 0x7e, 0x12, 0x54, 0xd8, 0x90, 0x91, 0x57, 0xf0, 0x5f,
	0x4e, 0x95, 0x1e, 0xe5, 0x62, 0x32, 0xca, 0x38, 0xc3, 0x1f, 0x61, 0xc3, 0x5c, 0xd0, 0x2e, 0xd1,
	0x2b, 0x31, 0x19, 0x96, 0x18, 0x89, 0xa1, 0x53, 0xb1, 0x4c, 0x94, 0x5d, 0x43, 0x0a, 0x1c, 0xe9,
	0x06, 0xe5, 0x34, 0x1e, 0x40, 0xdb, 0xe6, 0xa3, 0x0a, 0xc1, 0x15, 0x56, 0x09, 0xed, 0x2c, 0x27,
	0x74, 0x2f, 0x34, 0x8e, 0x26, 0x92, 0x72, 0x8d, 0xcc, 0xc4, 0xf2, 0x93, 0xa0, 0xc4, 0x3e, 0x59,
	0x28, 0xfe, 0xe3, 0xc1, 0x49, 0xbf, 0x28, 0x90, 0xb3, 0x01, 0xd7, 0x32, 0x43, 0xf5, 0x2f, 0x81,
	0xe7, 0xb0, 0x9f, 0x23, 0x65, 0x28, 0x1f, 0xd4, 0xf9, 0x16, 0x18, 0xb2, 0x32, 0xe9, 0x42, 0xe2,
	0x9a, 0xb2, 0xa0, 0x90, 0x58, 0x09, 0x8b, 0xa0, 0xbd, 0xe0, 0xd4, 0x74, 0x81, 0xa5, 0x94, 0xb2,
	0xc8, 0x4b, 0xe8, 0xb8, 0x10, 0xa9, 0x98, 0x4e, 0x33, 0x1d, 0x36, 0x5d, 0x7d, 0x0c, 0xf8, 0xd1,
	0x60, 0xe4, 0x35, 0x34, 0x91, 0x6b, 0x39, 0x0f, 0x5b, 0x51, 0xa3, 0x17, 0x5c, 0x1e, 0xd8, 0x56,
	0x5d, 0x5c, 0x89, 0x49, 0x29, 0x62, 0x9e, 0x58, 0x2f, 0x79, 0x02, 0x7e, 0xa6, 0x46, 0xb4, 0x28,
	0xf2, 0x79, 0xb8, 0x67, 0xa4, 0xef, 0x65, 0xaa, 0x5f, 0x9a, 0xf1, 0x1d, 0x1c, 0x7e, 0x46, 0x2a,
	0xf5, 0x2d, 0x52, 0xbd, 0xb5, 0xe2, 0xb5, 0x5c, 0x1b, 0xeb, 0xb9, 0xc6, 0x7d, 0x38, 0xaa, 0x45,
	0x5a, 0x69, 0x56, 0x3d, 0x54, 0x08, 0x7b, 0x6a, 0x96, 0xa6, 0xa8, 0x94, 0x09, 0xe4, 0x27, 0x0b,
	0x33, 0xfe, 0xed, 0xc1, 0xe9, 0x90, 0x2b, 0x4d, 0xf3, 0xfc, 0x9a, 0xd3, 0x42, 0xdd, 0x89, 0xed,
	0x73, 0xbe, 0x80, 0x63, 0x33, 0x5a, 0x19, 0x4f, 0xf3, 0x19, 0x43, 0xb6, 0xd4, 0xab, 0xa3, 0xd2,
	0x35, 0x74, 0x1e, 0xdb, 0xb1, 0x37, 0x40, 0x96, 0xf9, 0xb5, 0xbe, 0x1d, 0xd6, 0xe9, 0xa6, 0x7b,
	0xa7, 0xd0, 0x12, 0xe3, 0xb1, 0xc2, 0x45, 0xdb, 0x9c, 0x55, 0xa6, 0xc9, 0xa8, 0xa6, 0x61, 0x2b,
	0xf2, 0x7a, 0xed, 0xc4, 0x9c, 0x0d, 0x26, 0x38, 0xba, 0xce, 0x98, 0x73, 0xfc, 0x16, 0xce, 0xd6,
	0x84, 0x3e, 0x5e, 0xb2, 0x78, 0x00, 0xff, 0xaf, 0xcc, 0xee, 0x56, 0xf5, 0x95, 0xe0, 0x2f, 0x46,
	0x87, 0x44, 0x10, 0xa4, 0x33, 0x29, 0x91, 0xeb, 0x9b, 0x87, 0x0b, 0xea, 0x10, 0x39, 0x81, 0xa6,
	0xad, 0x99, 0xdd, 0x34, 0x6b, 0x54, 0x0a, 0x1b, 0x35, 0x85, 0x5d, 0xf0, 0x0b, 0xa1, 0x32, 0x9d,
	0x09, 0xee, 0x2a, 0x56, 0xd9, 0x97, 0xbf, 0x76, 0x20, 0x48, 0xe8, 0x58, 0x5f, 0xdb, 0xc7, 0x86,
	0xbc, 0x83, 0xdd, 0x72, 0x9d, 0x09, 0x71, 0xb3, 0x5c, 0x7b, 0x6b, 0xba, 0xc7, 0x4b, 0x98, 0x93,
	0xf8, 0x05, 0x3a, 0x4b, 0xda, 0xc9, 0xb9, 0x63, 0x6d, 0xda, 0xe6, 0xee, 0xd3, 0xcd, 0x4e, 0x77,
	0xd7, 0x57, 0x38, 0x58, 0x29, 0x3b, 0x79, 0xe6, 0x7e, 0xd8, 0x3c, 0x77, 0xdd, 0xe7, 0x8f, 0xb9,
	0xdd, 0x8d, 0x1f, 0x60, 0xbf, 0x9a, 0x7a, 0x72, 0xe6, 0xc8, 0xab, 0x1b, 0xd7, 0x0d, 0xd7, 0x1d,
	0xf6, 0xff, 0xdb, 0x96, 0x71, 0xbc, 0xff, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x3e, 0x6d, 0xfb,
	0x8d, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RaftServiceClient is the client API for RaftService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RaftServiceClient interface {
	// 投票方法
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
	// 日志复制
	AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	// 安装快照
	InstallSnapshot(ctx context.Context, in *InstallSnapshotRequest, opts ...grpc.CallOption) (*InstallSnapshotResponse, error)
	//心跳检查
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type raftServiceClient struct {
	cc *grpc.ClientConn
}

func NewRaftServiceClient(cc *grpc.ClientConn) RaftServiceClient {
	return &raftServiceClient{cc}
}

func (c *raftServiceClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, "/proto.RaftService/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	out := new(AppendEntriesResponse)
	err := c.cc.Invoke(ctx, "/proto.RaftService/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) InstallSnapshot(ctx context.Context, in *InstallSnapshotRequest, opts ...grpc.CallOption) (*InstallSnapshotResponse, error) {
	out := new(InstallSnapshotResponse)
	err := c.cc.Invoke(ctx, "/proto.RaftService/InstallSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/proto.RaftService/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServiceServer is the server API for RaftService service.
type RaftServiceServer interface {
	// 投票方法
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
	// 日志复制
	AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	// 安装快照
	InstallSnapshot(context.Context, *InstallSnapshotRequest) (*InstallSnapshotResponse, error)
	//心跳检查
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
}

// UnimplementedRaftServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRaftServiceServer struct {
}

func (*UnimplementedRaftServiceServer) Vote(ctx context.Context, req *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (*UnimplementedRaftServiceServer) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (*UnimplementedRaftServiceServer) InstallSnapshot(ctx context.Context, req *InstallSnapshotRequest) (*InstallSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallSnapshot not implemented")
}
func (*UnimplementedRaftServiceServer) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}

func RegisterRaftServiceServer(s *grpc.Server, srv RaftServiceServer) {
	s.RegisterService(&_RaftService_serviceDesc, srv)
}

func _RaftService_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RaftService/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RaftService/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).AppendEntries(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_InstallSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallSnapshotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).InstallSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RaftService/InstallSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).InstallSnapshot(ctx, req.(*InstallSnapshotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RaftService/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RaftService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RaftService",
	HandlerType: (*RaftServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Vote",
			Handler:    _RaftService_Vote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _RaftService_AppendEntries_Handler,
		},
		{
			MethodName: "InstallSnapshot",
			Handler:    _RaftService_InstallSnapshot_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _RaftService_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
