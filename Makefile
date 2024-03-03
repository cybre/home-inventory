destroy:
	docker-compose down -v --rmi all
rebuild:
	docker-compose down
	docker-compose up --build
rebuild-inventory:
	docker-compose up --build -d --no-deps inventory
restart-inventory:
	docker-compose up -d --no-deps inventory