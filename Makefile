build:
	go build -v -i

test:
	go test ./...

run:
	./grafana-sync \
		--verbose \
		--username=admin \
		--password=admin \
		--directory=/tmp
