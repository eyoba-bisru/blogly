# Blogly Backend

This is the backend service for the Blogly project, built with the [Gin](https://gin-gonic.com/) Go web framework.

## Setup

1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. For development, use [Air](https://github.com/cosmtrek/air) for live reloading:
   ```bash
   go install github.com/cosmtrek/air@latest
   air
   ```
4. Configure environment variables as needed.

## Running the Server

```bash
go run main.go
```

## License

MIT
