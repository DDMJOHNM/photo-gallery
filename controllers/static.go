package controllers

import "../views"

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "home"),
		Contact: views.NewView("bootstrap", "contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
