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
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

// migratedPrefixes lists route prefixes that have already been rewritten
// in Go and should be routed to the new service. Everything else falls
// through to the legacy Nest service by default.
var migratedPrefixes = []string{
	"/api/v1/policyholders",
	"/api/v1/policies",
	"/api/v1/notes",
}

func main() {
	legacyURL := mustParseURL(envOr("LEGACY_NEST_URL", "http://localhost:3001"))
	goServiceURL := mustParseURL(envOr("GO_SERVICE_URL", "http://localhost:8080"))
	port := envOr("GATEWAY_PORT", "8000")

	legacyProxy := httputil.NewSingleHostReverseProxy(legacyURL)
	goProxy := httputil.NewSingleHostReverseProxy(goServiceURL)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to create static sub fs: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	mux := http.NewServeMux()

	// Serve the visual Migration Control Center Dashboard UI on the root route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If it's a documentation request, serve doc content directly
		if strings.HasPrefix(r.URL.Path, "/api/docs/") {
			docName := strings.TrimPrefix(r.URL.Path, "/api/docs/")
			var filePath string
			switch docName {
			case "adr-0001":
				filePath = "../docs/adr/0001-go-over-node.md"
			case "adr-0002":
				filePath = "../docs/adr/0002-pgx-sqlc-over-orm.md"
			case "adr-0003":
				filePath = "../docs/adr/0003-strangler-fig-over-big-bang.md"
			case "adr-0004":
				filePath = "../docs/adr/0004-shared-db-coexistence-over-cdc.md"
			case "adr-0005":
				filePath = "../docs/adr/0005-dark-launching-and-canary-telemetry.md"
			case "adr-0006":
				filePath = "../docs/adr/0006-react-typescript-dashboard-ui.md"
			case "claude":
				filePath = "../CLAUDE.md"
			case "ai-workflow":
				filePath = "../docs/AI_WORKFLOW.md"
			case "skill":
				filePath = "../.agents/skills/ai-pr-review/SKILL.md"
			default:
				http.Error(w, "doc not found", http.StatusNotFound)
				return
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				http.Error(w, "failed to read doc: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(content)
			return
		}

		// If it's an API request, route it through the Strangler Gateway
		if strings.HasPrefix(r.URL.Path, "/api/") {
			target := legacyProxy
			targetURL := legacyURL.String()
			routedTo := "legacy-nest (:3001)"

			// Feature Flag / Manual Override check (X-Force-Backend header or ?force query param)
			force := r.Header.Get("X-Force-Backend")
			if force == "" {
				force = r.URL.Query().Get("force")
			}

			if force == "go" {
				target = goProxy
				targetURL = goServiceURL.String()
				routedTo = "go-service (:8080) [forced]"
			} else if force == "legacy" {
				target = legacyProxy
				targetURL = legacyURL.String()
				routedTo = "legacy-nest (:3001) [forced]"
			} else if isMigrated(r.URL.Path) {
				target = goProxy
				targetURL = goServiceURL.String()
				routedTo = "go-service (:8080)"
			}

			w.Header().Set("X-Routed-To", routedTo)
			w.Header().Set("X-Target-URL", targetURL)
			log.Printf("%s %s (force=%s) -> %s (%s)", r.Method, r.URL.Path, force, routedTo, targetURL)
			target.ServeHTTP(w, r)
			return
		}

		// Serve the compiled React Dashboard UI static assets
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("gateway listening on :%s (legacy=%s, go=%s)", port, legacyURL, goServiceURL)
	log.Printf("open visual dashboard at http://localhost:%s", port)
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
