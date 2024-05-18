package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type indexCache struct {
	Mutex     sync.RWMutex
	Updated   time.Time
	Modified  time.Time
	PreSplit  []byte
	PostSplit []byte
}

type Injectable struct {
	Key     []interface{}
	Handler http.HandlerFunc
}

type injectableData struct {
	Key  []interface{}
	Data []byte
}

var (
	index       = &indexCache{}
	injectBegin = []byte(`<script type="text/javascript">window.regolithInjectables = [];`)
	injectEnd   = []byte(`</script>`)
)

const (
	IndexCacheTime  = 10 * time.Second
	InjectionMarker = "<!--@inject-->"
)

func InjectData(injectables []Injectable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []injectableData
		interceptor := &responseInterceptor{}
		for _, injectable := range injectables {
			interceptor.Reset()
			injectable.Handler(interceptor, r)
			data = append(data, injectableData{injectable.Key, slices.Clone(interceptor.Bytes())})
		}

		index.Mutex.RLock()
		now := time.Now()
		since := now.Sub(index.Updated)
		index.Mutex.RUnlock()

		if since > IndexCacheTime {
			index.Mutex.Lock()

			// Even though we could still error out, we still set the updated time to generate less error spam
			index.Updated = now

			indexFilename := filepath.Join(os.Getenv("PUBLIC_ROOT"), "index.html")
			stat, err := os.Stat(indexFilename)
			if err != nil {
				index.Mutex.Unlock()
				log.Error().Err(err).Msg("Failed to stat index.html")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			modded := stat.ModTime()
			if modded != index.Modified {
				index.Modified = modded
				indexBytes, err := os.ReadFile(indexFilename)
				if err != nil {
					index.Mutex.Unlock()
					log.Error().Err(err).Msg("Failed to read index.html")
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				indexString := string(indexBytes)
				injectPos := strings.Index(indexString, InjectionMarker)
				if injectPos == -1 {
					index.Mutex.Unlock()
					log.Error().Msg("Failed to find injection marker in index.html")
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				injectPos += len(InjectionMarker)

				index.PreSplit = []byte(indexString[:injectPos])
				index.PostSplit = []byte(indexString[injectPos:])
			}

			index.Mutex.Unlock()
		}

		index.Mutex.RLock()

		w.Write(index.PreSplit)
		w.Write(injectBegin)
		for _, value := range data {
			keyBytes, _ := json.Marshal(value.Key)
			w.Write([]byte(`window.regolithInjectables.push([`))
			w.Write(keyBytes)
			w.Write([]byte(","))
			w.Write(value.Data)
			w.Write([]byte(`]);\n`))
		}
		w.Write(injectEnd)
		w.Write(index.PostSplit)

		index.Mutex.RUnlock()
	}
}

type responseInterceptor struct {
	buf bytes.Buffer
}

func (r *responseInterceptor) Header() http.Header {
	return http.Header{}
}

func (r *responseInterceptor) Write(data []byte) (int, error) {
	return r.buf.Write(data)
}

func (r *responseInterceptor) WriteHeader(statusCode int) {}

func (r *responseInterceptor) Reset() {
	r.buf.Reset()
}

func (r *responseInterceptor) Bytes() []byte {
	return r.buf.Bytes()
}

var _ http.ResponseWriter = &responseInterceptor{}
