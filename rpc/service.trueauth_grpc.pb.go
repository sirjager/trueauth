// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: service.trueauth.proto

package rpc

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
	TrueAuth_Welcome_FullMethodName  = "/trueauth.TrueAuth/Welcome"
	TrueAuth_Health_FullMethodName   = "/trueauth.TrueAuth/Health"
	TrueAuth_Signup_FullMethodName   = "/trueauth.TrueAuth/Signup"
	TrueAuth_Signin_FullMethodName   = "/trueauth.TrueAuth/Signin"
	TrueAuth_Validate_FullMethodName = "/trueauth.TrueAuth/Validate"
	TrueAuth_Signout_FullMethodName  = "/trueauth.TrueAuth/Signout"
	TrueAuth_Refresh_FullMethodName  = "/trueauth.TrueAuth/Refresh"
	TrueAuth_Reset_FullMethodName    = "/trueauth.TrueAuth/Reset"
	TrueAuth_Verify_FullMethodName   = "/trueauth.TrueAuth/Verify"
	TrueAuth_Delete_FullMethodName   = "/trueauth.TrueAuth/Delete"
	TrueAuth_User_FullMethodName     = "/trueauth.TrueAuth/User"
)

// TrueAuthClient is the client API for TrueAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrueAuthClient interface {
	// Welcome endpoint returns a welcome message.
	Welcome(ctx context.Context, in *WelcomeRequest, opts ...grpc.CallOption) (*WelcomeResponse, error)
	// Health endpoint returns the health status of the API.
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	Signup(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error)
	Signin(ctx context.Context, in *SigninRequest, opts ...grpc.CallOption) (*SigninResponse, error)
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	Signout(ctx context.Context, in *SignoutRequest, opts ...grpc.CallOption) (*SignoutResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
	Reset(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error)
	Verify(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	User(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type trueAuthClient struct {
	cc grpc.ClientConnInterface
}

func NewTrueAuthClient(cc grpc.ClientConnInterface) TrueAuthClient {
	return &trueAuthClient{cc}
}

func (c *trueAuthClient) Welcome(ctx context.Context, in *WelcomeRequest, opts ...grpc.CallOption) (*WelcomeResponse, error) {
	out := new(WelcomeResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Welcome_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Health_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Signup(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error) {
	out := new(SignupResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Signup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Signin(ctx context.Context, in *SigninRequest, opts ...grpc.CallOption) (*SigninResponse, error) {
	out := new(SigninResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Signin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Validate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Signout(ctx context.Context, in *SignoutRequest, opts ...grpc.CallOption) (*SignoutResponse, error) {
	out := new(SignoutResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Signout_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Refresh_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Reset(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error) {
	out := new(ResetResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Reset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Verify(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Verify_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, TrueAuth_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trueAuthClient) User(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, TrueAuth_User_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrueAuthServer is the server API for TrueAuth service.
// All implementations must embed UnimplementedTrueAuthServer
// for forward compatibility
type TrueAuthServer interface {
	// Welcome endpoint returns a welcome message.
	Welcome(context.Context, *WelcomeRequest) (*WelcomeResponse, error)
	// Health endpoint returns the health status of the API.
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	Signup(context.Context, *SignupRequest) (*SignupResponse, error)
	Signin(context.Context, *SigninRequest) (*SigninResponse, error)
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
	Signout(context.Context, *SignoutRequest) (*SignoutResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	Reset(context.Context, *ResetRequest) (*ResetResponse, error)
	Verify(context.Context, *VerifyRequest) (*VerifyResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	User(context.Context, *UserRequest) (*UserResponse, error)
	mustEmbedUnimplementedTrueAuthServer()
}

// UnimplementedTrueAuthServer must be embedded to have forward compatible implementations.
type UnimplementedTrueAuthServer struct {
}

func (UnimplementedTrueAuthServer) Welcome(context.Context, *WelcomeRequest) (*WelcomeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Welcome not implemented")
}
func (UnimplementedTrueAuthServer) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedTrueAuthServer) Signup(context.Context, *SignupRequest) (*SignupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}
func (UnimplementedTrueAuthServer) Signin(context.Context, *SigninRequest) (*SigninResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signin not implemented")
}
func (UnimplementedTrueAuthServer) Validate(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedTrueAuthServer) Signout(context.Context, *SignoutRequest) (*SignoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signout not implemented")
}
func (UnimplementedTrueAuthServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedTrueAuthServer) Reset(context.Context, *ResetRequest) (*ResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}
func (UnimplementedTrueAuthServer) Verify(context.Context, *VerifyRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedTrueAuthServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTrueAuthServer) User(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method User not implemented")
}
func (UnimplementedTrueAuthServer) mustEmbedUnimplementedTrueAuthServer() {}

// UnsafeTrueAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrueAuthServer will
// result in compilation errors.
type UnsafeTrueAuthServer interface {
	mustEmbedUnimplementedTrueAuthServer()
}

func RegisterTrueAuthServer(s grpc.ServiceRegistrar, srv TrueAuthServer) {
	s.RegisterService(&TrueAuth_ServiceDesc, srv)
}

func _TrueAuth_Welcome_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WelcomeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Welcome(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Welcome_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Welcome(ctx, req.(*WelcomeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Signup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Signup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Signup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Signup(ctx, req.(*SignupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Signin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SigninRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Signin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Signin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Signin(ctx, req.(*SigninRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Validate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Signout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Signout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Signout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Signout(ctx, req.(*SignoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Refresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Reset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Reset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Reset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Reset(ctx, req.(*ResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Verify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Verify(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrueAuth_User_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrueAuthServer).User(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrueAuth_User_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrueAuthServer).User(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TrueAuth_ServiceDesc is the grpc.ServiceDesc for TrueAuth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TrueAuth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trueauth.TrueAuth",
	HandlerType: (*TrueAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Welcome",
			Handler:    _TrueAuth_Welcome_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _TrueAuth_Health_Handler,
		},
		{
			MethodName: "Signup",
			Handler:    _TrueAuth_Signup_Handler,
		},
		{
			MethodName: "Signin",
			Handler:    _TrueAuth_Signin_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _TrueAuth_Validate_Handler,
		},
		{
			MethodName: "Signout",
			Handler:    _TrueAuth_Signout_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _TrueAuth_Refresh_Handler,
		},
		{
			MethodName: "Reset",
			Handler:    _TrueAuth_Reset_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _TrueAuth_Verify_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TrueAuth_Delete_Handler,
		},
		{
			MethodName: "User",
			Handler:    _TrueAuth_User_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.trueauth.proto",
}
