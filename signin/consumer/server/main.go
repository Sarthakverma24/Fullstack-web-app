package main

import (
	"awesomeProject/db"
	"awesomeProject/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	loggers := db.Logger.Sugar()

	scyllaWrapper, _ := db.NewScyllaDBSession()
	if scyllaWrapper == nil {
		loggers.Fatal("Failed to initialize scyllaWrapper")
	}
	loggers.Infoln("Successfully initialized ScyllaDB connection")

	http.HandleFunc("/signin", handleSignIn(scyllaWrapper))
	http.HandleFunc("/login", handleLogin(scyllaWrapper))
	log.Println("HTTP server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleSignIn(scyllaWrapper db.Database) http.HandlerFunc {
	http.HandleFunc("/recover", handleForgetPassword(scyllaWrapper))
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user utils.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Bad request: invalid JSON or data types", http.StatusBadRequest)
			return
		}
		if user.UserID == "" || user.UserName == "" || user.PhoneNo == "" || user.Gmail == "" || user.Password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if user.Age <= 0 || user.Age > 120 {
			http.Error(w, "Invalid age value", http.StatusBadRequest)
			return
		}
		queryCheck := `SELECT user_id FROM User.User_data WHERE gmail = ? ALLOW FILTERING`
		resultRaw, err := scyllaWrapper.GetData(context.Background(), queryCheck, user.Gmail)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if resultRaw != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("User %s already exist with same User ID", user.UserName),
			})
			return
		}

		query := `INSERT INTO User.User_data (user_id, user_name, phone_no, age, password, gmail) VALUES (?, ?, ?, ?, ?, ?)`
		err = scyllaWrapper.SetData(context.Background(), query,
			user.UserID, user.UserName, user.PhoneNo, user.Age, user.Password, user.Gmail)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("User %s created successfully!", user.UserName),
		})
	}
}

func handleLogin(scyllaWrapper db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var userLogin utils.UserLogin
		if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
			http.Error(w, "Bad request: invalid JSON or data types", http.StatusBadRequest)
			return
		}

		if userLogin.UserName == "" || userLogin.Password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		query := `SELECT user_id FROM User.User_data WHERE user_name = ? AND password = ? ALLOW FILTERING`
		resultRaw, err := scyllaWrapper.GetData(context.Background(), query, userLogin.UserName, userLogin.Password)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		results, ok := resultRaw.([]map[string]interface{})
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(results) == 0 {
			http.Error(w, "User not found or incorrect password", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("User %s is logged in", userLogin.UserName),
		})
	}
}
func handleForgetPassword(scyllaWrapper db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req utils.ForgetPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request: invalid JSON or data types", http.StatusBadRequest)
			return
		}

		if req.Gmail == "" || req.Phone == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		query := `SELECT user_name, password FROM User.User_data WHERE gmail = ? AND phone_no = ? ALLOW FILTERING`
		resultRaw, err := scyllaWrapper.GetData(context.Background(), query, req.Gmail, req.Phone)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		results, ok := resultRaw.([]map[string]interface{})
		if !ok || len(results) == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user := results[0]
		password := user["password"]

		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("- %v", password),
		})
	}
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	}
}
