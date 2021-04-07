# Introduction

If you’re a modern web developer, it is inevitable to ignore a database for long as it plays an important part in the application.

So in this post, I’ll be demonstrating how we can connect to a PostgreSQL database and perform basic SQL statements using Go.

# Prerequisites

You'll need [Go version 1.16+](https://golang.org/dl/) and [PostgreSQL](https://www.postgresql.org/download/) installed on your development machine.

> In order to connect with PostgreSQL we need driver, So we'll use [pgx](https://github.com/jackc/pgx) as our driver.

# Code time 🚀

## Code _v0.1_ 🌎

#### Aim

_Let's Start with a simple "Hello World.!" code._

Let’s create a new `main.go` file. Within this, we’ll import a few packages and set up a simple connection to an already running local database. for this tutorial, I'm using `postgres` as username, `123` as password, `localhost` network address, `:5432` default port, and `test` database.

You can change according to your setup.

```
postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
```

Now open `main.go` and write the following code.

```go
package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func main() {

    // Open up our database connection.
	conn, _ := pgx.Connect(context.Background(),
            "postgres://postgres:123@localhost:5432/test")

    // defer the close till after the main function has finished
    // executing
	defer conn.Close(context.Background())
	var greeting string
    //
	conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	fmt.Println(greeting)
}
```

`pgx.Connec()` establishes a connection with a PostgreSQL server with a connection string, this will return `pgx.Conn` is a PostgreSQL connection handle.

`conn.QueryRow()` executes sql query on the database, After that we store the response of data using `.Scan()`

## Code v1.0

So, now that we’ve successfully created a connection and build hello world with the database. Now let's start with a table and perform some queries over it.

#### Aim

Now we build a program where we can insert and fetch user data. We will understand various functions in the pgx package.

### Creating Table into PostgreSQL

```sql
CREATE TABLE IF NOT EXISTS USERS(
        ID          SERIAL   PRIMARY KEY,
        USERNAME    VARCHAR(20) NOT NULL UNIQUE
    );
```

Table name is `Users` with _ID_ and _USERNAME_ columns.

### Creating user struct

```go
// User is the model present in the database
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}
```

### Inserting User into database

So if we want to store a user into the database.

```go
//Creating temporary user object.
tmpUser := User{UserName: "Captain K"}
//Calling InsertUser Method
InsertUser(&tmpUser, conn)

func InsertUser(u *User, conn *pgx.Conn) {
	// Executing SQL query for insertion
	if _, err := conn.Exec(context.Background(), "INSERT INTO USERS(USERNAME) VALUES($1)", u.UserName); err != nil {
		// Handling error, if occur
		fmt.Println("Unable to insert due to: ", err)
		return
	}
	fmt.Println("Insertion Succesfull")
}
```

### Querying Multiple Rows

When we want to read all the users stored in the database.

```go
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
```

### Querying a Single Row

Find a user using user's ID.

```go
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
```

## Conclusion 🎉

In this post, we managed to set up a connection to a PostgreSQL and then perform some simple queries to that database and marshal the returned responses into a struct. This should hopefully give you everything you need in order to take things further and build your own Go applications on top of PostgreSQL.

> - [When to use db.Exec or db.Query ?](https://stackoverflow.com/a/50666083/8791826)
> - We can use pgx and pgxpool interchangeably but [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v4@v4.11.0/pgxpool) is a concurrency-safe connection pool for pgx. It is not safe for concurrent usage. Using a connection pool to manage access to multiple database connections from multiple goroutines.

### Recommended Reading: [REST server with Go in 5 minutes](https://dev.to/kushagra_mehta/rest-server-with-go-in-5-minutes-3n8l)

_[post banner](https://medium.com/@amoghagarwal/insert-optimisations-in-golang-26884b183b35)_
