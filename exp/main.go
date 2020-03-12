package main

import (
	"fmt"
	"../models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "testgallerydb"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
		us, err := models.NewUserService(psqlInfo)
		if err != nil {
			panic(err)
		}
		defer us.Close()
		us.DestructiveReset()
		user := models.User{
			Name: "Michael Scott",
			Email: "michael@dundermifflin.com",
			Password: "bestboss",
		}
		err = us.Create(&user)
		if err != nil {
			panic(err)
		}
		// Verify that the user has a Remember and RememberHash
		fmt.Printf("%+v\n", user)
		if user.Remember == "" {
			panic("Invalid remember token")
		}
		// Now verify that we can lookup a user with that remember
		// token
		user2, err := us.ByRemember(user.Remember)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", *user2)

}