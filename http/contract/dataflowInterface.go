package contract

import (
	"context"
	"io"
	"net/http"
	"net/textproto"
)

// RequestDataflowInterface 是一个 Http 请求构建器, 建议将注释中的私有方法实现到内部
type RequestDataflowInterface interface {
	// 获取引用的 Client
	// getClient() ClientInterface
	// 获取引用的 Round Request Handle
	// getMiddlewareHandle() RequestHandle

	WithContext(ctx context.Context) RequestDataflowInterface
	Method(method string) RequestDataflowInterface
	Uri(uri string) RequestDataflowInterface
	Url(url string) RequestDataflowInterface
	Header(key string, values ...string) RequestDataflowInterface
	Query(key string, values ...string) RequestDataflowInterface

	Json(jsonAny interface{}) RequestDataflowInterface
	Body(body io.Reader) RequestDataflowInterface
	Any(data BodyEncoder) RequestDataflowInterface
	Xml(xmlAny interface{}) RequestDataflowInterface
	Multipart(multipartDf func(multipart MultipartDfInterface)) RequestDataflowInterface

	Err() error

	// 在发送前应该检查错误
	// validateRequest() error

	Request() (response *http.Response, err error)
	Result(result interface{}) (err error)
	RequestResHelper() (response ResponseHelper, err error)
}

type BodyEncoder interface {
	Encode() (io.Reader, error)
}

type ResponseHelper interface {
	GetStatusCode() int
	GetHeader(key string) string
	GetBody() io.Reader
	GetBodyBytes() ([]byte, error)
	GetBodyJsonAsMap() (map[string]interface{}, error)
}

type MultipartDfInterface interface {
	Boundary(b string) MultipartDfInterface
	FileByPath(fieldName string, filePath string) MultipartDfInterface
	FileMem(fieldName string, fileName string, reader io.Reader) MultipartDfInterface
	Part(header textproto.MIMEHeader, reader io.Reader) MultipartDfInterface
	FieldValue(fieldName string, value string) MultipartDfInterface
	Field(fieldName string, reader io.Reader) MultipartDfInterface
	Close() error
	GetBoundary() string
	GetReader() io.Reader
	GetContentType() string
	Err() error
}
