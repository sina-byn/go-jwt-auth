package user

import (
	"database/sql"
	"errors"

	"github.com/sina-byn/go-jwt-auth/pkg/db"
	"github.com/sina-byn/go-jwt-auth/pkg/utils"
)

func CreateUser(email, password string) (*int64, error) {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	result, err := db.DB.Exec(query, email, hashedPassword)

	if err != nil {
		return nil, err
	}

	userId, err := result.LastInsertId()

	return &userId, err
}

func GetUserByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, email)
	var user User

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Fullname)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func UpdateUser(id int64, user *User) error {
	stmt, err := db.DB.Prepare("UPDATE users SET email=?, password=?, fullname=? WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(user.Email, user.Password, user.Fullname, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int64) error {
	stmt, err := db.DB.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)

	return err
}
