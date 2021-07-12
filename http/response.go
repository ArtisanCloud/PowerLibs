package http

import (
	"io"
	"net/http"
)

type HttpResponse struct {
	*http.Response
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		Response: &http.Response{},
	}
}

func (rs HttpResponse) GetBody() io.ReadCloser {
	return rs.Response.Body
}
func (rs HttpResponse) GetHeader() http.Header {
	return rs.Response.Header
}
func (rs HttpResponse) GetStatusCode() int {
	return rs.Response.StatusCode
}
