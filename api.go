package main

import (
	"examen/app"
	"examen/db"
	"fmt"
	"net/http"
)

func main() {
	db.Conn = db.NewDB()

	http.HandleFunc("GET /{$}", healthCheck)

	http.HandleFunc("POST /login/{$}", app.Login)

	http.HandleFunc("GET /users/{$}", app.GetAllUsers)
	http.HandleFunc("GET /users/{userId}", app.GetUserById)
	http.HandleFunc("POST /users/{$}", app.CreateUser)
	http.HandleFunc("PUT /users/{userId}", app.UpdateUser)
	http.HandleFunc("DELETE /users/{userId}", app.DeleteUser)

	fmt.Println("Listening at http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Conn.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Ping to DB successfull from %s", r.UserAgent())
}
