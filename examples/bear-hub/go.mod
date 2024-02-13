module github.com/nathan-hello/no-magic-stack/examples/bear-hub

go 1.21

require (
	github.com/a-h/templ v0.2.543
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/joho/godotenv v1.5.1
	github.com/justinas/alice v1.2.0
	github.com/lib/pq v1.10.9
	github.com/nathan-hello/htmx-template v0.0.0-20240212150651-a58549fa8023
	golang.org/x/crypto v0.19.0
)

require golang.org/x/net v0.21.0 // indirect

replace github.com/nathan-hello/no-magic-stack/examples/bear-hub => ./
