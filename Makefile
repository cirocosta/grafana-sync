build:
	go build -v -i

run:
	./grafana-sync \
		--verbose \
		--username=admin \
		--password=admin \
		--directory=/tmp
