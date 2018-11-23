build:
	go install -v

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
