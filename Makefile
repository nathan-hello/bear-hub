dev/templ:
	templ generate -path src/components --watch

dev/air:
	go run github.com/cosmtrek/air@v1.52.0

dev/tailwind:
	bunx tailwindcss -i ./public/css/tw-base.css -o ./public/css/tw-output.css --minify --watch

dev: 
	make -j3 dev/templ dev/air dev/tailwind

test/chat:
	./scripts/start_chat.sh
	
