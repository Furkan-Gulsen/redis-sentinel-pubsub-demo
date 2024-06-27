COMPOSE_FILE := docker-compose.yaml

.PHONY: start stop clean redis-mock-data redis-benchmark

start:
	docker-compose -f $(COMPOSE_FILE) up -d --build

stop:
	docker-compose -f $(COMPOSE_FILE) stop

clean:
	docker-compose -f $(COMPOSE_FILE) down

redis-mock-data:
	chmod +x redis-mock-data.sh
	./redis-mock-data.sh

redis-benchmark:
	redis-benchmark -h localhost -p 6383 -c 1000 -n 1000000
