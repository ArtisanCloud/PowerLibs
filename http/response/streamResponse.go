package response

import (
	"os"
	"path/filepath"
)

type StreamResponse struct {
	*HttpResponse
	Data []byte
}

type RequestDownload struct {
	HashType    string `json:"hash_type""`
	HashValue   string `json:"hash_value""`
	DownloadURL string `json:"download_url""`
}


func NewStreamResponse(code int) *StreamResponse {

	return &StreamResponse{
		NewHttpResponse(code),
		nil,
	}
}

func (rs *StreamResponse)SetHttpResponse(httpResponse *HttpResponse){
	rs.HttpResponse = httpResponse
}


func (rs *StreamResponse) Save(directory string, fileName string ) (int, error) {

	path := filepath.Join(directory, fileName)

	saveFile, err := os.Create(path)
	if err != nil {
		return 0,err
	}
	defer saveFile.Close()

	totalSize, err := saveFile.Write(rs.Data)
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

func (rs *StreamResponse) SaveAs(directory string, fileName string) (int, error) {
	return rs.Save(directory, fileName)
}
