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
	EditView *views.View
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
		EditView: views.NewView("bootstrap", "galleries/edit"),
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

func (g *Galleries) galleryByID(w http.ResponseWriter,
	r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "you do not have permission to edit this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	var form GalleryForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	gallery.Title = form.Title

	err = g.gs.Update(gallery)

	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery updated successfully",
		}
	}

	g.EditView.Render(w, vd)

}
