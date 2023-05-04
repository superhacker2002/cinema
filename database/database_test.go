package database

import (
	"database/sql"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	type User struct {
		Id             int
		Username       string
		HashedPassword string
		Email          string
		CreditCardInfo string
		RoleId         int
	}

	db, err := sql.Open("postgres", "postgresql://<username>:<password>@<host>:<port>/<database>")
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.HashedPassword, &user.Email, &user.CreditCardInfo, &user.RoleId)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, but got %d", len(users))
	}
}
