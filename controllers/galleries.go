package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"../context"
	"../models"
	"../views"
	"github.com/gorilla/mux"
)

const (
	ShowGallery = "show_gallery"
)

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		gs:       gs,
	}
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		fmt.Print("There was an error")
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	if r.Method == "GET" {
		url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Redirect(w, r, url.Path, http.StatusNotFound)
	}
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idstr := vars["id"]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "invalid gallery id", http.StatusNotFound)
	}

	//_ = id

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	// gallery := models.Gallery{
	// 	Title: "A temporary fake gallery with ID: " + idstr,
	// }

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}
