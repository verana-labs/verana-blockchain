// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: verana/perm/v1/query.proto

package permv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Query_Params_FullMethodName                 = "/verana.perm.v1.Query/Params"
	Query_ListPermissions_FullMethodName        = "/verana.perm.v1.Query/ListPermissions"
	Query_GetPermission_FullMethodName          = "/verana.perm.v1.Query/GetPermission"
	Query_GetPermissionSession_FullMethodName   = "/verana.perm.v1.Query/GetPermissionSession"
	Query_ListPermissionSessions_FullMethodName = "/verana.perm.v1.Query/ListPermissionSessions"
	Query_FindPermissionsWithDID_FullMethodName = "/verana.perm.v1.Query/FindPermissionsWithDID"
	Query_FindBeneficiaries_FullMethodName      = "/verana.perm.v1.Query/FindBeneficiaries"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	ListPermissions(ctx context.Context, in *QueryListPermissionsRequest, opts ...grpc.CallOption) (*QueryListPermissionsResponse, error)
	GetPermission(ctx context.Context, in *QueryGetPermissionRequest, opts ...grpc.CallOption) (*QueryGetPermissionResponse, error)
	GetPermissionSession(ctx context.Context, in *QueryGetPermissionSessionRequest, opts ...grpc.CallOption) (*QueryGetPermissionSessionResponse, error)
	ListPermissionSessions(ctx context.Context, in *QueryListPermissionSessionsRequest, opts ...grpc.CallOption) (*QueryListPermissionSessionsResponse, error)
	FindPermissionsWithDID(ctx context.Context, in *QueryFindPermissionsWithDIDRequest, opts ...grpc.CallOption) (*QueryFindPermissionsWithDIDResponse, error)
	FindBeneficiaries(ctx context.Context, in *QueryFindBeneficiariesRequest, opts ...grpc.CallOption) (*QueryFindBeneficiariesResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListPermissions(ctx context.Context, in *QueryListPermissionsRequest, opts ...grpc.CallOption) (*QueryListPermissionsResponse, error) {
	out := new(QueryListPermissionsResponse)
	err := c.cc.Invoke(ctx, Query_ListPermissions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetPermission(ctx context.Context, in *QueryGetPermissionRequest, opts ...grpc.CallOption) (*QueryGetPermissionResponse, error) {
	out := new(QueryGetPermissionResponse)
	err := c.cc.Invoke(ctx, Query_GetPermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetPermissionSession(ctx context.Context, in *QueryGetPermissionSessionRequest, opts ...grpc.CallOption) (*QueryGetPermissionSessionResponse, error) {
	out := new(QueryGetPermissionSessionResponse)
	err := c.cc.Invoke(ctx, Query_GetPermissionSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListPermissionSessions(ctx context.Context, in *QueryListPermissionSessionsRequest, opts ...grpc.CallOption) (*QueryListPermissionSessionsResponse, error) {
	out := new(QueryListPermissionSessionsResponse)
	err := c.cc.Invoke(ctx, Query_ListPermissionSessions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FindPermissionsWithDID(ctx context.Context, in *QueryFindPermissionsWithDIDRequest, opts ...grpc.CallOption) (*QueryFindPermissionsWithDIDResponse, error) {
	out := new(QueryFindPermissionsWithDIDResponse)
	err := c.cc.Invoke(ctx, Query_FindPermissionsWithDID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FindBeneficiaries(ctx context.Context, in *QueryFindBeneficiariesRequest, opts ...grpc.CallOption) (*QueryFindBeneficiariesResponse, error) {
	out := new(QueryFindBeneficiariesResponse)
	err := c.cc.Invoke(ctx, Query_FindBeneficiaries_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	ListPermissions(context.Context, *QueryListPermissionsRequest) (*QueryListPermissionsResponse, error)
	GetPermission(context.Context, *QueryGetPermissionRequest) (*QueryGetPermissionResponse, error)
	GetPermissionSession(context.Context, *QueryGetPermissionSessionRequest) (*QueryGetPermissionSessionResponse, error)
	ListPermissionSessions(context.Context, *QueryListPermissionSessionsRequest) (*QueryListPermissionSessionsResponse, error)
	FindPermissionsWithDID(context.Context, *QueryFindPermissionsWithDIDRequest) (*QueryFindPermissionsWithDIDResponse, error)
	FindBeneficiaries(context.Context, *QueryFindBeneficiariesRequest) (*QueryFindBeneficiariesResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) ListPermissions(context.Context, *QueryListPermissionsRequest) (*QueryListPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPermissions not implemented")
}
func (UnimplementedQueryServer) GetPermission(context.Context, *QueryGetPermissionRequest) (*QueryGetPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermission not implemented")
}
func (UnimplementedQueryServer) GetPermissionSession(context.Context, *QueryGetPermissionSessionRequest) (*QueryGetPermissionSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermissionSession not implemented")
}
func (UnimplementedQueryServer) ListPermissionSessions(context.Context, *QueryListPermissionSessionsRequest) (*QueryListPermissionSessionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPermissionSessions not implemented")
}
func (UnimplementedQueryServer) FindPermissionsWithDID(context.Context, *QueryFindPermissionsWithDIDRequest) (*QueryFindPermissionsWithDIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindPermissionsWithDID not implemented")
}
func (UnimplementedQueryServer) FindBeneficiaries(context.Context, *QueryFindBeneficiariesRequest) (*QueryFindBeneficiariesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindBeneficiaries not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListPermissions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListPermissions(ctx, req.(*QueryListPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetPermission(ctx, req.(*QueryGetPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetPermissionSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetPermissionSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetPermissionSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetPermissionSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetPermissionSession(ctx, req.(*QueryGetPermissionSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListPermissionSessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListPermissionSessionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListPermissionSessions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListPermissionSessions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListPermissionSessions(ctx, req.(*QueryListPermissionSessionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FindPermissionsWithDID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFindPermissionsWithDIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FindPermissionsWithDID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_FindPermissionsWithDID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FindPermissionsWithDID(ctx, req.(*QueryFindPermissionsWithDIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FindBeneficiaries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFindBeneficiariesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FindBeneficiaries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_FindBeneficiaries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FindBeneficiaries(ctx, req.(*QueryFindBeneficiariesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "verana.perm.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ListPermissions",
			Handler:    _Query_ListPermissions_Handler,
		},
		{
			MethodName: "GetPermission",
			Handler:    _Query_GetPermission_Handler,
		},
		{
			MethodName: "GetPermissionSession",
			Handler:    _Query_GetPermissionSession_Handler,
		},
		{
			MethodName: "ListPermissionSessions",
			Handler:    _Query_ListPermissionSessions_Handler,
		},
		{
			MethodName: "FindPermissionsWithDID",
			Handler:    _Query_FindPermissionsWithDID_Handler,
		},
		{
			MethodName: "FindBeneficiaries",
			Handler:    _Query_FindBeneficiaries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "verana/perm/v1/query.proto",
}
