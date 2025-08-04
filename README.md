# Blogly

A simple blogging platform built with Next.js and Go Gin framework.

## Features

- Create, edit, and delete blog posts
- User authentication
- Responsive design

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/blogly.git
   cd blogly
   ```
2. Install dependencies for Next.js:
   ```bash
   cd frontend
   bun install
   ```
3. Install dependencies for Go Gin backend:
   ```bash
   cd ../backend
   go mod tidy
   ```
4. Run the backend server:
   ```bash
   go run main.go
   ```
5. Run the frontend:
   ```bash
   cd ../frontend
   bun run dev
   ```

## License

MIT
