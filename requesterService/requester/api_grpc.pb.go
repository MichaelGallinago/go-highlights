// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/proto/api.proto

package requester

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RequesterService_GetTopLongMemes_FullMethodName        = "/api.RequesterService/GetTopLongMemes"
	RequesterService_SearchMemesBySubstring_FullMethodName = "/api.RequesterService/SearchMemesBySubstring"
	RequesterService_GetMemesByMonth_FullMethodName        = "/api.RequesterService/GetMemesByMonth"
	RequesterService_GetRandomMeme_FullMethodName          = "/api.RequesterService/GetRandomMeme"
)

// RequesterServiceClient is the client API for RequesterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RequesterServiceClient interface {
	GetTopLongMemes(ctx context.Context, in *TopLongMemesRequest, opts ...grpc.CallOption) (*MemesResponse, error)
	SearchMemesBySubstring(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*MemesResponse, error)
	GetMemesByMonth(ctx context.Context, in *MonthRequest, opts ...grpc.CallOption) (*MemesResponse, error)
	GetRandomMeme(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MemeResponse, error)
}

type requesterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRequesterServiceClient(cc grpc.ClientConnInterface) RequesterServiceClient {
	return &requesterServiceClient{cc}
}

func (c *requesterServiceClient) GetTopLongMemes(ctx context.Context, in *TopLongMemesRequest, opts ...grpc.CallOption) (*MemesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MemesResponse)
	err := c.cc.Invoke(ctx, RequesterService_GetTopLongMemes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requesterServiceClient) SearchMemesBySubstring(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*MemesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MemesResponse)
	err := c.cc.Invoke(ctx, RequesterService_SearchMemesBySubstring_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requesterServiceClient) GetMemesByMonth(ctx context.Context, in *MonthRequest, opts ...grpc.CallOption) (*MemesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MemesResponse)
	err := c.cc.Invoke(ctx, RequesterService_GetMemesByMonth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requesterServiceClient) GetRandomMeme(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MemeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MemeResponse)
	err := c.cc.Invoke(ctx, RequesterService_GetRandomMeme_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RequesterServiceServer is the server API for RequesterService service.
// All implementations must embed UnimplementedRequesterServiceServer
// for forward compatibility.
type RequesterServiceServer interface {
	GetTopLongMemes(context.Context, *TopLongMemesRequest) (*MemesResponse, error)
	SearchMemesBySubstring(context.Context, *SearchRequest) (*MemesResponse, error)
	GetMemesByMonth(context.Context, *MonthRequest) (*MemesResponse, error)
	GetRandomMeme(context.Context, *Empty) (*MemeResponse, error)
	mustEmbedUnimplementedRequesterServiceServer()
}

// UnimplementedRequesterServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRequesterServiceServer struct{}

func (UnimplementedRequesterServiceServer) GetTopLongMemes(context.Context, *TopLongMemesRequest) (*MemesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopLongMemes not implemented")
}
func (UnimplementedRequesterServiceServer) SearchMemesBySubstring(context.Context, *SearchRequest) (*MemesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchMemesBySubstring not implemented")
}
func (UnimplementedRequesterServiceServer) GetMemesByMonth(context.Context, *MonthRequest) (*MemesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMemesByMonth not implemented")
}
func (UnimplementedRequesterServiceServer) GetRandomMeme(context.Context, *Empty) (*MemeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRandomMeme not implemented")
}
func (UnimplementedRequesterServiceServer) mustEmbedUnimplementedRequesterServiceServer() {}
func (UnimplementedRequesterServiceServer) testEmbeddedByValue()                          {}

// UnsafeRequesterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RequesterServiceServer will
// result in compilation errors.
type UnsafeRequesterServiceServer interface {
	mustEmbedUnimplementedRequesterServiceServer()
}

func RegisterRequesterServiceServer(s grpc.ServiceRegistrar, srv RequesterServiceServer) {
	// If the following call pancis, it indicates UnimplementedRequesterServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RequesterService_ServiceDesc, srv)
}

func _RequesterService_GetTopLongMemes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopLongMemesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequesterServiceServer).GetTopLongMemes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequesterService_GetTopLongMemes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequesterServiceServer).GetTopLongMemes(ctx, req.(*TopLongMemesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequesterService_SearchMemesBySubstring_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequesterServiceServer).SearchMemesBySubstring(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequesterService_SearchMemesBySubstring_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequesterServiceServer).SearchMemesBySubstring(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequesterService_GetMemesByMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequesterServiceServer).GetMemesByMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequesterService_GetMemesByMonth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequesterServiceServer).GetMemesByMonth(ctx, req.(*MonthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequesterService_GetRandomMeme_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequesterServiceServer).GetRandomMeme(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequesterService_GetRandomMeme_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequesterServiceServer).GetRandomMeme(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// RequesterService_ServiceDesc is the grpc.ServiceDesc for RequesterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RequesterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.RequesterService",
	HandlerType: (*RequesterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopLongMemes",
			Handler:    _RequesterService_GetTopLongMemes_Handler,
		},
		{
			MethodName: "SearchMemesBySubstring",
			Handler:    _RequesterService_SearchMemesBySubstring_Handler,
		},
		{
			MethodName: "GetMemesByMonth",
			Handler:    _RequesterService_GetMemesByMonth_Handler,
		},
		{
			MethodName: "GetRandomMeme",
			Handler:    _RequesterService_GetRandomMeme_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/api.proto",
}
