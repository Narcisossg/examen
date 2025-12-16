package app

import (
	"2025A3/db"
	"2025A3/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users, err = db.GetAllUsers()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(
			w,
			"erreur de récupération des users",
			http.StatusInternalServerError)
		return
	}
	encodedUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(
			w,
			"erreur de conversion des users",
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(encodedUsers))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var userDto models.User
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// pour faire une regex
	//match, _ := regexp.Match("e([a-z]+)gi", []byte(userDto.Username))

	var errMsgs []string
	if len(userDto.Username) < 4 || len(userDto.Username) > 50 {
		errMsgs = append(errMsgs,
			"Username must have a length between 4 and 50")
	}
	if strings.Contains(userDto.Username, "Langage C") {
		errMsgs = append(errMsgs,
			"Username must not contains the forbidden word")
	}
	if len(userDto.Password) < 4 || len(userDto.Password) > 50 {
		errMsgs = append(errMsgs,
			"Password must have a length between 4 and 50")
	}
	if strings.ContainsAny(userDto.Password, "-!?") {
		errMsgs = append(errMsgs,
			"Password must have at least 1 special char [-!?]")
	}
	if userDto.Credit < 0 {
		errMsgs = append(errMsgs,
			"Credit must not be negative")
	}

	if len(errMsgs) > 0 {
		jsonErrs, _ := json.Marshal(errMsgs)
		http.Error(w, string(jsonErrs), http.StatusBadRequest)
		return
	}

	duplicates, err := db.GetAllUsersByName(userDto.Username)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}
	if len(duplicates) > 0 {
		http.Error(w, "Username must be unique", http.StatusConflict)
		return
	}

	err = db.CreateUser(userDto)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}
