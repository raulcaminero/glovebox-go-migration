// Package main implements a strangler-fig gateway: a thin reverse proxy
// that sits in front of both the legacy NestJS service and the new Go
// service, and routes each request to one or the other based on a feature
// flag. This is what lets a migration ship incrementally — module by
// module, or even route by route — instead of a risky big-bang cutover.
//
// See docs/adr/0003-strangler-fig-over-big-bang.md for why this pattern
// was chosen, and docs/MIGRATION.md for the rollout phases that use it.
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// migratedPrefixes lists route prefixes that have already been rewritten
// in Go and should be routed to the new service. Everything else falls
// through to the legacy Nest service by default. In a real rollout this
// would be a config value (or a proper feature-flag service) updated as
// each module's migration completes — not a code change per cutover.
var migratedPrefixes = []string{
	"/api/v1/policyholders",
	"/api/v1/policies",
	"/api/v1/notes",
}

func main() {
	legacyURL := mustParseURL(envOr("LEGACY_NEST_URL", "http://localhost:3001"))
	goServiceURL := mustParseURL(envOr("GO_SERVICE_URL", "http://localhost:8080"))
	port := envOr("PORT", "8000")

	legacyProxy := httputil.NewSingleHostReverseProxy(legacyURL)
	goProxy := httputil.NewSingleHostReverseProxy(goServiceURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := legacyProxy
		routedTo := "legacy-nest"
		if isMigrated(r.URL.Path) {
			target = goProxy
			routedTo = "go-service"
		}
		w.Header().Set("X-Routed-To", routedTo) // makes the routing decision visible in the demo
		log.Printf("%s %s -> %s", r.Method, r.URL.Path, routedTo)
		target.ServeHTTP(w, r)
	})

	log.Printf("gateway listening on :%s (legacy=%s, go=%s)", port, legacyURL, goServiceURL)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func isMigrated(path string) bool {
	for _, prefix := range migratedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid URL %q: %v", raw, err)
	}
	return u
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
