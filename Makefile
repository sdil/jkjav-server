run:
	docker-compose build
	docker-compose down
	docker-compose up
	# go run .

load-get-test:
	hey -c 100 -z 5s http://localhost:3000/stations/PWTC

load-post-test:
	hey -T "application/json" \
	-c 100 \
	-z 5s \
	-m POST \
	-d '{"address":"Kuala Lumpur","date":"20210516","firstName":"Fadhil","lastName":"Yaacob","location":"PWTC","mysejahteraId":"900127015527","phoneNumber":"0123456789"}' \
	http://localhost:3000/booking

redis-cli:
	 docker-compose exec -it redis redis-cli

test:
	go test -v ./api/routes/
