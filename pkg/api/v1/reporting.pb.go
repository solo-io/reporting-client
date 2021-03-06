// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/v1/reporting.proto

package v1

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// information that every reporting source should include
type Product struct {
	Product              string   `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Arch                 string   `protobuf:"bytes,3,opt,name=arch,proto3" json:"arch,omitempty"`
	Os                   string   `protobuf:"bytes,4,opt,name=os,proto3" json:"os,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_9df4dc8e97d3e9ee, []int{0}
}

func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (m *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(m, src)
}
func (m *Product) XXX_Size() int {
	return xxx_messageInfo_Product.Size(m)
}
func (m *Product) XXX_DiscardUnknown() {
	xxx_messageInfo_Product.DiscardUnknown(m)
}

var xxx_messageInfo_Product proto.InternalMessageInfo

func (m *Product) GetProduct() string {
	if m != nil {
		return m.Product
	}
	return ""
}

func (m *Product) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Product) GetArch() string {
	if m != nil {
		return m.Arch
	}
	return ""
}

func (m *Product) GetOs() string {
	if m != nil {
		return m.Os
	}
	return ""
}

type InstanceMetadata struct {
	Product *Product `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
	// should be unique per reporting source, and unchanging over the life of the reporting source
	// this repo offers a SignatureManager type that can maintain a unique signature
	Signature            string   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstanceMetadata) Reset()         { *m = InstanceMetadata{} }
func (m *InstanceMetadata) String() string { return proto.CompactTextString(m) }
func (*InstanceMetadata) ProtoMessage()    {}
func (*InstanceMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_9df4dc8e97d3e9ee, []int{1}
}

func (m *InstanceMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstanceMetadata.Unmarshal(m, b)
}
func (m *InstanceMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstanceMetadata.Marshal(b, m, deterministic)
}
func (m *InstanceMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstanceMetadata.Merge(m, src)
}
func (m *InstanceMetadata) XXX_Size() int {
	return xxx_messageInfo_InstanceMetadata.Size(m)
}
func (m *InstanceMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_InstanceMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_InstanceMetadata proto.InternalMessageInfo

func (m *InstanceMetadata) GetProduct() *Product {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *InstanceMetadata) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type UsageRequest struct {
	InstanceMetadata *InstanceMetadata `protobuf:"bytes,1,opt,name=instance_metadata,json=instanceMetadata,proto3" json:"instance_metadata,omitempty"`
	// arbitrary key/value pairs - each reporting source can choose what to include here
	Payload              map[string]string `protobuf:"bytes,2,rep,name=payload,proto3" json:"payload,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UsageRequest) Reset()         { *m = UsageRequest{} }
func (m *UsageRequest) String() string { return proto.CompactTextString(m) }
func (*UsageRequest) ProtoMessage()    {}
func (*UsageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9df4dc8e97d3e9ee, []int{2}
}

func (m *UsageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UsageRequest.Unmarshal(m, b)
}
func (m *UsageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UsageRequest.Marshal(b, m, deterministic)
}
func (m *UsageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UsageRequest.Merge(m, src)
}
func (m *UsageRequest) XXX_Size() int {
	return xxx_messageInfo_UsageRequest.Size(m)
}
func (m *UsageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UsageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UsageRequest proto.InternalMessageInfo

func (m *UsageRequest) GetInstanceMetadata() *InstanceMetadata {
	if m != nil {
		return m.InstanceMetadata
	}
	return nil
}

func (m *UsageRequest) GetPayload(ctx context.Context) map[string]string {
	if m != nil {
		return m.Payload
	}
	return nil
}

type UsageResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UsageResponse) Reset()         { *m = UsageResponse{} }
func (m *UsageResponse) String() string { return proto.CompactTextString(m) }
func (*UsageResponse) ProtoMessage()    {}
func (*UsageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9df4dc8e97d3e9ee, []int{3}
}

func (m *UsageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UsageResponse.Unmarshal(m, b)
}
func (m *UsageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UsageResponse.Marshal(b, m, deterministic)
}
func (m *UsageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UsageResponse.Merge(m, src)
}
func (m *UsageResponse) XXX_Size() int {
	return xxx_messageInfo_UsageResponse.Size(m)
}
func (m *UsageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UsageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UsageResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Product)(nil), "reportingapi.solo.io.Product")
	proto.RegisterType((*InstanceMetadata)(nil), "reportingapi.solo.io.InstanceMetadata")
	proto.RegisterType((*UsageRequest)(nil), "reportingapi.solo.io.UsageRequest")
	proto.RegisterMapType((map[string]string)(nil), "reportingapi.solo.io.UsageRequest.PayloadEntry")
	proto.RegisterType((*UsageResponse)(nil), "reportingapi.solo.io.UsageResponse")
}

func init() { proto.RegisterFile("api/v1/reporting.proto", fileDescriptor_9df4dc8e97d3e9ee) }

