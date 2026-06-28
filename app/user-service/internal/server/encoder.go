package server

import (
	"net/http"

	kratos_http "github.com/go-kratos/kratos/v3/transport/http"
	kratos_status "github.com/go-kratos/kratos/v3/transport/http/status"
	"google.golang.org/grpc/status"
)

type httpResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(kratos_http.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := kratos_http.CodecForRequest(r, "Accept")
	resp := &httpResponse{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    v,
	}
	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, err = w.Write(data)
	return err
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	resp := new(httpResponse)
	if gs, ok := status.FromError(err); ok {
		resp = &httpResponse{
			Success: false,
			Code:    kratos_status.FromGRPCCode(gs.Code()),
			Msg:     gs.Message(),
			Data:    nil,
		}
	} else {
		resp = &httpResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data:    nil,
		}
	}
	codec, _ := kratos_http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(resp.Code)
	_, _ = w.Write(body)
}
