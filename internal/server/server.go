package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"portfolio/internal/handlers"

	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	templates *template.Template
}

func New() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	s.loadTemplates()
	s.setupRoutes()
	return s
}

func (s *Server) loadTemplates() {
	var err error
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	s.templates, err = template.New("").Funcs(funcMap).ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}
}

func (s *Server) setupRoutes() {
	// Static file serving
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Handlers
	homeHandler := handlers.NewHomeHandler(s.templates)
	blogHandler := handlers.NewBlogHandler(s.templates)

	// Routes
	s.router.Handle("/", homeHandler).Methods("GET")
	s.router.HandleFunc("/blog", blogHandler.ListPosts).Methods("GET")
	s.router.HandleFunc("/blog/{slug}", blogHandler.ViewPost).Methods("GET")

	// Add middleware
	s.router.Use(s.loggingMiddleware)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start(port string) error {
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s.router)
}
