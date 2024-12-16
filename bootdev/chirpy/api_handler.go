package main

import (
	"fmt"
	"net/http"
)

func readinessHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err := rw.Write([]byte("OK"))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/html")
	responseContent := fmt.Sprintf(`
    <html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
    </html>
    `, cfg.fileServerHits.Load())
	_, err := rw.Write([]byte(responseContent))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) resetAppHandler(rw http.ResponseWriter, r *http.Request) {
	err := cfg.resetUsers(r)
	if err != nil {
		if err.Error() == "Cannot reset users in prod !" {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg.resetMetrics()
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.Write([]byte("Hits and users reset to zero!"))
}

func (cfg *apiConfig) resetMetrics() {
	cfg.fileServerHits.Store(0)
}
