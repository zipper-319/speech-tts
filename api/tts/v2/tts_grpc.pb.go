// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: tts/v2/tts.proto

package v2

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
	CloudMindsTTS_Call_FullMethodName               = "/ttsschema.CloudMindsTTS/Call"
	CloudMindsTTS_GetVersion_FullMethodName         = "/ttsschema.CloudMindsTTS/GetVersion"
	CloudMindsTTS_GetTtsConfig_FullMethodName       = "/ttsschema.CloudMindsTTS/GetTtsConfig"
	CloudMindsTTS_GetUserSpeakers_FullMethodName    = "/ttsschema.CloudMindsTTS/GetUserSpeakers"
	CloudMindsTTS_GetTtsConfigByUser_FullMethodName = "/ttsschema.CloudMindsTTS/GetTtsConfigByUser"
	CloudMindsTTS_Register_FullMethodName           = "/ttsschema.CloudMindsTTS/Register"
)

// CloudMindsTTSClient is the client API for CloudMindsTTS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudMindsTTSClient interface {
	// 合成音频流的流式接口
	Call(ctx context.Context, in *TtsReq, opts ...grpc.CallOption) (CloudMindsTTS_CallClient, error)
	// 获取服务版本信息
	GetVersion(ctx context.Context, in *VerVersionReq, opts ...grpc.CallOption) (*VerVersionRsp, error)
	// 获取服务端配置信息
	GetTtsConfig(ctx context.Context, in *VerReq, opts ...grpc.CallOption) (*RespGetTtsConfig, error)
	// 获取指定用户发音人信息
	GetUserSpeakers(ctx context.Context, in *GetUserSpeakersRequest, opts ...grpc.CallOption) (*GetUserSpeakersResponse, error)
	GetTtsConfigByUser(ctx context.Context, in *GetTtsConfigByUserRequest, opts ...grpc.CallOption) (*RespGetTtsConfig, error)
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
}

type cloudMindsTTSClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudMindsTTSClient(cc grpc.ClientConnInterface) CloudMindsTTSClient {
	return &cloudMindsTTSClient{cc}
}

