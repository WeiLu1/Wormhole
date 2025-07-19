run-local:
	set -a && source .env && go run main.go $(file)

run-local-with-sample:
	set -a && source .env && go run main.go sample.yaml

test:
	go test ./...