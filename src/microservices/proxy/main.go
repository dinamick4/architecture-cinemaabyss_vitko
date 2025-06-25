package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"os"
	"strings"
	"strconv"
	"math/rand"
	"time"
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
	http.HandleFunc("/api/users", handleUsers)
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

func handleUsers(w http.ResponseWriter, r *http.Request) {
    var targetURL string = os.Getenv("MONOLITH_URL") + "/api/users"
    proxyRequest(targetURL, w, r)
}


func handleMovies(w http.ResponseWriter, r *http.Request) {
    gradualMigrationStr := os.Getenv("GRADUAL_MIGRATION")
    migrationPercentInt, _ := strconv.Atoi(os.Getenv("MOVIES_MIGRATION_PERCENT"))

    shouldMigrate := strings.ToLower(gradualMigrationStr) == "true"
    var endpoint string = "/api/movies"
    var targetURL string
    // Случайная миграция при миграционном проценте равном 50%
    if shouldMigrate && migrationPercentInt == 50 {
        rand.Seed(int64(time.Now().UnixNano())) // Используем случайность на основе текущего времени

        // Берем случайное число от 0 до 100 и выбираем один из двух сервисов
        randomNum := rand.Intn(100)
        log.Printf("random %s", randomNum)
        if randomNum < 50 { // 50% вероятность отправить запрос в MOVIES_SERVICE_URL
            targetURL = os.Getenv("MOVIES_SERVICE_URL") + endpoint
        } else { // Остальные 50% отправляются в MONOLITH_URL
            targetURL = os.Getenv("MONOLITH_URL") + endpoint
        }
    } else if shouldMigrate && migrationPercentInt > 0 && migrationPercentInt <= 100 { // Полная миграция в случае процента больше нуля
        targetURL = os.Getenv("MOVIES_SERVICE_URL") + endpoint
    } else { // Иначе используем монолит
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