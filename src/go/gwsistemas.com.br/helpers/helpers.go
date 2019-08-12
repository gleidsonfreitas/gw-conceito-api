package helpers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Connection returns a database connection whose
// url connection is loaded from environment variables.
// The caller function is responsible for closing the connection
// invoking defer db.Close().
func Connection() (*sql.DB, error) {
	defer TimeTrack("Connection", time.Now())

	params := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// TimeTrack
func TimeTrack(name string, start time.Time) {
	elapsed := time.Now().Sub(start).Seconds() * 1e3
	fmt.Println(" > " + name + ": " + fmt.Sprintf("%0.2f", elapsed) + " ms")
}

// MD5
func MD5(plain string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(plain)))
}
