package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"./views"

	"database/sql"
	"fmt"
	_"github.com/lib/pq"

	)


var (
	homeView *views.View
 	contactView *views.View
	signupView *views.View
	)

func home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	must(homeView.Render(w,nil))

}

func contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	must(contactView.Render(w,nil))
}

func signup(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	must(signupView.Render(w,nil))
}


func must(err error)  {
	if err != nil{
		panic(err)
	}
}

const (
	host = "localhost"
	port = "5432"
	user = "postgres"
	password="postgres"
	dbname="testgallerydb"
)

func main (){

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db ,err := sql.Open("postgres",psqlInfo)
	if err != nil{
		panic(err)
	}
	err = db.Ping()
	if err != nil{
		panic(err)
	}

	fmt.Print("Db successfully connected")
	db.Close()


	//var err error
	homeView =	views.NewView("bootstrap","views/home.gohtml")
	contactView = views.NewView("bootstrap","views/contact.gohtml")
	signupView = views.NewView("bootstrap","views/signup.gohtml")

	var h http.Handler = http.HandlerFunc(home)
	r:= mux.NewRouter()
	r.HandleFunc("/",home)
	r.HandleFunc("/contact",contact)
	r.HandleFunc("/signup",signup)
	r.NotFoundHandler = h
	http.ListenAndServe(":3000",r)
}