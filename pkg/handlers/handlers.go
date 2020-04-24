package handlers

import (
	"QGTodo/pkg/db"
	"QGTodo/pkg/util/auth"
	"QGTodo/pkg/util/jwtauth"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Signup(queries *DB.Queries) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = queries.CreateUser(
			r.Context(),
			DB.CreateUserParams{
				Username: sql.NullString{String: creds.Username, Valid: true},
				Password: auth.HashPassword(creds.Password),
				CreatedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				UpdatedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			})
		if err != nil {
			print(err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)

		claims := &jwtauth.Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtauth.JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}

func Signin(queries *DB.Queries) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var creds Credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		user, err := queries.GetUserByName(r.Context(), sql.NullString{
			String: creds.Username,
			Valid:  true,
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return

		}

		if !auth.CheckPasswordHash(creds.Password, user.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expirationTime := time.Now().Add(5 * time.Minute)

		claims := &jwtauth.Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtauth.JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}

func Welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims, err := jwtauth.CheckClaims(w, r)
	if err != nil {
		log.Print(err.Error())
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
func GetTasksFromUser(queries *DB.Queries) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		claims, err := jwtauth.CheckClaims(w, r)
		if err != nil {
			log.Print(err.Error())
			return
		}
		ctx := r.Context()
		user, err := queries.GetUserByName(ctx, sql.NullString{
			String: claims.Username,
			Valid:  false,
		})
		if err != nil {
			log.Print(err.Error())
			return
		}
		tasks, err := queries.ParanoidListTasksFromUser(ctx, sql.NullInt32{
			Int32: user.ID,
			Valid: false,
		})
		if err != nil {
			log.Print(err.Error())
			return
		}
		jsonTasks, err := json.Marshal(tasks)
		if err != nil {
			log.Print(err.Error())
			return
		}
		w.Write(jsonTasks)
	}
}
func Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims, err := jwtauth.CheckClaims(w, r)

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtauth.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
