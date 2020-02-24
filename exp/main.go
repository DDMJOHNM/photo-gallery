package main

import (
	"fmt"

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
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Order{})

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()
	us.AutoMigrate()

	var user User

	db.First(&user)
	if db.Error != nil {
		panic(db.Error)
	}

	//createOrder(db, user, 1001, "Fake Description #1")
	//createOrder(db, user, 9999, "Fake Description #2")
	//createOrder(db, user, 8800, "Fake Description #3")
}
func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID: user.ID, Amount: amount, Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
