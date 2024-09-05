gomod:
	go mod tidy

tailwind:
	npx tailwindcss build -i static/index.css -o static/tailwind.css

templ: gomod tailwind
	templ generate

build: templ
	go build -o collage cmd/main.go
