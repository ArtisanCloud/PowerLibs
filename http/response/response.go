package response

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	*http.Response
}

func NewHttpResponse() *HttpResponse {
	body := ""
	return &HttpResponse{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
			ContentLength: int64(len(body)),
			Request:       nil,
			Header:        make(http.Header, 20),
		},
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
