package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/dogeorg/reflector/pkg/database"
	"github.com/go-chi/chi/v5"
)

type Entry struct {
	Token string `json:"token"`
	IP    string `json:"ip"`
}

func CreateEntry(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if !isValidToken(entry.Token) || !isValidIP(entry.IP) {
			http.Error(w, "Invalid token or IP format", http.StatusBadRequest)
			return
		}

		if err := db.SaveEntry(entry.Token, entry.IP); err != nil {
			http.Error(w, "Failed to save entry", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetIP(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := chi.URLParam(r, "token")

		if token == "me" {
			ip := r.RemoteAddr
			if fwdIP := r.Header.Get("X-Forwarded-For"); fwdIP != "" {
				ip = fwdIP
			}
			// Strip the port from the IP address
			if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
				ip = ip[:colonIndex]
			}
			json.NewEncoder(w).Encode(map[string]string{"ip": ip})
			return
		}

		ip, err := db.GetIP(token)
		if err != nil {
			http.Error(w, "IP not found", http.StatusNotFound)
			return
		}

		if err := db.DeleteEntry(token); err != nil {
			http.Error(w, "Failed to remove entry", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"ip": ip})
	}
}

func isValidToken(token string) bool {
	return len(token) <= 20
}

func isValidIP(ip string) bool {
	ipPattern := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(ipPattern, ip)
	return match && len(ip) <= 15
}
