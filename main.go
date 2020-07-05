package main

import (
	"fmt"
	"net/http"
	"./rand"

	"./controllers"
	"./hash"
	"./middleware"
	"./models"
	"./views"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	home    *views.View
	contact *views.View
	//signupView  *views.View
)

// func signup(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(signupView.Render(w, nil))
// }

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

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "secret"
// 	dbname   = "testgallerydb"
// )

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

type PostgresConfig struct{
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Name string `json:"name"`
	Pepper: "secret-random-string", 
	HMACKey: "secret-hmac-key",
}

type Config struct{
	Port int
	Env string
	Pepper string `json:"pepper"` 
	HMACKey string `json:"hmac_key"`
}

func (c Config) isProd() bool{
	return c.Env == "prod"
}

func DefaultConfig() Config{
	return Config{
		Port :3000,
		Env : "dev",
	}
}

func (c PostgresConfig) Dialect() string{
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string{
	if c.Password == ""{
		return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
	}

	 return fmt.Sprintf("host=%s port=%d user=%s "+
	 "password=%s dbname=%s sslmode=disable",
	 c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host: "localhost",
		Port: 5432,	
		User: "postgres",
		Password: "secret",
		Name: "testgallerydb",
	}
}

func main() {

	//fmt.Println(rand.String(10))
	//fmt.Println(rand.RememberToken())

	cfg:= DefaultConfig()
	dbCfg := DefaultPostgresConfig()

	hmac := hash.NewHMAC("my-secret-key")

	fmt.Println(hmac.Hash("this is my string to hash"))

	//  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//  	"password=%s dbname=%s sslmode=disable",
	//  	host, port, user, password, dbname)

	// us, err := models.NewUserService(psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer us.Close()
	// us.DestructiveReset()

	services, err := models.NewServices(dbCfg.Dialect(),dbCfg.ConnectionInfo())
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	defer services.Close()
	services.AutoMigrate()

	home = views.NewView("bootstrap", "home")
	contact = views.NewView("bootstrap", "contact")
	staticC := controllers.NewStatic()
	//usersC := controllers.NewUsers(us)
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	userMw := middleware.User{
		UserService: services.User,
	}

	requireUserMw := middleware.RequireUser{}

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	isProd := false
	b, err := rand.Bytes(32)
	if err != nil{
		panic(err)
	}

	csrfMw := csrf.Protect(b,csrf.Secure(isProd))

	//var h http.Handler = http.HandlerFunc(home)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	//r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	// r.Handle("/galleries",
	// 	requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries",
		requireUserMw.ApplyFn(galleriesC.Index)).
		Methods("GET").
		Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesC.Edit)).
		Methods("GET").
		Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")
	assetHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/assets/").Handler(assetHandler)
	
	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(":3000", csrfMw(userMw.Apply(r)))

}

//632
