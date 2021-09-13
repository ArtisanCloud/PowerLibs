package response

import (
	"io/ioutil"
	"os"
)

type StreamResponse struct {
	*HttpResponse
}

type RequestDownload struct {
	HashType    string `json:"hash_type""`
	HashValue   string `json:"hash_value""`
	DownloadURL string `json:"download_url""`
}

func (rs *StreamResponse) Save(directory string, fileName string ) (int, error) {

	saveFile, err := os.Create(directory)
	if err != nil {
		return 0, err
	}
	defer saveFile.Close()

	data, err := ioutil.ReadAll(rs.GetBody())

	//fileMd5 := sha256.New()
	totalSize := 0
	_, err = saveFile.Write(data)
	if err != nil {
		return 0, err
	}

	//fileMd5.Write(data)
	//totalSize += len(data)

	//if totalSize != d.fileSize {
	//	return errors.New("文件不完整")
	//}
	//
	//if d.md5 != "" {
	//	if hex.EncodeToString(fileMd5.Sum(nil)) != d.md5 {
	//		return errors.New("文件损坏")
	//	} else {
	//		log.Println("文件SHA-256校验成功")
	//	}
	//}
	//
	//if d.md5 != "" {
	//	if hex.EncodeToString(fileMd5.Sum(nil)) != d.md5 {
	//		return 0, errors.New("文件损坏")
	//	} else {
	//		log.Println("文件SHA-256校验成功")
	//	}
	//}

	return totalSize, nil
}

func (rs *StreamResponse) SaveAs(directory string, fileName string, appendSuffix bool) (int, error) {
	return rs.Save(directory, fileName)
}
