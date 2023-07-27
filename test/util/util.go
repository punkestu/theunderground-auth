package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
