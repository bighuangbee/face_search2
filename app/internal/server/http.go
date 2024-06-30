package server

import (
	"context"
	"encoding/json"
	"github.com/bighuangbee/face_search2/app/internal/service/face"
	"github.com/bighuangbee/face_search2/pkg/middleware"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	pb "github.com/bighuangbee/face_search2/api/biz/v1"
	"github.com/bighuangbee/face_search2/pkg/conf"
)

func jsonMarshal(res *pb.CommonReply) ([]byte, error) {
	newProto := protojson.MarshalOptions{EmitUnpopulated: true}
	output, err := newProto.Marshal(res)
	if err != nil {
		return nil, err
	}

	var stuff map[string]any
	if err := json.Unmarshal(output, &stuff); err != nil {
		return nil, err
	}

	if stuff["data"] != nil {
		delete(stuff["data"].(map[string]any), "@type")
	}
	return json.MarshalIndent(stuff, "", "  ")
}

func EncoderResponseSuccess() http.ServerOption {
	return http.ResponseEncoder(
		func(w stdhttp.ResponseWriter, request *stdhttp.Request, i interface{}) error {
			resp := &pb.CommonReply{
				Code:    200,
				Message: "",
			}
			var data []byte
			var err error
			if m, ok := i.(proto.Message); ok {
				payload, err := anypb.New(m)
				if err != nil {
					return err
				}
				resp.Data = payload
				data, err = jsonMarshal(resp)
				if err != nil {
					return err
				}
			} else {
				dataMap := map[string]interface{}{
					"code":    200,
					"message": "",
					"data":    i,
				}
				data, err = json.Marshal(dataMap)
				if err != nil {
					return err
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(data)
			if err != nil {
				return err
			}
			return nil
		})
}

func EncoderResponseError() http.ServerOption {
	return http.ErrorEncoder(func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
		se := errors.FromError(err)

		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(&pb.CommonReply{
			Message: se.Message,
			Code:    se.Code,
			Reason:  se.Reason,
		})
		if err != nil {
			w.WriteHeader(stdhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		code := stdhttp.StatusOK
		if se.Code == stdhttp.StatusUnauthorized {
			code = stdhttp.StatusUnauthorized
		}
		w.WriteHeader(code)
		w.Write(body)
		return
	})
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	logger log.Logger,
	faceRecognize *face.FaceRecognizeApp,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.LogFile(logger),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}),
			handlers.AllowCredentials(),
		)),
		EncoderResponseSuccess(),
		EncoderResponseError(),
	}

	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	pb.RegisterFaceRecognizeHTTPServer(srv, faceRecognize)
	srv.Handle("/debug/pprof/", pprof.NewHandler())

	route := srv.Route("/")

	route.POST("/face/search", func(ctx http.Context) error {
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return faceRecognize.Search(ctx)
		})

		reply, err := h(ctx, nil)
		if err != nil {
			return err
		}

		result := reply.(*pb.SearchResultReply)

		resp := &pb.CommonReply{}
		if err != nil {
			resp.Code = int32(errors.Code(err))
			resp.Reason = errors.Reason(err)
			resp.Message = "人脸库搜索失败, " + err.Error()
		} else {
			resp.Code = 200
			resp.Message = "人脸搜索成功"
			resp.Data, _ = anypb.New(result)
		}
		d, _ := jsonMarshal(resp)
		ctx.Response().Write(d)
		return nil
	})

	return srv
}