var fileDescriptor_9df4dc8e97d3e9ee = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0x86, 0xbf, 0xa4, 0xfd, 0x2c, 0x9d, 0x56, 0x8d, 0x4b, 0x91, 0x50, 0x14, 0x4a, 0x04, 0xe9,
	0x29, 0xa5, 0xf5, 0xa0, 0xf4, 0x28, 0x78, 0xe8, 0x41, 0x28, 0x29, 0x82, 0x78, 0x91, 0x31, 0x5d,
	0xea, 0x62, 0xdc, 0x5d, 0x77, 0x37, 0x81, 0xfe, 0x65, 0x7f, 0x85, 0x34, 0xbb, 0xb1, 0xb5, 0x04,
	0xbd, 0xcd, 0xbc, 0x33, 0x93, 0xe7, 0x49, 0x08, 0x9c, 0xa2, 0x64, 0xa3, 0x62, 0x3c, 0x52, 0x54,
	0x0a, 0x65, 0x18, 0x5f, 0xc5, 0x52, 0x09, 0x23, 0x48, 0xef, 0x3b, 0x40, 0xc9, 0x62, 0x2d, 0x32,
	0x11, 0x33, 0x11, 0x21, 0xb4, 0xe6, 0x4a, 0x2c, 0xf3, 0xd4, 0x90, 0x10, 0x5a, 0xd2, 0x96, 0xa1,
	0x37, 0xf0, 0x86, 0xed, 0xa4, 0x6a, 0x37, 0x93, 0x82, 0x2a, 0xcd, 0x04, 0x0f, 0x7d, 0x3b, 0x71,
	0x2d, 0x21, 0xd0, 0x44, 0x95, 0xbe, 0x86, 0x8d, 0x32, 0x2e, 0x6b, 0x72, 0x04, 0xbe, 0xd0, 0x61,
	0xb3, 0x4c, 0x7c, 0xa1, 0x23, 0x06, 0xc1, 0x8c, 0x6b, 0x83, 0x3c, 0xa5, 0xf7, 0xd4, 0xe0, 0x12,
	0x0d, 0x92, 0xeb, 0x9f, 0xac, 0xce, 0xe4, 0x3c, 0xae, 0xd3, 0x8b, 0x9d, 0xdb, 0x56, 0xe5, 0x0c,
	0xda, 0x9a, 0xad, 0x38, 0x9a, 0x5c, 0x51, 0x27, 0xb3, 0x0d, 0xa2, 0x4f, 0x0f, 0xba, 0x0f, 0x1a,
	0x57, 0x34, 0xa1, 0x1f, 0x39, 0xd5, 0x86, 0x2c, 0xe0, 0x84, 0x39, 0xf6, 0xf3, 0xbb, 0x83, 0x3b,
	0xe2, 0x65, 0x3d, 0x71, 0x5f, 0x35, 0x09, 0xd8, 0xbe, 0xfc, 0x0c, 0x5a, 0x12, 0xd7, 0x99, 0xc0,
	0x65, 0xe8, 0x0f, 0x1a, 0xc3, 0xce, 0x64, 0x54, 0xff, 0xa8, 0x5d, 0x93, 0x78, 0x6e, 0x2f, 0xee,
	0xb8, 0x51, 0xeb, 0xa4, 0xba, 0xef, 0x4f, 0xa1, 0xbb, 0x3b, 0x20, 0x01, 0x34, 0xde, 0xe8, 0xda,
	0x7d, 0xff, 0x4d, 0x49, 0x7a, 0xf0, 0xbf, 0xc0, 0x2c, 0xaf, 0x5e, 0xd6, 0x36, 0x53, 0xff, 0xc6,
	0x8b, 0x8e, 0xe1, 0xd0, 0x11, 0xb4, 0x14, 0x5c, 0xd3, 0x49, 0x06, 0x41, 0x52, 0x79, 0x2c, 0xa8,
	0x2a, 0x58, 0x4a, 0xc9, 0x23, 0x74, 0x6c, 0x56, 0xae, 0x92, 0xe8, 0x6f, 0xd3, 0xfe, 0xc5, 0xaf,
	0x3b, 0x96, 0x15, 0xfd, 0xbb, 0x6d, 0x3e, 0xf9, 0xc5, 0xf8, 0xe5, 0xa0, 0xfc, 0xb9, 0xae, 0xbe,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x91, 0x89, 0xea, 0x76, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReportingServiceClient is the client API for ReportingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReportingServiceClient interface {
	ReportUsage(ctx context.Context, in *UsageRequest, opts ...grpc.CallOption) (*UsageResponse, error)
}

type reportingServiceClient struct {
	cc *grpc.ClientConn
}

func NewReportingServiceClient(cc *grpc.ClientConn) ReportingServiceClient {
	return &reportingServiceClient{cc}
}

func (c *reportingServiceClient) ReportUsage(ctx context.Context, in *UsageRequest, opts ...grpc.CallOption) (*UsageResponse, error) {
	out := new(UsageResponse)
	err := c.cc.Invoke(ctx, "/reportingapi.solo.io.ReportingService/ReportUsage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportingServiceServer is the server API for ReportingService service.
type ReportingServiceServer interface {
	ReportUsage(context.Context, *UsageRequest) (*UsageResponse, error)
}

// UnimplementedReportingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedReportingServiceServer struct {
}

func (*UnimplementedReportingServiceServer) ReportUsage(ctx context.Context, req *UsageRequest) (*UsageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportUsage not implemented")
}

func RegisterReportingServiceServer(s *grpc.Server, srv ReportingServiceServer) {
	s.RegisterService(&_ReportingService_serviceDesc, srv)
}

func _ReportingService_ReportUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UsageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportingServiceServer).ReportUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/reportingapi.solo.io.ReportingService/ReportUsage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportingServiceServer).ReportUsage(ctx, req.(*UsageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReportingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "reportingapi.solo.io.ReportingService",
	HandlerType: (*ReportingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportUsage",
			Handler:    _ReportingService_ReportUsage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/reporting.proto",
}
