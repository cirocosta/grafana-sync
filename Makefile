build:
	go build -v -i

test:
	go test ./...

fmt:
	go fmt ./...

run:
	./grafana-sync \
		--verbose \
		--username=admin \
		--password=admin \
		--directory=/tmp
