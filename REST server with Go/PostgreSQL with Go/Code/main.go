package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// User is the model present in the database
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

func main() {

	conn, _ := pgx.Connect(context.Background(),
		"postgres://postgres:123@localhost:5432/test")

	// defer the close till after the main function has finished
	// executing
	defer conn.Close(context.Background())

	//Creating temporary user object.
	tmpUser := User{UserName: "Captain K"}
	//Calling InsertUser Method
	InsertUser(&tmpUser, conn)

}

func InsertUser(u *User, conn *pgx.Conn) {

	// Executing SQL query for insertion
	if _, err := conn.Exec(context.Background(), "INSERT INTO USERS(USERNAME) VALUES($1)", u.UserName); err != nil {
		// Handling error, if occur
		fmt.Println("Unable to insert due to: ", err)
		return
	}
	fmt.Println("Insertion Succesfull")
}

func GetAllUsers(conn *pgx.Conn) {

	// Execute the query
	if rows, err := conn.Query(context.Background(), "SELECT * FROM USERS"); err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return
	} else {

		// carefully deferring Queries closing
		defer rows.Close()

		// Using tmp variable for reading
		var tmp User

		// Next prepares the next row for reading.
		for rows.Next() {
			// Scan reads the values from the current row into tmp
			rows.Scan(&tmp.ID, &tmp.UserName)
			fmt.Printf("%+v\n", tmp)
		}
		if rows.Err() != nil {
			// if any error occurred while reading rows.
			fmt.Println("Error will reading user table: ", err)
			return
		}
	}
}

func GetAnUser(id int, conn *pgx.Conn) {

	// variable to store username
	var username string

	// Executing query for single row
	if err := conn.QueryRow(context.Background(), "SELECT USERNAME WHERE ID=$1", id).Scan(&username); err != nil {
		fmt.Println("Error occur while finding user: ", err)
		return
	}
	fmt.Printf("User with id=%v is %v\n", id, username)
}
