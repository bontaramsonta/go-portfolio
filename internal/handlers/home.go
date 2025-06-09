package handlers

import (
	"html/template"
	"net/http"
)

type HomeHandler struct {
	templates *template.Template
}

func NewHomeHandler(templates *template.Template) *HomeHandler {
	return &HomeHandler{
		templates: templates,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title    string
		Subtitle string
		Year     int
	}{
		Title:    "DevOps Engineer / Full Stack Developer",
		Subtitle: "I don't have time to create an actual portfolio but check out my devlogs and projects",
		Year:     2024,
	}

	err := h.templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}