run:
	docker run -d -p 6379:6379 --name redis-jkjav redis
	go run .

test-submit:
	curl -i -X POST localhost:3000/submit -d \
	'{"mysejahteraId": "abc123", "firstName": "Larry", "lastName": "Ellison", "location": "PWTCx", "phoneNumber": "0123456789", "date": "20210601"}' \
	-H "Content-type: application/json"

# this requires hey
load-test:
	hey -c 10 -z 1s http://localhost:3000/submit

redis-cli:
	 docker exec -it redis-jkjav redis-cli
