package req

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var client http.Client

func Request(method, url string, header map[string]any, body io.Reader) (resp []byte, err error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		slog.Error("http.NewRequest error" + err.Error())
		return
	}

	for k, v := range header {
		switch v := v.(type) {
		case string:
			req.Header.Set(k, v)
		case []string:
			for _, item := range v {
				req.Header.Add(k, item)
			}
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("http request response status error code: %d", res.StatusCode)
		slog.Error(err.Error())
		return
	}

	resp, err = io.ReadAll(res.Body)

	slog.Debug("Request response status",
		"method", method, "url", url,
		"header", req.Header,
		"response body", resp,
	)

	return
}

func Get(url string, header map[string]any, body io.Reader) ([]byte, error) {
	return Request("GET", url, header, body)
}

func Post(url string, header map[string]any, body io.Reader) ([]byte, error) {
	return Request("POST", url, header, body)
}
