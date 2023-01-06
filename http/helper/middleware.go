package helper

import (
	"bytes"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func HttpDebugMiddleware(debug bool) contract.RequestMiddleware {
	return func(handle contract.RequestHandle) contract.RequestHandle {
		return func(request *http.Request) (response *http.Response, err error) {
			var output bytes.Buffer
			if debug {
				output.WriteString(fmt.Sprintf("%s %s ", request.Method, request.URL.String()))

				// print out request header
				output.Write([]byte("\r\nrequest header: { \r\n"))
				for k, vals := range request.Header {
					for _, v := range vals {
						output.Write([]byte(fmt.Sprintf("\t%s:%s\r\n", k, v)))
					}
				}
				output.Write([]byte("} \r\n"))

				// print out request body
				if request.Body != nil {

					output.Write([]byte("request body:"))
					var buf bytes.Buffer
					reader := io.TeeReader(request.Body, &buf)
					body, _ := io.ReadAll(reader)
					request.Body = io.NopCloser(&buf)
					output.Write(body)
				}
			}

			response, err = handle(request)
			if err != nil {
				return response, err
			}

			if debug {

				output.Write([]byte("\r\n------------------\r\n"))
				output.Write([]byte("response content:\r\n"))
				dumpRes, _ := httputil.DumpResponse(response, true)
				output.Write(dumpRes)
				log.Println(output.String())
			}

			return
		}
	}
}
