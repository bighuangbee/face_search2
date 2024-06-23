// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.0
// source: biz/v1/face_recognize.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationFaceRecognizeFaceSearchByDatetime = "/api.biz.v1.FaceRecognize/FaceSearchByDatetime"
const OperationFaceRecognizeRegisteByPath = "/api.biz.v1.FaceRecognize/RegisteByPath"
const OperationFaceRecognizeRegisteStatus = "/api.biz.v1.FaceRecognize/RegisteStatus"
const OperationFaceRecognizeUnRegisteAll = "/api.biz.v1.FaceRecognize/UnRegisteAll"

type FaceRecognizeHTTPServer interface {
	// FaceSearchByDatetime人脸搜索-按时间日期范围
	FaceSearchByDatetime(context.Context, *FaceSearchByDatetimeRequest) (*SearchResultReply, error)
	// RegisteByPath人脸注册-从默认目录读取注册图
	RegisteByPath(context.Context, *EmptyRequest) (*RegisteByPathReply, error)
	// RegisteStatus人脸注册-获取状态
	RegisteStatus(context.Context, *EmptyRequest) (*RegisteStatusReply, error)
	// UnRegisteAll人脸注销-所有人脸
	UnRegisteAll(context.Context, *EmptyRequest) (*EmptyReply, error)
}

func RegisterFaceRecognizeHTTPServer(s *http.Server, srv FaceRecognizeHTTPServer) {
	r := s.Route("/")
	r.POST("/face/registe/path", _FaceRecognize_RegisteByPath0_HTTP_Handler(srv))
	r.GET("/face/registe/status", _FaceRecognize_RegisteStatus0_HTTP_Handler(srv))
	r.POST("/face/search/datetime", _FaceRecognize_FaceSearchByDatetime0_HTTP_Handler(srv))
	r.POST("/face/unregiste/all", _FaceRecognize_UnRegisteAll0_HTTP_Handler(srv))
}

func _FaceRecognize_RegisteByPath0_HTTP_Handler(srv FaceRecognizeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in EmptyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFaceRecognizeRegisteByPath)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RegisteByPath(ctx, req.(*EmptyRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisteByPathReply)
		return ctx.Result(200, reply)
	}
}

func _FaceRecognize_RegisteStatus0_HTTP_Handler(srv FaceRecognizeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in EmptyRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFaceRecognizeRegisteStatus)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RegisteStatus(ctx, req.(*EmptyRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisteStatusReply)
		return ctx.Result(200, reply)
	}
}

func _FaceRecognize_FaceSearchByDatetime0_HTTP_Handler(srv FaceRecognizeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FaceSearchByDatetimeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFaceRecognizeFaceSearchByDatetime)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FaceSearchByDatetime(ctx, req.(*FaceSearchByDatetimeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SearchResultReply)
		return ctx.Result(200, reply)
	}
}

func _FaceRecognize_UnRegisteAll0_HTTP_Handler(srv FaceRecognizeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in EmptyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFaceRecognizeUnRegisteAll)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UnRegisteAll(ctx, req.(*EmptyRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*EmptyReply)
		return ctx.Result(200, reply)
	}
}

type FaceRecognizeHTTPClient interface {
	FaceSearchByDatetime(ctx context.Context, req *FaceSearchByDatetimeRequest, opts ...http.CallOption) (rsp *SearchResultReply, err error)
	RegisteByPath(ctx context.Context, req *EmptyRequest, opts ...http.CallOption) (rsp *RegisteByPathReply, err error)
	RegisteStatus(ctx context.Context, req *EmptyRequest, opts ...http.CallOption) (rsp *RegisteStatusReply, err error)
	UnRegisteAll(ctx context.Context, req *EmptyRequest, opts ...http.CallOption) (rsp *EmptyReply, err error)
}

type FaceRecognizeHTTPClientImpl struct {
	cc *http.Client
}

func NewFaceRecognizeHTTPClient(client *http.Client) FaceRecognizeHTTPClient {
	return &FaceRecognizeHTTPClientImpl{client}
}

func (c *FaceRecognizeHTTPClientImpl) FaceSearchByDatetime(ctx context.Context, in *FaceSearchByDatetimeRequest, opts ...http.CallOption) (*SearchResultReply, error) {
	var out SearchResultReply
	pattern := "/face/search/datetime"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationFaceRecognizeFaceSearchByDatetime))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FaceRecognizeHTTPClientImpl) RegisteByPath(ctx context.Context, in *EmptyRequest, opts ...http.CallOption) (*RegisteByPathReply, error) {
	var out RegisteByPathReply
	pattern := "/face/registe/path"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationFaceRecognizeRegisteByPath))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FaceRecognizeHTTPClientImpl) RegisteStatus(ctx context.Context, in *EmptyRequest, opts ...http.CallOption) (*RegisteStatusReply, error) {
	var out RegisteStatusReply
	pattern := "/face/registe/status"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationFaceRecognizeRegisteStatus))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FaceRecognizeHTTPClientImpl) UnRegisteAll(ctx context.Context, in *EmptyRequest, opts ...http.CallOption) (*EmptyReply, error) {
	var out EmptyReply
	pattern := "/face/unregiste/all"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationFaceRecognizeUnRegisteAll))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
