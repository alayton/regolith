package main

import (
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/alayton/regolith/src/api"
)

func initRouter( /* pass your database handle or any other resources needed by your route handlers */ ) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   CORSOrigins,
			AllowedMethods:   CORSMethods,
			AllowedHeaders:   CORSHeaders,
			AllowCredentials: true,
			MaxAge:           900,
		}),
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/", func(r chi.Router) {
		r.Get("/", InjectData([]Injectable{{[]interface{}{"dummy"}, api.GetDummyData}}))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/dummy", api.GetDummyData)
	})

	staticRoot := os.Getenv("PUBLIC_ROOT")
	fs := http.FileServer(http.Dir(staticRoot))

	if os.Getenv("ENABLE_GZIP") == "true" {
		fs = gziphandler.GzipHandler(fs)
	}

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{w, 0}
		fs.ServeHTTP(writer, r)
		if writer.Status == http.StatusNotFound {
			// These headers are written by FileServer when a file isn't found
			writer.Header().Del("Content-Type")
			writer.Header().Del("X-Content-Type-Options")

			http.ServeFile(w, r, staticRoot+"/index.html")
		}
	})

	return r
}

type responseWriter struct {
	writer http.ResponseWriter
	Status int
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.Status == http.StatusNotFound {
		// Discard writes if the file wasn't found so we can serve index.html instead
		return 0, nil
	}
	return w.writer.Write(b)
}

func (w *responseWriter) WriteHeader(status int) {
	w.Status = status
	if status != http.StatusNotFound {
		w.writer.WriteHeader(status)
	}
}
