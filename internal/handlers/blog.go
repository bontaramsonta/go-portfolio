package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"portfolio/internal/models"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type BlogHandler struct {
	templates *template.Template
	posts     []models.BlogPost
}

func NewBlogHandler(templates *template.Template) *BlogHandler {
	handler := &BlogHandler{
		templates: templates,
	}
	handler.loadBlogPosts()
	return handler
}

func (h *BlogHandler) loadBlogPosts() {
	contentDir := "web/content/blogs"
	files, err := os.ReadDir(contentDir)
	if err != nil {
		fmt.Printf("Error reading blog directory: %v\n", err)
		return
	}

	var posts []models.BlogPost
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			post, err := h.parseBlogPost(filepath.Join(contentDir, file.Name()))
			if err != nil {
				fmt.Printf("Error parsing blog post %s: %v\n", file.Name(), err)
				continue
			}
			if post.Published {
				posts = append(posts, post)
			}
		}
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	h.posts = posts
}

func (h *BlogHandler) parseBlogPost(filePath string) (models.BlogPost, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return models.BlogPost{}, err
	}

	// Split frontmatter and content
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return models.BlogPost{}, fmt.Errorf("invalid frontmatter format")
	}

	// Parse frontmatter
	var metadata models.BlogMetadata
	err = yaml.Unmarshal([]byte(parts[1]), &metadata)
	if err != nil {
		return models.BlogPost{}, fmt.Errorf("error parsing frontmatter: %v", err)
	}

	// Parse date
	date, err := time.Parse("2006-01-02", metadata.Date)
	if err != nil {
		return models.BlogPost{}, fmt.Errorf("error parsing date: %v", err)
	}

	// Convert markdown to HTML
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(parts[2]), &buf); err != nil {
		return models.BlogPost{}, fmt.Errorf("error converting markdown: %v", err)
	}

	// Generate slug from filename
	filename := filepath.Base(filePath)
	slug := strings.TrimSuffix(filename, ".md")

	post := models.BlogPost{
		ID:            slug,
		Title:         metadata.Title,
		Slug:          slug,
		Excerpt:       metadata.Excerpt,
		Content:       buf.String(),
		Author:        metadata.Author,
		Date:          date,
		Tags:          metadata.Tags,
		Category:      metadata.Category,
		ReadTime:      metadata.ReadTime,
		Published:     metadata.Published,
		HasMermaid:    metadata.HasMermaid,
		HasCodeBlocks: metadata.HasCodeBlocks,
	}

	return post, nil
}

func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Posts []models.BlogPost
		Title string
	}{
		Posts: h.posts,
		Title: "Blog Posts",
	}

	err := h.templates.ExecuteTemplate(w, "blog-list.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *BlogHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	var post *models.BlogPost
	for _, p := range h.posts {
		if p.Slug == slug {
			post = &p
			break
		}
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Post  models.BlogPost
		Title string
	}{
		Post:  *post,
		Title: post.Title,
	}

	err := h.templates.ExecuteTemplate(w, "blog-post.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
