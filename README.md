# Regolith

Regolith serves as a starter template for a React SPA and Go webserver + API monorepo.

Additionally, it allows for preloading page data to speed up initial page rendering. This is accomplished by (optionally) defining a route for each page with a list of endpoints that should be injected into the response.

Example from `cmd/api/routes.go`:
```
	r.Route("/", func(r chi.Router) {
		r.Get("/", InjectData([]Injectable{{[]interface{}{"dummy"}, api.GetDummyData}}))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/dummy", api.GetDummyData)
	})
```
The first route handles requests to the domain root and defines one endpoint to inject the result of into the response with the React Query key `dummy`. When the client starts running the React app, it takes each injected response and, using React Query, sets the query data for that key. When the app tries to request that endpoint, it finds the cached data and doesn't have to make a separate request. This only applies when loading data with React Query with the same key on both server and client.

If no route exists for a URL the app handles, the page renders without any data being injected.

## Directory Structure

```bash
├── app - The React app
│   ├── assets - Assets to include as part of the build
│   ├── pages - Contains components grouped by pages they belong on
│   └── shared - Contains components and code not belonging to a specific page
├── cache - Static SSL certs go here, named server.crt and server.key
├── cmd
│   └── api - Go API main module
├── dist - Building the React app places the built files and the contents of /public here
├── node_modules
├── public - Files that should always be publically accessible
└── src - Your Go API modules
```

Placing main modules in `cmd` is a typical Go convention. Other modules could go anywhere you like, including directly within the root.
Using `src` by default is purely a matter of taste, and can be changed by updating any imports referencing those files.

## Building & Running

Go:
```bash
go build ./cmd/api
./api
```

React:
```base
yarn build
[or]
yarn dev
```

## Environment Variables

In order of priority:
1. `.env.[environment].local`
2. `.env.local`
3. `.env.[environment]`
4. `.env`

Currently used variables:
- SSL_MODE - `static` or anything else. With `static`, SSL certs are read from `CERTS_DIR/server.(crt|key)`. Otherwise, LetsEncrypt is used, with the whitelisted hosts and contact email set in `cmd/api/main.go`
- USE_SSL - `true` or `false`. If `false`, only `HTTP_LISTEN_ADDR` is used.
- HTTP_LISTEN_ADDR - The local address to listen for HTTP traffic on, e.g. `:80`. When using LetsEncrypt, this automatically redirects to `https://CANONICAL_DOMAIN`
- HTTPS_LISTEN_ADDR - The local address to listen for HTTPS traffic on, e.g. `:443`
- CANONICAL_DOMAIN - The "official" domain the app expects to be served on, e.g. `example.com`
- PUBLIC_ROOT - The directory the Go application will serve static files from. Should usually be `./dist`
- CERTS_DIR - The directory the SSL certificate is read from when `SSL_MODE=static`
- ENABLE_GZIP - `true` or `false`. Whether or not to gzip static content served by the Go application
- VITE_API_URL - The URL that the React app should use for API requests, e.g. `/api` or `http://localhost/api` for development

For a variable to be accessible within the React app, it must be prefixed with `VITE_`.
