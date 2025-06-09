# Portfolio Site - Project Overview

## Tech Stack

### Backend
- **Go (Golang)** - Web server using `net/http` and `html/template`
- **Gorilla Mux** - HTTP router for clean URL routing

### Frontend
- **Tailwind CSS** - Utility-first CSS framework
- **Bun** - JavaScript runtime for build tools
- **Tailwind CLI** - For CSS compilation
- **Shiki.js** - Syntax highlighting for code blocks
- **Mermaid.js** - Diagram rendering for architecture diagrams

### Build Tools
- **Bun** - Package manager and build tool
- **Tailwind CLI** - CSS processing
- **Go modules** - Dependency management

## Project Structure

```
portfolio/
├── .rules/
│   └── project-overview.md
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── home.go             # Landing page handler
│   │   └── blog.go             # Blog handlers
│   ├── models/
│   │   └── blog.go             # Blog post models
│   └── server/
│       └── server.go           # Server setup and routing
├── web/
│   ├── templates/
│   │   ├── base.html           # Base template with common layout
│   │   ├── home.html           # Landing page template
│   │   ├── blog-list.html      # Blog listing template
│   │   └── blog-post.html      # Individual blog post template
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css       # Compiled Tailwind CSS
│   │   └── js/
│   │       ├── mermaid.min.js  # Mermaid diagram library
│   │       └── shiki.min.js    # Shiki syntax highlighter
│   └── content/
│       └── blogs/
│           ├── devops-k8s-architecture.md
│           └── go-concurrency-patterns.md
├── tailwind.config.js          # Tailwind configuration
├── package.json                # Bun/Node dependencies
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies checksum
└── README.md                   # Project documentation
```

## Features

### Landing Page
- Hero section with centered title and subtitle
- Clean, minimal design using Tailwind CSS
- Responsive layout
- Navigation to blog posts

### Blog System
1. **DevOps Blog Post**
   - Topic: Kubernetes Architecture on AWS
   - Rendered Mermaid diagrams showing cloud architecture
   - Technical content about container orchestration

2. **Backend Go Blog Post**
   - Topic: Go Concurrency Patterns
   - Syntax-highlighted Go code blocks using Shiki.js
   - Practical examples and explanations

### Technical Implementation
- Server-side rendering with Go templates
- Static asset serving for CSS/JS
- Markdown parsing for blog content
- Client-side diagram rendering with Mermaid
- Syntax highlighting with Shiki.js

## Development Workflow

1. **Setup**: Install dependencies with `bun install`
2. **CSS Build**: Generate Tailwind CSS with `bunx tailwindcss -i input.css -o web/static/css/style.css --watch`
3. **Server**: Run Go server with `go run cmd/server/main.go`
4. **Development**: Hot reload for CSS changes, manual restart for Go changes

## Deployment Considerations

- Single binary deployment (Go compiled executable)
- Static assets bundled or served separately
- Environment-based configuration
- Docker containerization ready
- Can be deployed to any platform supporting Go binaries