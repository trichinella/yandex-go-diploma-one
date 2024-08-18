package middleware

import (
	"compress/gzip"
	"diploma1/internal/app/service/logging"
	"io"
	"net/http"
	"slices"
	"strings"
)

var compressibleContentTypes []string

func init() {
	compressibleContentTypes = []string{"application/json", "text/html", "text/plain"}
}

func Compress() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !slices.Contains(compressibleContentTypes, r.Header.Get("Content-Type")) {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := getCompressor(w)
			defer func() {
				if err = gz.Close(); err != nil {
					logging.Sugar.Error(err)
				}
			}()

			if err != nil {
				logging.Sugar.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		})
	}
}

func getCompressor(w http.ResponseWriter) (*gzip.Writer, error) {
	gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)

	if err != nil {
		return nil, err
	}

	return gz, nil
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
