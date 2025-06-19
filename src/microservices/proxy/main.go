package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"os"
	"strings"
)

// Models
type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
}

func main() {

	// Set up HTTP routes
	http.HandleFunc("/api/movies", handleMovies)
	http.HandleFunc("/health", handleHealth)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Starting movies proxy on port %s", port)
	log.Fatal(http.ListenAndServe(":"+ port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}


func handleMovies(w http.ResponseWriter, r *http.Request) {
    migrationStr := os.Getenv("GRADUAL_MIGRATION")
    shouldMigrate := strings.ToLower(migrationStr) == "true"

    var endpoint string = "/api/movies"
    var targetURL string
    if shouldMigrate {
        targetURL = os.Getenv("MOVIES_SERVICE_URL") + endpoint
    } else {
        targetURL = os.Getenv("MONOLITH_URL") + endpoint
    }
    log.Printf("Invoke %s", targetURL)
    proxyRequest(targetURL, w, r)
}

// proxyRequest выполняет переадресацию запроса на другой микросервис
func proxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
    // Копируем запрос
    req, err := http.NewRequest(r.Method, targetURL+r.URL.RawQuery, r.Body)
    if err != nil {
        log.Fatalf("Ошибка копирования запроса: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Копируем заголовки
    for k, vv := range r.Header {
        for _, v := range vv {
            req.Header.Add(k, v)
        }
    }

    // Отправляем запрос в другой микросервис
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatalf("Ошибка при выполнении запроса: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Читаем тело ответа
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Ошибка при чтении тела ответа: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Возвращаем клиенту полученные данные
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}