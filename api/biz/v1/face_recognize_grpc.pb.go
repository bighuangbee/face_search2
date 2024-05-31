// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: biz/v1/face_recognize.proto

package v1

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
	FaceRecognize_RegisteByPath_FullMethodName = "/api.biz.v1.FaceRecognize/RegisteByPath"
	FaceRecognize_UnRegisteAll_FullMethodName  = "/api.biz.v1.FaceRecognize/UnRegisteAll"
)

// FaceRecognizeClient is the client API for FaceRecognize service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FaceRecognizeClient interface {
	// 人脸注册-从默认目录读取注册图
	RegisteByPath(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyReply, error)
	// 人脸注销-所有人脸
	UnRegisteAll(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyReply, error)
}

type faceRecognizeClient struct {
	cc grpc.ClientConnInterface
}

func NewFaceRecognizeClient(cc grpc.ClientConnInterface) FaceRecognizeClient {
	return &faceRecognizeClient{cc}
}

func (c *faceRecognizeClient) RegisteByPath(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, FaceRecognize_RegisteByPath_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *faceRecognizeClient) UnRegisteAll(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, FaceRecognize_UnRegisteAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FaceRecognizeServer is the server API for FaceRecognize service.
// All implementations must embed UnimplementedFaceRecognizeServer
// for forward compatibility
type FaceRecognizeServer interface {
	// 人脸注册-从默认目录读取注册图
	RegisteByPath(context.Context, *EmptyRequest) (*EmptyReply, error)
	// 人脸注销-所有人脸
	UnRegisteAll(context.Context, *EmptyRequest) (*EmptyReply, error)
	mustEmbedUnimplementedFaceRecognizeServer()
}

// UnimplementedFaceRecognizeServer must be embedded to have forward compatible implementations.
type UnimplementedFaceRecognizeServer struct {
}

func (UnimplementedFaceRecognizeServer) RegisteByPath(context.Context, *EmptyRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisteByPath not implemented")
}
func (UnimplementedFaceRecognizeServer) UnRegisteAll(context.Context, *EmptyRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnRegisteAll not implemented")
}
func (UnimplementedFaceRecognizeServer) mustEmbedUnimplementedFaceRecognizeServer() {}

// UnsafeFaceRecognizeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FaceRecognizeServer will
// result in compilation errors.
type UnsafeFaceRecognizeServer interface {
	mustEmbedUnimplementedFaceRecognizeServer()
}

func RegisterFaceRecognizeServer(s grpc.ServiceRegistrar, srv FaceRecognizeServer) {
	s.RegisterService(&FaceRecognize_ServiceDesc, srv)
}

func _FaceRecognize_RegisteByPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaceRecognizeServer).RegisteByPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaceRecognize_RegisteByPath_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaceRecognizeServer).RegisteByPath(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FaceRecognize_UnRegisteAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FaceRecognizeServer).UnRegisteAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FaceRecognize_UnRegisteAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FaceRecognizeServer).UnRegisteAll(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FaceRecognize_ServiceDesc is the grpc.ServiceDesc for FaceRecognize service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FaceRecognize_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.biz.v1.FaceRecognize",
	HandlerType: (*FaceRecognizeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisteByPath",
			Handler:    _FaceRecognize_RegisteByPath_Handler,
		},
		{
			MethodName: "UnRegisteAll",
			Handler:    _FaceRecognize_UnRegisteAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "biz/v1/face_recognize.proto",
}