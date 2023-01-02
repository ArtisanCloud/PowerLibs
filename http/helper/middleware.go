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
				if request.Body != nil {
					var buf bytes.Buffer
					reader := io.TeeReader(request.Body, &buf)
					body, _ := io.ReadAll(reader)
					request.Body = io.NopCloser(&buf)
					output.Write(body)
				}
			}

			response, err = handle(request)

			if debug {
				dumpRes, _ := httputil.DumpResponse(response, true)
				output.Write(dumpRes)
				log.Println(output.String())
			}

			return
		}
	}
}