func (c *cloudMindsTTSClient) Call(ctx context.Context, in *TtsReq, opts ...grpc.CallOption) (CloudMindsTTS_CallClient, error) {
	stream, err := c.cc.NewStream(ctx, &CloudMindsTTS_ServiceDesc.Streams[0], CloudMindsTTS_Call_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &cloudMindsTTSCallClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CloudMindsTTS_CallClient interface {
	Recv() (*TtsRes, error)
	grpc.ClientStream
}

type cloudMindsTTSCallClient struct {
	grpc.ClientStream
}

func (x *cloudMindsTTSCallClient) Recv() (*TtsRes, error) {
	m := new(TtsRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cloudMindsTTSClient) GetVersion(ctx context.Context, in *VerVersionReq, opts ...grpc.CallOption) (*VerVersionRsp, error) {
	out := new(VerVersionRsp)
	err := c.cc.Invoke(ctx, CloudMindsTTS_GetVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudMindsTTSClient) GetTtsConfig(ctx context.Context, in *VerReq, opts ...grpc.CallOption) (*RespGetTtsConfig, error) {
	out := new(RespGetTtsConfig)
	err := c.cc.Invoke(ctx, CloudMindsTTS_GetTtsConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudMindsTTSClient) GetUserSpeakers(ctx context.Context, in *GetUserSpeakersRequest, opts ...grpc.CallOption) (*GetUserSpeakersResponse, error) {
	out := new(GetUserSpeakersResponse)
	err := c.cc.Invoke(ctx, CloudMindsTTS_GetUserSpeakers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudMindsTTSClient) GetTtsConfigByUser(ctx context.Context, in *GetTtsConfigByUserRequest, opts ...grpc.CallOption) (*RespGetTtsConfig, error) {
	out := new(RespGetTtsConfig)
	err := c.cc.Invoke(ctx, CloudMindsTTS_GetTtsConfigByUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudMindsTTSClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	out := new(RegisterResp)
	err := c.cc.Invoke(ctx, CloudMindsTTS_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudMindsTTSServer is the server API for CloudMindsTTS service.
// All implementations must embed UnimplementedCloudMindsTTSServer
// for forward compatibility
type CloudMindsTTSServer interface {
	// 合成音频流的流式接口
	Call(*TtsReq, CloudMindsTTS_CallServer) error
	// 获取服务版本信息
	GetVersion(context.Context, *VerVersionReq) (*VerVersionRsp, error)
	// 获取服务端配置信息
	GetTtsConfig(context.Context, *VerReq) (*RespGetTtsConfig, error)
	// 获取指定用户发音人信息
	GetUserSpeakers(context.Context, *GetUserSpeakersRequest) (*GetUserSpeakersResponse, error)
	GetTtsConfigByUser(context.Context, *GetTtsConfigByUserRequest) (*RespGetTtsConfig, error)
	Register(context.Context, *RegisterReq) (*RegisterResp, error)
	mustEmbedUnimplementedCloudMindsTTSServer()
}

// UnimplementedCloudMindsTTSServer must be embedded to have forward compatible implementations.
type UnimplementedCloudMindsTTSServer struct {
}

func (UnimplementedCloudMindsTTSServer) Call(*TtsReq, CloudMindsTTS_CallServer) error {
	return status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedCloudMindsTTSServer) GetVersion(context.Context, *VerVersionReq) (*VerVersionRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedCloudMindsTTSServer) GetTtsConfig(context.Context, *VerReq) (*RespGetTtsConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTtsConfig not implemented")
}
func (UnimplementedCloudMindsTTSServer) GetUserSpeakers(context.Context, *GetUserSpeakersRequest) (*GetUserSpeakersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSpeakers not implemented")
}
func (UnimplementedCloudMindsTTSServer) GetTtsConfigByUser(context.Context, *GetTtsConfigByUserRequest) (*RespGetTtsConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTtsConfigByUser not implemented")
}
func (UnimplementedCloudMindsTTSServer) Register(context.Context, *RegisterReq) (*RegisterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedCloudMindsTTSServer) mustEmbedUnimplementedCloudMindsTTSServer() {}

// UnsafeCloudMindsTTSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudMindsTTSServer will
// result in compilation errors.
type UnsafeCloudMindsTTSServer interface {
	mustEmbedUnimplementedCloudMindsTTSServer()
}

func RegisterCloudMindsTTSServer(s grpc.ServiceRegistrar, srv CloudMindsTTSServer) {
	s.RegisterService(&CloudMindsTTS_ServiceDesc, srv)
}

func _CloudMindsTTS_Call_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TtsReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CloudMindsTTSServer).Call(m, &cloudMindsTTSCallServer{stream})
}

type CloudMindsTTS_CallServer interface {
	Send(*TtsRes) error
	grpc.ServerStream
}

type cloudMindsTTSCallServer struct {
	grpc.ServerStream
}

func (x *cloudMindsTTSCallServer) Send(m *TtsRes) error {
	return x.ServerStream.SendMsg(m)
}

func _CloudMindsTTS_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerVersionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudMindsTTSServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudMindsTTS_GetVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudMindsTTSServer).GetVersion(ctx, req.(*VerVersionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudMindsTTS_GetTtsConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudMindsTTSServer).GetTtsConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudMindsTTS_GetTtsConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudMindsTTSServer).GetTtsConfig(ctx, req.(*VerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudMindsTTS_GetUserSpeakers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserSpeakersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudMindsTTSServer).GetUserSpeakers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudMindsTTS_GetUserSpeakers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudMindsTTSServer).GetUserSpeakers(ctx, req.(*GetUserSpeakersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudMindsTTS_GetTtsConfigByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTtsConfigByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudMindsTTSServer).GetTtsConfigByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudMindsTTS_GetTtsConfigByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudMindsTTSServer).GetTtsConfigByUser(ctx, req.(*GetTtsConfigByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudMindsTTS_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudMindsTTSServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudMindsTTS_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudMindsTTSServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CloudMindsTTS_ServiceDesc is the grpc.ServiceDesc for CloudMindsTTS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloudMindsTTS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ttsschema.CloudMindsTTS",
	HandlerType: (*CloudMindsTTSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _CloudMindsTTS_GetVersion_Handler,
		},
		{
			MethodName: "GetTtsConfig",
			Handler:    _CloudMindsTTS_GetTtsConfig_Handler,
		},
		{
			MethodName: "GetUserSpeakers",
			Handler:    _CloudMindsTTS_GetUserSpeakers_Handler,
		},
		{
			MethodName: "GetTtsConfigByUser",
			Handler:    _CloudMindsTTS_GetTtsConfigByUser_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _CloudMindsTTS_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Call",
			Handler:       _CloudMindsTTS_Call_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tts/v2/tts.proto",
}
