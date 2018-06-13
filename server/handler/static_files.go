package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/src-d/gitbase-playground/server/assets"
)

const (
	staticDirName = "static"
	indexFileName = "index.html"

	serverValuesPlaceholder = "window.REPLACE_BY_SERVER"
)

// Static contains handlers to serve static using go-bindata
type Static struct {
	dir     string
	options options
}

// NewStatic creates new Static
func NewStatic(dir, serverURL string, selectLimit int) *Static {
	return &Static{
		dir: dir,
		options: options{
			ServerURL:   serverURL,
			SelectLimit: selectLimit,
		},
	}
}

// struct which will be marshalled and exposed to frontend
type options struct {
	ServerURL   string `json:"SERVER_URL"`
	SelectLimit int    `json:"SELECT_LIMIT"`
}

// ServeHTTP serves any static file from static directory or fallbacks on index.hml
func (s *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join(s.dir, r.URL.Path)
	b, err := assets.Asset(filepath)
	if err != nil {
		if strings.HasPrefix(filepath, path.Join(s.dir, staticDirName)) {
			http.NotFound(w, r)
			return
		}

		s.ServeIndexHTML(nil)(w, r)
		return
	}

	s.serveAsset(w, r, filepath, b)
}

// ServeIndexHTML serves index.html file
func (s *Static) ServeIndexHTML(initialState interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := path.Join(s.dir, indexFileName)
		b, err := assets.Asset(filepath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		options := s.options
		bData, err := json.Marshal(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		b = bytes.Replace(b, []byte(serverValuesPlaceholder), bData, 1)
		s.serveAsset(w, r, filepath, b)
	}
}

func (s *Static) serveAsset(w http.ResponseWriter, r *http.Request, filepath string, content []byte) {
	info, err := assets.AssetInfo(filepath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(content))
}
