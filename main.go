package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type User struct {
	FullName string `json:"Name"`
	Age      int
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/test", sendData)
	r.Get("/create", createData)
	http.ListenAndServe(":3000", r)
}

func sendData(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		message := map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		}
		json.NewEncoder(w).Encode(message)
		return
	}
	makeArray := make(map[string]interface{})

	makeArray["id"] = id
	json.NewEncoder(w).Encode(makeArray)
}

func createData(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		message := map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		}

		json.NewEncoder(w).Encode(message)
		return
	}

	message := map[string]interface{}{
		"message": user.FullName,
		"status":  http.StatusCreated,
		"payload": user,
	}

	marshalData, _ := json.Marshal(message)
	fmt.Println(marshalData)

	json.NewEncoder(w).Encode(string(marshalData))
	return
}
