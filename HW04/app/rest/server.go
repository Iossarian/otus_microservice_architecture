package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var db *sql.DB

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ListenAndServer(_ context.Context, dsn string) {
	var err error
	db, err = initDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/user/{userId}", getUser).Methods("GET")
	r.HandleFunc("/user/{userId}", updateUser).Methods("PUT")
	r.HandleFunc("/user/{userId}", deleteUser).Methods("DELETE")

	fmt.Println("Server is running on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("error 1", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users(name, age) VALUES($1, $2)", user.Name, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)

		return
	}

	var user User

	err = db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = db.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", user.Name, user.Age, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	user.ID = userId
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)

		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id=$1", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
