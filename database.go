package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "upper.io/db.v3"
	"upper.io/db.v3/mysql"
)

type User struct {
	ID int `db:"id"`
	// Maps the "Name" property to the "name" column
	// of the "birthday" table.
	Username string `db:"username"`
	// Maps the "Born" property to the "born" column
	// of the "birthday" table.
	Password string `db:"password"`
}

func upperDemo(c *gin.Context) {

	c.String(http.StatusOK, "DB Openning...\n")
	sess, err := mysql.Open(DBConfig)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close()

	c.String(http.StatusOK, "Table Using...\n")
	userCollection := sess.Collection("users")

	// go func() {
	// 	c.String(http.StatusOK, "Table Truncate...")
	// }()
	// err = userCollection.Truncate()
	// if err != nil {
	// 	log.Fatalf("Truncate(): %q\n", err)
	// }

	c.String(http.StatusOK, "Inseting User...\n")
	userCollection.Insert(User{
		Username: "Zuolar",
		Password: "qwe123",
	})

	var res db.Result
	res = userCollection.Find()
	var users []User

	c.String(http.StatusOK, "Selecting User...\n")
	err = res.All(&users)
	if err != nil {
		log.Fatalf("res.All(): %q\n", err)
	}

	for _, user := range users {
		c.String(http.StatusOK, fmt.Sprintf("%d: %s => %s.\n",
			user.ID,
			user.Username,
			user.Password,
		))
	}

}
