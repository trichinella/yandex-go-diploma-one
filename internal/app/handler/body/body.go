package body

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

func Content(r *http.Request) ([]byte, error) {
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		return CompressedContent(r)
	}

	return UnCompressedContent(r)
}

func CompressedContent(r *http.Request) ([]byte, error) {
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return []byte{}, nil
		}

		return nil, err
	}

	defer func() {
		_ = gz.Close()
	}()

	// при чтении вернётся распакованный слайс байт
	body, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func UnCompressedContent(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = r.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, err
}
