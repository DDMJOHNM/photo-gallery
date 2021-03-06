package controllers

import (
	"fmt"
	"net/http"
	"time"

	"../context"
	"../models"
	"../rand"
	"../views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)

}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		fmt.Println(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that address")
		default:
			vd.SetAlert(err)
		}
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		fmt.Println(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	alert := views.Alert{
		Level: views.AlertLvlSuccess, Message: "Welcome to your Photo Gallery",
	}

	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}

		cookie := http.Cookie{
			Name:     "remember_token",
			Value:    user.Remember,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}

	return nil
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, user)
}

func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	user := context.User(r.Context())

	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)

}
