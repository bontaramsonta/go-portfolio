# Portfolio Site

A simple portfolio website built with Go, Tailwind CSS, and modern web technologies.

## Tech Stack

- **Backend**: Go with Gorilla Mux router
- **Frontend**: Tailwind CSS, Mermaid.js, Shiki.js
- **Build Tools**: Bun, Tailwind CLI
- **Content**: Markdown with YAML frontmatter

## Features

- Clean, responsive design with Tailwind CSS
- Blog system with markdown content
- Syntax highlighting for code blocks (Shiki.js)
- Diagram rendering with Mermaid.js
- Server-side rendering with Go templates

## Setup

### Prerequisites

- Go 1.21 or later
- Bun (JavaScript runtime)

### Installation

1. Clone the repository and navigate to the portfolio directory:

```bash
cd portfolio
```

2. Install JavaScript dependencies:

```bash
bun install
```

3. Install Go dependencies:

```bash
go mod tidy
```

4. Build CSS (for production):

```bash
bun run build:css:prod
```

Or for development with watch mode:

```bash
bun run build:css
```

## Running the Server

### Quick Start

```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default. Visit http://localhost:8080 to view the site.

To specify a different port:

```bash
PORT=3000 go run cmd/server/main.go
```

## Project Structure

```
portfolio/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   └── server/                 # Server setup
├── web/
│   ├── templates/              # HTML templates
│   ├── static/                 # CSS, JS, assets
│   └── content/blogs/          # Blog posts
├── go.mod                      # Go dependencies
├── package.json                # Node/Bun dependencies
└── tailwind.config.js          # Tailwind configuration
```

## Adding Blog Posts

Create new markdown files in `web/content/blogs/` with YAML frontmatter:

```markdown
---
title: "Your Post Title"
excerpt: "Brief description"
author: "Your Name"
date: "2024-01-15"
tags: ["tag1", "tag2"]
category: "Category"
read_time: 5
published: true
has_mermaid: false
has_code_blocks: true
---

# Your content here...
```

## Development

1. Start CSS watch mode:

```bash
bun run build:css
```

2. In another terminal, start the Go server:

```bash
go run cmd/server/main.go
```

3. Make changes to templates, CSS, or Go code
4. CSS changes are automatically rebuilt
5. Go changes require server restart

## Deployment

1. Build production CSS:

```bash
bun run build:css:prod
```

2. Build Go binary:

```bash
go build -o portfolio cmd/server/main.go
```

3. Run the binary:

```bash
./portfolio
```

## Environment Variables

- `PORT`: Server port (default: 8080)

## Sample Content

The project includes two sample blog posts:

1. **DevOps**: Kubernetes Architecture on AWS with Mermaid diagrams
2. **Backend**: Go Concurrency Patterns with syntax-highlighted code
