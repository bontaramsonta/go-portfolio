package models

import (
	"time"
)

type BlogPost struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Excerpt     string    `json:"excerpt"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	Category    string    `json:"category"`
	ReadTime    int       `json:"read_time"` // in minutes
	Published   bool      `json:"published"`
	HasMermaid  bool      `json:"has_mermaid"`
	HasCodeBlocks bool    `json:"has_code_blocks"`
}

type BlogMetadata struct {
	Title       string   `yaml:"title"`
	Excerpt     string   `yaml:"excerpt"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Category    string   `yaml:"category"`
	ReadTime    int      `yaml:"read_time"`
	Published   bool     `yaml:"published"`
	HasMermaid  bool     `yaml:"has_mermaid"`
	HasCodeBlocks bool   `yaml:"has_code_blocks"`
}