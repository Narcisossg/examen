package app

import (
	"encoding/json"
	"examen/db"
	"examen/models"
	"examen/utils"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	existing, err := db.GetUserByUsername(creds.Username)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error requesting user by username", http.StatusInternalServerError)
		return
	}

	if existing == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token":"%s"}`, token)
}
