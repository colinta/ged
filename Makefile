.PHONY: build test readme clean

build:
	go build -o ged ./cmd/ged

test:
	go test ./...

readme:
	./scripts/generate-readme.sh > README.md

clean:
	rm -f ged
