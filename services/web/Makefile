run:
	@make -j2 gow tailwind

tailwind:
	./tailwindcss -i ./assets/main.css -o ./static/main.css --watch

gow:
	ENV="dev" ~/go/bin/gow -e=go,mod,html -s run ../../cmd/web
