package db

import (
	"2025A3/models"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	var rows *sql.Rows

	rows, err := Conn.Query("SELECT id, username, password, credit FROM esgi.users")
	if err != nil {
		return nil, fmt.Errorf("db getAllUsers : %v", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Credit)
		if err != nil {
			return nil, fmt.Errorf("db getAllUsers : %v", err.Error())
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("db getAllUsers : %v", err.Error())
	}

	return users, nil
}

func GetAllUsersByName(name string) ([]models.User, error) {
	var users []models.User

	var rows *sql.Rows

	rows, err := Conn.Query("SELECT id, username, password, credit FROM esgi.users WHERE username = $1", name)
	if err != nil {
		return nil, fmt.Errorf("db getAllUsersByName : %v", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Credit)
		if err != nil {
			return nil, fmt.Errorf("db getAllUsersByName : %v", err.Error())
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("db getAllUsersByName : %v", err.Error())
	}

	return users, nil
}

func CreateUser(user models.User) error {
	_, err := Conn.Exec("INSERT INTO esgi.users (username, password, credit) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Credit)
	if err != nil {
		return fmt.Errorf("db createUser : %v", err.Error())
	}
	return nil
}
