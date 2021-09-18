run:
	docker run -d -p 6379:6379 --name redis-jkjav redis
	go run .

load-test:
	hey -T "application/json" \
	-c 100 \
	-z 5s \
	-m POST \
	-d '{"address":"Kuala Lumpur","date":"20210516","firstName":"Fadhil","lastName":"Yaacob","location":"PWTC","mysejahteraId":"900127015527","phoneNumber":"0123456789"}' \
	http://localhost:3000/booking

redis-cli:
	 docker-compose exec -it redis redis-cli
