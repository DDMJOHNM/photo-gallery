package main

import (
	"net/http"

	"./models"
	"./views"
	"github.com/gorilla/mux"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	homeView    *views.View
	contactView *views.View
	signupView  *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}

const (
	host     = "localhost"
	port     = 5555
	user     = "postgres"
	password = "postgres"
	dbname   = "testgallerydb"
)

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
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

	/*user := models.User{
		Name:  "Michael Scott",
		Email: "michael@test.com",
	}

	/*if err := us.Create(&user); err != nil {
		fmt.Print(err.Error())
	}

	user.Name = "Updated Name"
	if err := us.Update(&user); err != nil {
		panic(err)
	}

	founduser, err := us.ByID(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(founduser)

	/*foundUser, err := us.ByEmail("michael@test.com")
	if err != nil {
		panic(err)
	}
	// Because of an update, the name should now // be "Updated Name"
	fmt.Println(foundUser)*/

	//usersC := controllers.NewUsers(us)

	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	var h http.Handler = http.HandlerFunc(home)
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", signup)
	r.NotFoundHandler = h
	http.ListenAndServe(":3000", r)
}

//201
