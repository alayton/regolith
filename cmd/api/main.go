package main

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
)

var (
	AutocertHostWhitelist = []string{""}
	AutocertEmail         = ""
	CORSOrigins           = []string{"https://*", "http://*"}
	CORSMethods           = []string{"GET", "POST", "PUT", "DELETE"}
	CORSHeaders           = []string{"Accept", "Authorization", "Content-Type", "Cache-Control"}
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	env := os.Getenv("REGOLITH_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	godotenv.Load(".env.local")
	godotenv.Load(".env." + env)
	godotenv.Load()

	// Create database connection, etc

	router := initRouter( /* pass your database connection or any other resources needed by your route handlers */ )

	useSSL := false
	sslString := os.Getenv("USE_SSL")
	if len(sslString) > 0 {
		if v, err := strconv.ParseBool(sslString); err == nil {
			useSSL = v
		}
	}

	if useSSL {
		if os.Getenv("SSL_MODE") == "static" {
			cache := os.Getenv("CERTS_DIR")
			err := http.ListenAndServeTLS(os.Getenv("LISTEN_ADDR"), path.Join(cache, "server.crt"), path.Join(cache, "server.key"), router)
			if err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Server exiting")
			}
		} else {
			certman := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(AutocertHostWhitelist...),
				Cache:      autocert.DirCache(os.Getenv("CACHE_DIR")),
				Email:      AutocertEmail,
			}

			httpServer := &http.Server{
				Addr:    os.Getenv("HTTP_LISTEN_ADDR"),
				Handler: certman.HTTPHandler(http.HandlerFunc(httpsRedirect)),
			}

			go httpServer.ListenAndServe()

			httpsServer := &http.Server{
				Addr:      os.Getenv("HTTPS_LISTEN_ADDR"),
				TLSConfig: certman.TLSConfig(),
				Handler:   router,
			}

			err := httpsServer.ListenAndServeTLS("", "")
			if err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Server exiting")
			}
		}
	} else {
		err := http.ListenAndServe(os.Getenv("LISTEN_ADDR"), router)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server exiting")
		}
	}
}

func httpsRedirect(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	url.Scheme = "https"
	url.Host = os.Getenv("CANONICAL_DOMAIN")
	http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
}
