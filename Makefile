include .env

db:
	rm database.db && touch database.db && sqlite3 database.db < src/db/schema.sql
	DB_URI=${DB_URI} go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate 

chat:
	./chat.sh

dev/templ:
	templ generate -path src/components --watch

dev/air:
	go run github.com/cosmtrek/air@v1.52.0

dev/tailwind:
	bunx tailwindcss -i ./public/css/tw-base.css -o ./public/css/tw-output.css --minify --watch

dev: 
	make -j3 dev/templ dev/air dev/tailwind

