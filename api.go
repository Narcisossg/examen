package main

import (
	"2025A3/app"
	"2025A3/db"
	"fmt"
	"net/http"
)

func main() {
	db.Conn = db.NewDB()

	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /users/{$}", app.GetAllUsers)
	http.HandleFunc("POST /users/{$}", app.CreateUser)
	http.HandleFunc("DELETE /users/{userId}", app.DeleteUser)
	// http.HandleFunc("GET /users/{user}", app.GetUserById)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8080", nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Conn.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Ping to DB successfull from %s", r.UserAgent())
}

// patch/ put / le redex / melenger déux "structures" differentes / POSTMAN !!!
