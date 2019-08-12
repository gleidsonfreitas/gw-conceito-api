package repository

import (
	"database/sql"
	"fmt"
	"time"

	"gwsistemas.com.br/helpers"
	"gwsistemas.com.br/models"

	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

/* User */

// CreateUser
func CreateUser(db *sql.DB, name string, email string, password string) (*models.User, error) {
	defer helpers.TimeTrack("CreateUser", time.Now())

	sql :=
		`INSERT INTO users
			(name, email, passwordDigest)
		VALUES
			($1, $2, $3)
		RETURNING id`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Hash the password
	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(password), 31)

	var userId int
	err = stmt.QueryRow(name, email, string(passwordDigest)).Scan(&userId)
	if err != nil {
		return nil, err
	}

	return FindUser(db, userId)
}

// FindUser
func FindUser(db *sql.DB, userId int) (*models.User, error) {
	defer helpers.TimeTrack("FindUser", time.Now())

	user := new(models.User)
	row := db.QueryRow(`SELECT id, name, email, password_digest, enabled FROM users WHERE id = $1`, userId)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.Enabled)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

// FindUserByEmail
func FindUserByEmail(db *sql.DB, email string) (*models.User, error) {
	defer helpers.TimeTrack("FindUserByEmail", time.Now())

	user := new(models.User)
	row := db.QueryRow(`SELECT id, name, email, password_digest, enabled FROM users WHERE email = $1`, email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.Enabled)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

// FindUsers
func FindUsers(db *sql.DB) ([]*models.User, error) {
	defer helpers.TimeTrack("FindUsers", time.Now())

	sql := `SELECT id, name, email, password_digest, enabled, created_at, updated_at FROM users ORDER BY name`

	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordDigest, &user.Enabled, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}
