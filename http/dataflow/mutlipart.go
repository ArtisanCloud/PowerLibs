package dataflow

import (
	"bytes"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
)

type MultipartDf struct {
	buf     bytes.Buffer
	mWriter *multipart.Writer
	errs    []error
}

func NewMultipartHelper() contract.MultipartDfInterface {
	df := MultipartDf{}
	mWriter := multipart.NewWriter(&df.buf)
	df.mWriter = mWriter
	return &df
}

func (m *MultipartDf) Boundary(b string) contract.MultipartDfInterface {
	err := m.mWriter.SetBoundary(b)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "set boundary failed"))
	}
	return m
}

func (m *MultipartDf) FileByPath(fieldName string, filePath string) contract.MultipartDfInterface {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	defer file.Close()
	_, fileName := path.Split(filePath)

	writer, err := m.mWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	return m
}

func (m *MultipartDf) FileMem(fieldName string, fileName string, reader io.Reader) contract.MultipartDfInterface {
	writer, err := m.mWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create file part failed"))
	}
	return m
}

func (m *MultipartDf) Part(header textproto.MIMEHeader, reader io.Reader) contract.MultipartDfInterface {
	writer, err := m.mWriter.CreatePart(header)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create part failed"))
		return m
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create part read failed"))
		return m
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create part write failed"))
		return m
	}
	return m
}

func (m *MultipartDf) FieldValue(fieldName string, value string) contract.MultipartDfInterface {
	err := m.mWriter.WriteField(fieldName, value)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "set field failed"))
	}
	return m
}

func (m *MultipartDf) Field(fieldName string, reader io.Reader) contract.MultipartDfInterface {
	writer, err := m.mWriter.CreateFormField(fieldName)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create field failed"))
		return m
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create field read failed"))
		return m
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		m.errs = append(m.errs, errors.Wrap(err, "create field write failed"))
		return m
	}
	return m
}

func (m *MultipartDf) Close() error {
	return m.mWriter.Close()
}

func (m *MultipartDf) GetBoundary() string {
	return m.mWriter.Boundary()
}

func (m *MultipartDf) GetReader() io.Reader {
	return &m.buf
}

func (m *MultipartDf) GetContentType() string {
	return m.mWriter.FormDataContentType()
}

func (m *MultipartDf) Err() error {
	if len(m.errs) == 0 {
		return nil
	} else {
		return m.errs[0]
	}
}
