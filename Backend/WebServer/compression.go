package webserver

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var compressibleTypes = map[string]bool{
	"text/html":               true,
	"text/css":                true,
	"text/javascript":         true,
	"text/plain":              true,
	"text/xml":                true,
	"text/csv":                true,
	"application/javascript":  true,
	"application/json":        true,
	"application/xml":         true,
	"application/wasm":        true,
	"application/xhtml+xml":   true,
	"application/rss+xml":     true,
	"application/atom+xml":    true,
	"application/manifest+json": true,
	"image/svg+xml":           true,
}

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
		return w
	},
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzWriter      *gzip.Writer
	headerWritten bool
	compressed    bool
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	if !g.headerWritten {
		g.detectAndSetHeaders()
	}
	if g.compressed {
		return g.gzWriter.Write(data)
	}
	return g.ResponseWriter.Write(data)
}

func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	if !g.headerWritten {
		g.detectAndSetHeaders()
	}
	g.ResponseWriter.WriteHeader(statusCode)
}

func (g *gzipResponseWriter) detectAndSetHeaders() {
	g.headerWritten = true
	contentType := g.ResponseWriter.Header().Get("Content-Type")
	if contentType == "" {
		return
	}

	mimeType := strings.SplitN(contentType, ";", 2)[0]
	mimeType = strings.TrimSpace(mimeType)

	if compressibleTypes[mimeType] {
		g.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		g.ResponseWriter.Header().Add("Vary", "Accept-Encoding")
		g.ResponseWriter.Header().Del("Content-Length")
		g.compressed = true
	}
}

func (g *gzipResponseWriter) Flush() {
	if g.compressed {
		g.gzWriter.Flush()
	}
	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (g *gzipResponseWriter) close() {
	if g.compressed {
		g.gzWriter.Close()
		gzipWriterPool.Put(g.gzWriter)
	}
}

func clientAcceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func isAlreadyCompressed(header http.Header) bool {
	return header.Get("Content-Encoding") != ""
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	gz := gzipWriterPool.Get().(*gzip.Writer)
	gz.Reset(w)
	return &gzipResponseWriter{
		ResponseWriter: w,
		gzWriter:       gz,
	}
}
