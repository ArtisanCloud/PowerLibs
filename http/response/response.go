package response

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	*http.Response
}

func NewHttpResponse(code int) *HttpResponse {
	body := ""
	return &HttpResponse{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    code,
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

func (rs HttpResponse) GetBodyData() ([]byte, error) {
	body, err := ioutil.ReadAll(rs.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("read Response body err: %v", err)
	}

	_ = rs.Response.Body.Close()
	rs.Response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (rs HttpResponse) GetHeader() http.Header {
	return rs.Response.Header
}
func (rs HttpResponse) GetStatusCode() int {
	return rs.Response.StatusCode
}

func (rs HttpResponse) Send(writer http.ResponseWriter) (err error) {

	// set header code
	writer.WriteHeader(rs.GetStatusCode())

	// set write body
	body, err := ioutil.ReadAll(rs.GetBody())
	if err != nil {
		return err
	}
	_, err = writer.Write(body)

	return err
}
