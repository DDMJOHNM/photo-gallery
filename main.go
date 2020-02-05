package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"./views"
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

func main (){

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