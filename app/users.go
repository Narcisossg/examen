package app

import (
	"2025A3/db"
	"2025A3/models"
	"2025A3/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users, err = db.GetAllUsers()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de récupération des users", http.StatusInternalServerError)
		return
	}

	encodedUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "erreur de conversion des users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(encodedUsers))
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserById(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de récupération du user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	encodedUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "erreur de conversion du user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(encodedUser))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto models.User
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var errMsgs []string
	if len(userDto.Username) < 4 || len(userDto.Username) > 50 {
		errMsgs = append(errMsgs, "Username must have a length between 4 and 50")
	}
	if strings.Contains(userDto.Username, "Langage C") {
		errMsgs = append(errMsgs, "Username must not contains the forbidden word")
	}
	if len(userDto.Password) < 4 || len(userDto.Password) > 50 {
		errMsgs = append(errMsgs, "Password must have a length between 4 and 50")
	}
	if !strings.ContainsAny(userDto.Password, "-!?") {
		errMsgs = append(errMsgs, "Password must have at least 1 special char [-!?]")
	}
	if userDto.Credit < 0 {
		errMsgs = append(errMsgs, "Credit must not be negative")
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de hashage", http.StatusInternalServerError)
		return
	}
	userDto.Password = string(hashed)

	err = db.CreateUser(userDto)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	existing, err := db.GetUserById(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var userDto models.User
	err = json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var errMsgs []string
	if len(userDto.Username) < 4 || len(userDto.Username) > 50 {
		errMsgs = append(errMsgs, "Username must have a length between 4 and 50")
	}
	if strings.Contains(userDto.Username, "Langage C") {
		errMsgs = append(errMsgs, "Username must not contains the forbidden word")
	}
	if len(userDto.Password) < 4 || len(userDto.Password) > 50 {
		errMsgs = append(errMsgs, "Password must have a length between 4 and 50")
	}
	if !strings.ContainsAny(userDto.Password, "-!?") {
		errMsgs = append(errMsgs, "Password must have at least 1 special char [-!?]")
	}
	if userDto.Credit < 0 {
		errMsgs = append(errMsgs, "Credit must not be negative")
	}

	if len(errMsgs) > 0 {
		jsonErrs, _ := json.Marshal(errMsgs)
		http.Error(w, string(jsonErrs), http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de hashage", http.StatusInternalServerError)
		return
	}
	userDto.Password = string(hashed)

	err = db.UpdateUser(id, userDto)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	username, err := utils.VerifyJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userToDelete, err := db.GetUserById(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	if userToDelete == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	connectedUser, err := db.GetUserByUsername(username)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	if connectedUser.Id != userToDelete.Id {
		http.Error(w, "Forbidden: You can only delete your own account", http.StatusForbidden)
		return
	}

	err = db.DeleteUser(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur inconnue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
