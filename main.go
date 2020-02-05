package main

import(
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	)


var homeTemplate *template.Template
var contactTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	if err := homeTemplate.Execute(w,nil); err != nil{
		panic(err)
	}

}

func contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	if err := contactTemplate.Execute(w,nil); err != nil{
		panic(err)
	}
}



func main (){

	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if  err != nil{
		panic(err)
	}

	contactTemplate, err = template.ParseFiles("views/contact.gohtml")
	if  err != nil{
		panic(err)
	}

	var h http.Handler = http.HandlerFunc(home)
	r:= mux.NewRouter()
	r.HandleFunc("/",home)
	r.HandleFunc("/contact",contact)
	r.NotFoundHandler = h
	http.ListenAndServe(":3000",r)
}