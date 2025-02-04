package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ListenAndServe(_ context.Context, dsn string) {
	var err error

	db, err = initDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/auth/register", register).Methods("POST")
	r.HandleFunc("/auth/login", login).Methods("POST")
	r.HandleFunc("/auth/logout", logout).Methods("GET")
	r.HandleFunc("/auth/users/{userId}", auth).Methods("GET")

	fmt.Println("Server is running on port 8001")

	log.Fatal(http.ListenAndServe(":8001", r))
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	_, err = db.Exec("INSERT INTO profile(username, password) VALUES($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var userID int

	err = db.QueryRow("SELECT id FROM profile WHERE username = $1", user.Username).Scan(&userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	user.ID = userID

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-User-Id", fmt.Sprintf("%v", userID))
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var userID int
	var storedHashedPassword string

	err = db.QueryRow("SELECT id,password FROM profile WHERE username = $1", user.Username).
		Scan(&userID, &storedHashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		}

		http.Error(w, "Failed to get user", http.StatusInternalServerError)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "User not found", http.StatusForbidden)

		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	sessions.Lock()
	sessions.m[sessionID] = userID
	sessions.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Failed to get session cookie", http.StatusBadRequest)

		return
	}

	sessionID := cookie.Value
	sessions.Lock()
	delete(sessions.m, sessionID)
	sessions.Unlock()

	// Сразу "протухаем" куку
	cookie = &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func auth(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized: no cookie", http.StatusUnauthorized)

		return
	}

	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)

		return
	}

	sessionID := sessionCookie.Value

	sessions.RLock()
	storedUserID, ok := sessions.m[sessionID]
	sessions.RUnlock()

	if !ok || userID != storedUserID {
		http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)

		return
	}

	w.WriteHeader(http.StatusOK)
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
