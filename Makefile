build:
	@tailwindcss -i src/public/css/styles.css -o src/public/styles.css
	@templ generate view

tailwind:
	./tailwindcss -i src/views/css/styles.css -o src/public/styles.css --watch

tailwind--min:
	./tailwindcss -i src/views/css/styles.css -o src/public/styles.css --minify

templ:
	@templ generate -watch


