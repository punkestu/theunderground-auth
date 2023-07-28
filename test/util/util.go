package util

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

type GetBodyOptions struct {
	Verbose bool
}

func GetBody(resp *http.Response, dst any, opt *GetBodyOptions) error {
	resBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return err
	}
	if opt != nil {
		if opt.Verbose {
			println(string(resBodyByte))
		}
	}
	err = json.Unmarshal(resBodyByte, &dst)
	if err != nil {
		println(err.Error())
	}
	return nil
}

func SendRequest(method, endpoint string, body any, headers map[string]string) (req *http.Request, err error) {
	reqBody, err := json.Marshal(&body)
	req = httptest.NewRequest(method, endpoint, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	for name, header := range headers {
		req.Header.Set(name, header)
	}
	return
}

func SendFileRequest(method, endpoint, filePath string, headers map[string]string) (req *http.Request, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	part1, err := writer.CreateFormFile("credential", filepath.Base(pwd+filePath))
	if err != nil {
		return
	}
	_, err = io.Copy(part1, file)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req = httptest.NewRequest(method, endpoint, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for name, header := range headers {
		req.Header.Set(name, header)
	}

	return
}
