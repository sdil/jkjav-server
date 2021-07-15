run:
	docker run -d -p 6379:6379 redis
	go run .
