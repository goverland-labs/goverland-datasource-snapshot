// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: votingpb/voting.proto

package votingpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ValidateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Voter    string `protobuf:"bytes,1,opt,name=voter,proto3" json:"voter,omitempty"`
	Proposal string `protobuf:"bytes,2,opt,name=proposal,proto3" json:"proposal,omitempty"`
}

func (x *ValidateRequest) Reset() {
	*x = ValidateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateRequest) ProtoMessage() {}

func (x *ValidateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateRequest.ProtoReflect.Descriptor instead.
func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{0}
}

func (x *ValidateRequest) GetVoter() string {
	if x != nil {
		return x.Voter
	}
	return ""
}

func (x *ValidateRequest) GetProposal() string {
	if x != nil {
		return x.Proposal
	}
	return ""
}

type ValidateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok              bool             `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	VotingPower     float64          `protobuf:"fixed64,2,opt,name=voting_power,json=votingPower,proto3" json:"voting_power,omitempty"`
	ValidationError *ValidationError `protobuf:"bytes,3,opt,name=validation_error,json=validationError,proto3,oneof" json:"validation_error,omitempty"`
}

func (x *ValidateResponse) Reset() {
	*x = ValidateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateResponse) ProtoMessage() {}

func (x *ValidateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateResponse.ProtoReflect.Descriptor instead.
func (*ValidateResponse) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{1}
}

func (x *ValidateResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *ValidateResponse) GetVotingPower() float64 {
	if x != nil {
		return x.VotingPower
	}
	return 0
}

func (x *ValidateResponse) GetValidationError() *ValidationError {
	if x != nil {
		return x.ValidationError
	}
	return nil
}

type ValidationError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Code    uint32 `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *ValidationError) Reset() {
	*x = ValidationError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationError) ProtoMessage() {}

func (x *ValidationError) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationError.ProtoReflect.Descriptor instead.
func (*ValidationError) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{2}
}

func (x *ValidationError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ValidationError) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type PrepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Voter    string     `protobuf:"bytes,1,opt,name=voter,proto3" json:"voter,omitempty"`
	Proposal string     `protobuf:"bytes,2,opt,name=proposal,proto3" json:"proposal,omitempty"`
	Choice   *anypb.Any `protobuf:"bytes,3,opt,name=choice,proto3" json:"choice,omitempty"`
	Reason   *string    `protobuf:"bytes,6,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
}

func (x *PrepareRequest) Reset() {
	*x = PrepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareRequest) ProtoMessage() {}

func (x *PrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareRequest.ProtoReflect.Descriptor instead.
func (*PrepareRequest) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{3}
}

func (x *PrepareRequest) GetVoter() string {
	if x != nil {
		return x.Voter
	}
	return ""
}

func (x *PrepareRequest) GetProposal() string {
	if x != nil {
		return x.Proposal
	}
	return ""
}

func (x *PrepareRequest) GetChoice() *anypb.Any {
	if x != nil {
		return x.Choice
	}
	return nil
}

func (x *PrepareRequest) GetReason() string {
	if x != nil && x.Reason != nil {
		return *x.Reason
	}
	return ""
}

type PrepareResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TypedData string `protobuf:"bytes,2,opt,name=typed_data,json=typedData,proto3" json:"typed_data,omitempty"`
}

func (x *PrepareResponse) Reset() {
	*x = PrepareResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareResponse) ProtoMessage() {}

func (x *PrepareResponse) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareResponse.ProtoReflect.Descriptor instead.
func (*PrepareResponse) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{4}
}

func (x *PrepareResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PrepareResponse) GetTypedData() string {
	if x != nil {
		return x.TypedData
	}
	return ""
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sig string `protobuf:"bytes,2,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{5}
}

