.PHONY: dev up down restart

# Re-create volumes from scratch and start everything
dev:
	docker-compose down -v
	docker-compose up --build

# Standard startup
up:
	docker-compose up

# Stop services
down:
	docker-compose down

# Restart without wiping volumes
restart:
	docker-compose restart
