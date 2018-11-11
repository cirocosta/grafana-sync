build:
	go build -v -i

run:
	./grafana-sync \
		--verbose \
		--access-token=eyJrIjoiMHJZYlEyd1M2dXZON01hTmdQTG5kR29wVjJZUDhhVjciLCJuIjoidmlld2VyIiwiaWQiOjF9 \
		--directory=/tmp