func (x *VoteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *VoteRequest) GetSig() string {
	if x != nil {
		return x.Sig
	}
	return ""
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Ipfs    string   `protobuf:"bytes,2,opt,name=ipfs,proto3" json:"ipfs,omitempty"`
	Relayer *Relayer `protobuf:"bytes,3,opt,name=relayer,proto3" json:"relayer,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{6}
}

func (x *VoteResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *VoteResponse) GetIpfs() string {
	if x != nil {
		return x.Ipfs
	}
	return ""
}

func (x *VoteResponse) GetRelayer() *Relayer {
	if x != nil {
		return x.Relayer
	}
	return nil
}

type Relayer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Receipt string `protobuf:"bytes,2,opt,name=receipt,proto3" json:"receipt,omitempty"`
}

func (x *Relayer) Reset() {
	*x = Relayer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Relayer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relayer) ProtoMessage() {}

func (x *Relayer) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relayer.ProtoReflect.Descriptor instead.
func (*Relayer) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{7}
}

func (x *Relayer) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Relayer) GetReceipt() string {
	if x != nil {
		return x.Receipt
	}
	return ""
}

type GetVoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetVoteRequest) Reset() {
	*x = GetVoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVoteRequest) ProtoMessage() {}

func (x *GetVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVoteRequest.ProtoReflect.Descriptor instead.
func (*GetVoteRequest) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{8}
}

func (x *GetVoteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetVoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Ipfs          string     `protobuf:"bytes,2,opt,name=ipfs,proto3" json:"ipfs,omitempty"`
	Voter         string     `protobuf:"bytes,3,opt,name=voter,proto3" json:"voter,omitempty"`
	Created       int64      `protobuf:"varint,4,opt,name=created,proto3" json:"created,omitempty"`
	OriginalDaoId string     `protobuf:"bytes,5,opt,name=original_dao_id,json=originalDaoId,proto3" json:"original_dao_id,omitempty"`
	ProposalId    string     `protobuf:"bytes,6,opt,name=proposal_id,json=proposalId,proto3" json:"proposal_id,omitempty"`
	Choice        *anypb.Any `protobuf:"bytes,7,opt,name=choice,proto3" json:"choice,omitempty"`
	Reason        string     `protobuf:"bytes,8,opt,name=reason,proto3" json:"reason,omitempty"`
	App           string     `protobuf:"bytes,9,opt,name=app,proto3" json:"app,omitempty"`
	Vp            float64    `protobuf:"fixed64,10,opt,name=vp,proto3" json:"vp,omitempty"`
	VpByStrategy  []float64  `protobuf:"fixed64,11,rep,packed,name=vp_by_strategy,json=vpByStrategy,proto3" json:"vp_by_strategy,omitempty"`
	VpState       string     `protobuf:"bytes,12,opt,name=vp_state,json=vpState,proto3" json:"vp_state,omitempty"`
}

func (x *GetVoteResponse) Reset() {
	*x = GetVoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_votingpb_voting_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVoteResponse) ProtoMessage() {}

func (x *GetVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_votingpb_voting_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVoteResponse.ProtoReflect.Descriptor instead.
func (*GetVoteResponse) Descriptor() ([]byte, []int) {
	return file_votingpb_voting_proto_rawDescGZIP(), []int{9}
}

func (x *GetVoteResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetVoteResponse) GetIpfs() string {
	if x != nil {
		return x.Ipfs
	}
	return ""
}

func (x *GetVoteResponse) GetVoter() string {
	if x != nil {
		return x.Voter
	}
	return ""
}

func (x *GetVoteResponse) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *GetVoteResponse) GetOriginalDaoId() string {
	if x != nil {
		return x.OriginalDaoId
	}
	return ""
}

func (x *GetVoteResponse) GetProposalId() string {
	if x != nil {
		return x.ProposalId
	}
	return ""
}

func (x *GetVoteResponse) GetChoice() *anypb.Any {
	if x != nil {
		return x.Choice
	}
	return nil
}

func (x *GetVoteResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *GetVoteResponse) GetApp() string {
	if x != nil {
		return x.App
	}
	return ""
}

func (x *GetVoteResponse) GetVp() float64 {
	if x != nil {
		return x.Vp
	}
	return 0
}

func (x *GetVoteResponse) GetVpByStrategy() []float64 {
	if x != nil {
		return x.VpByStrategy
	}
	return nil
}

func (x *GetVoteResponse) GetVpState() string {
	if x != nil {
		return x.VpState
	}
	return ""
}

var File_votingpb_voting_proto protoreflect.FileDescriptor

var file_votingpb_voting_proto_rawDesc = []byte{
	0x0a, 0x15, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2f, 0x76, 0x6f, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70,
	0x62, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x0f,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x6f, 0x74, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x22, 0xa5, 0x01, 0x0a, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67,
	0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x76, 0x6f,
	0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x49, 0x0a, 0x10, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00,
	0x52, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x88, 0x01, 0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3f, 0x0a, 0x0f, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x98, 0x01, 0x0a, 0x0e, 0x50,
	0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x6f, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x6f,
	0x74, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12,
	0x2c, 0x0a, 0x06, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x06, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a,
	0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x40, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x79, 0x70, 0x65,
	0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x79,
	0x70, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x22, 0x2f, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x69, 0x67, 0x22, 0x5f, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x70, 0x66, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x70, 0x66, 0x73, 0x12, 0x2b, 0x0a, 0x07,
	0x72, 0x65, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x52, 0x07, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x22, 0x3d, 0x0a, 0x07, 0x52, 0x65, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xd7, 0x02, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x70, 0x66, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x70,
	0x66, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x64,
	0x61, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x61, 0x6c, 0x44, 0x61, 0x6f, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x06, 0x63,
	0x68, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x06, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x61, 0x70, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x76, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x02, 0x76, 0x70, 0x12, 0x24, 0x0a, 0x0e, 0x76, 0x70, 0x5f, 0x62, 0x79, 0x5f, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x01, 0x52, 0x0c, 0x76, 0x70, 0x42,
	0x79, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x70, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x70, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x32, 0x82, 0x02, 0x0a, 0x06, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x41, 0x0a, 0x08, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x76, 0x6f,
	0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70,
	0x62, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x18, 0x2e,
	0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67,
	0x70, 0x62, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x15, 0x2e, 0x76, 0x6f, 0x74,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x47, 0x65, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x76,
	0x6f, 0x74, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_votingpb_voting_proto_rawDescOnce sync.Once
	file_votingpb_voting_proto_rawDescData = file_votingpb_voting_proto_rawDesc
)

func file_votingpb_voting_proto_rawDescGZIP() []byte {
	file_votingpb_voting_proto_rawDescOnce.Do(func() {
		file_votingpb_voting_proto_rawDescData = protoimpl.X.CompressGZIP(file_votingpb_voting_proto_rawDescData)
	})
	return file_votingpb_voting_proto_rawDescData
}

var file_votingpb_voting_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_votingpb_voting_proto_goTypes = []interface{}{
	(*ValidateRequest)(nil),  // 0: votingpb.ValidateRequest
	(*ValidateResponse)(nil), // 1: votingpb.ValidateResponse
	(*ValidationError)(nil),  // 2: votingpb.ValidationError
	(*PrepareRequest)(nil),   // 3: votingpb.PrepareRequest
	(*PrepareResponse)(nil),  // 4: votingpb.PrepareResponse
	(*VoteRequest)(nil),      // 5: votingpb.VoteRequest
	(*VoteResponse)(nil),     // 6: votingpb.VoteResponse
	(*Relayer)(nil),          // 7: votingpb.Relayer
	(*GetVoteRequest)(nil),   // 8: votingpb.GetVoteRequest
	(*GetVoteResponse)(nil),  // 9: votingpb.GetVoteResponse
	(*anypb.Any)(nil),        // 10: google.protobuf.Any
}
var file_votingpb_voting_proto_depIdxs = []int32{
	2,  // 0: votingpb.ValidateResponse.validation_error:type_name -> votingpb.ValidationError
	10, // 1: votingpb.PrepareRequest.choice:type_name -> google.protobuf.Any
	7,  // 2: votingpb.VoteResponse.relayer:type_name -> votingpb.Relayer
	10, // 3: votingpb.GetVoteResponse.choice:type_name -> google.protobuf.Any
	0,  // 4: votingpb.Voting.Validate:input_type -> votingpb.ValidateRequest
	3,  // 5: votingpb.Voting.Prepare:input_type -> votingpb.PrepareRequest
	5,  // 6: votingpb.Voting.Vote:input_type -> votingpb.VoteRequest
	8,  // 7: votingpb.Voting.GetVote:input_type -> votingpb.GetVoteRequest
	1,  // 8: votingpb.Voting.Validate:output_type -> votingpb.ValidateResponse
	4,  // 9: votingpb.Voting.Prepare:output_type -> votingpb.PrepareResponse
	6,  // 10: votingpb.Voting.Vote:output_type -> votingpb.VoteResponse
	9,  // 11: votingpb.Voting.GetVote:output_type -> votingpb.GetVoteResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_votingpb_voting_proto_init() }
func file_votingpb_voting_proto_init() {
	if File_votingpb_voting_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_votingpb_voting_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateRequest); i {
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
		file_votingpb_voting_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateResponse); i {
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
		file_votingpb_voting_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationError); i {
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
		file_votingpb_voting_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareRequest); i {
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
		file_votingpb_voting_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareResponse); i {
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
		file_votingpb_voting_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
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
		file_votingpb_voting_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
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
		file_votingpb_voting_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Relayer); i {
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
		file_votingpb_voting_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVoteRequest); i {
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
		file_votingpb_voting_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVoteResponse); i {
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
	file_votingpb_voting_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_votingpb_voting_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_votingpb_voting_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_votingpb_voting_proto_goTypes,
		DependencyIndexes: file_votingpb_voting_proto_depIdxs,
		MessageInfos:      file_votingpb_voting_proto_msgTypes,
	}.Build()
	File_votingpb_voting_proto = out.File
	file_votingpb_voting_proto_rawDesc = nil
	file_votingpb_voting_proto_goTypes = nil
	file_votingpb_voting_proto_depIdxs = nil
}
