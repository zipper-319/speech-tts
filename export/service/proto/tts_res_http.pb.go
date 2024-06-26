// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.0
// - protoc             v3.15.8
// source: ttsData/v2/tts_res.proto

package ttsData

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationTtsDataAddTtsData = "/tts_data.v2.TtsData/AddTtsData"
const OperationTtsDataDelTtsData = "/tts_data.v2.TtsData/DelTtsData"
const OperationTtsDataGeneToken = "/tts_data.v2.TtsData/GeneToken"
const OperationTtsDataGetAllResource = "/tts_data.v2.TtsData/GetAllResource"
const OperationTtsDataGetSpeakerModel = "/tts_data.v2.TtsData/GetSpeakerModel"
const OperationTtsDataGetTtsData = "/tts_data.v2.TtsData/GetTtsData"
const OperationTtsDataRegisterResService = "/tts_data.v2.TtsData/RegisterResService"
const OperationTtsDataUnRegisterResService = "/tts_data.v2.TtsData/UnRegisterResService"
const OperationTtsDataUpdateTtsData = "/tts_data.v2.TtsData/UpdateTtsData"

type TtsDataHTTPServer interface {
	AddTtsData(context.Context, *AddTtsDataRequest) (*emptypb.Empty, error)
	DelTtsData(context.Context, *DelTtsDataRequest) (*emptypb.Empty, error)
	GeneToken(context.Context, *GeneTokenRequest) (*GeneTokenResponse, error)
	GetAllResource(context.Context, *emptypb.Empty) (*GetAllResourceResult, error)
	GetSpeakerModel(context.Context, *emptypb.Empty) (*GetSpeakerModelResult, error)
	GetTtsData(context.Context, *GetTtsDataRequest) (*GetTtsDataResponse, error)
	RegisterResService(context.Context, *RegisterResServiceRequest) (*emptypb.Empty, error)
	UnRegisterResService(context.Context, *UnRegisterResServiceRequest) (*emptypb.Empty, error)
	UpdateTtsData(context.Context, *UpdateTtsDataRequest) (*emptypb.Empty, error)
}

func RegisterTtsDataHTTPServer(s *http.Server, srv TtsDataHTTPServer) {
	r := s.Route("/")
	r.GET("/api/ttsData/v1/resource/get", _TtsData_GetTtsData0_HTTP_Handler(srv))
	r.POST("/api/ttsData/v1/resource/add", _TtsData_AddTtsData0_HTTP_Handler(srv))
	r.POST("/api/ttsData/v1/resource/del", _TtsData_DelTtsData0_HTTP_Handler(srv))
	r.POST("/api/ttsData/v1/resource/update", _TtsData_UpdateTtsData0_HTTP_Handler(srv))
	r.GET("/api/ttsData/v1/resource/geneToken", _TtsData_GeneToken0_HTTP_Handler(srv))
	r.GET("/api/ttsData/v1/resource/get-all", _TtsData_GetAllResource0_HTTP_Handler(srv))
	r.GET("/api/ttsData/v1/resource/get-speaker-model", _TtsData_GetSpeakerModel0_HTTP_Handler(srv))
	r.POST("/api/ttsData/v1/resource/register", _TtsData_RegisterResService0_HTTP_Handler(srv))
	r.POST("/api/ttsData/v1/resource/unregister", _TtsData_UnRegisterResService0_HTTP_Handler(srv))
}

func _TtsData_GetTtsData0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetTtsDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataGetTtsData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetTtsData(ctx, req.(*GetTtsDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetTtsDataResponse)
		return ctx.Result(200, reply)
	}
}

func _TtsData_AddTtsData0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AddTtsDataRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataAddTtsData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddTtsData(ctx, req.(*AddTtsDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _TtsData_DelTtsData0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DelTtsDataRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataDelTtsData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DelTtsData(ctx, req.(*DelTtsDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _TtsData_UpdateTtsData0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateTtsDataRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataUpdateTtsData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateTtsData(ctx, req.(*UpdateTtsDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _TtsData_GeneToken0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GeneTokenRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataGeneToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GeneToken(ctx, req.(*GeneTokenRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GeneTokenResponse)
		return ctx.Result(200, reply)
	}
}

func _TtsData_GetAllResource0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataGetAllResource)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAllResource(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAllResourceResult)
		return ctx.Result(200, reply)
	}
}

func _TtsData_GetSpeakerModel0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataGetSpeakerModel)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetSpeakerModel(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetSpeakerModelResult)
		return ctx.Result(200, reply)
	}
}

func _TtsData_RegisterResService0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterResServiceRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataRegisterResService)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RegisterResService(ctx, req.(*RegisterResServiceRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _TtsData_UnRegisterResService0_HTTP_Handler(srv TtsDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UnRegisterResServiceRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationTtsDataUnRegisterResService)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UnRegisterResService(ctx, req.(*UnRegisterResServiceRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type TtsDataHTTPClient interface {
	AddTtsData(ctx context.Context, req *AddTtsDataRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DelTtsData(ctx context.Context, req *DelTtsDataRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GeneToken(ctx context.Context, req *GeneTokenRequest, opts ...http.CallOption) (rsp *GeneTokenResponse, err error)
	GetAllResource(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *GetAllResourceResult, err error)
	GetSpeakerModel(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *GetSpeakerModelResult, err error)
	GetTtsData(ctx context.Context, req *GetTtsDataRequest, opts ...http.CallOption) (rsp *GetTtsDataResponse, err error)
	RegisterResService(ctx context.Context, req *RegisterResServiceRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UnRegisterResService(ctx context.Context, req *UnRegisterResServiceRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateTtsData(ctx context.Context, req *UpdateTtsDataRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type TtsDataHTTPClientImpl struct {
	cc *http.Client
}

func NewTtsDataHTTPClient(client *http.Client) TtsDataHTTPClient {
	return &TtsDataHTTPClientImpl{client}
}

func (c *TtsDataHTTPClientImpl) AddTtsData(ctx context.Context, in *AddTtsDataRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/ttsData/v1/resource/add"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationTtsDataAddTtsData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) DelTtsData(ctx context.Context, in *DelTtsDataRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/ttsData/v1/resource/del"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationTtsDataDelTtsData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) GeneToken(ctx context.Context, in *GeneTokenRequest, opts ...http.CallOption) (*GeneTokenResponse, error) {
	var out GeneTokenResponse
	pattern := "/api/ttsData/v1/resource/geneToken"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTtsDataGeneToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) GetAllResource(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*GetAllResourceResult, error) {
	var out GetAllResourceResult
	pattern := "/api/ttsData/v1/resource/get-all"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTtsDataGetAllResource))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) GetSpeakerModel(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*GetSpeakerModelResult, error) {
	var out GetSpeakerModelResult
	pattern := "/api/ttsData/v1/resource/get-speaker-model"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTtsDataGetSpeakerModel))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) GetTtsData(ctx context.Context, in *GetTtsDataRequest, opts ...http.CallOption) (*GetTtsDataResponse, error) {
	var out GetTtsDataResponse
	pattern := "/api/ttsData/v1/resource/get"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationTtsDataGetTtsData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) RegisterResService(ctx context.Context, in *RegisterResServiceRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/ttsData/v1/resource/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationTtsDataRegisterResService))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) UnRegisterResService(ctx context.Context, in *UnRegisterResServiceRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/ttsData/v1/resource/unregister"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationTtsDataUnRegisterResService))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TtsDataHTTPClientImpl) UpdateTtsData(ctx context.Context, in *UpdateTtsDataRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/ttsData/v1/resource/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationTtsDataUpdateTtsData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
