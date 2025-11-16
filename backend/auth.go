
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte("replace-this-with-secure-secret")

type User struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	UserID string `json:"userId"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err!=nil { http.Error(w, err.Error(), 400); return }
	if u.Email=="" || u.Password=="" || u.Name=="" { http.Error(w, "missing fields", 400); return }
	hash, err := hashPassword(u.Password)
	if err!=nil { http.Error(w, "password error", 500); return }
	id := uuid.New().String()
	_, err = db.ExecContext(r.Context(), `INSERT INTO users(id,name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,now())`, id, u.Name, u.Email, hash)
	if err!=nil { http.Error(w, err.Error(), 500); return }
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&creds); err!=nil { http.Error(w, err.Error(), 400); return }
	var id, hash string
	err := db.QueryRowContext(r.Context(), `SELECT id,password_hash FROM users WHERE email=$1`, creds.Email).Scan(&id,&hash)
	if err==sql.ErrNoRows { http.Error(w, "invalid credentials", 401); return }
	if err!=nil { http.Error(w, err.Error(), 500); return }
	if err := checkPasswordHash(hash, creds.Password); err!=nil { http.Error(w, "invalid credentials", 401); return }
	// Create token
	exp := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		Email: creds.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(jwtKey)
	if err!=nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(map[string]string{"token": s})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth=="" { http.Error(w, "missing auth", 401); return }
		var tokenString string
		if len(auth) > 7 && auth[:7] == "Bearer " { tokenString = auth[7:] } else { http.Error(w, "invalid auth header", 401); return }
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { return jwtKey, nil })
		if err!=nil || !token.Valid { http.Error(w, "invalid token", 401); return }
		// Store user id in context
		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
