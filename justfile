build:
    go build -o ged ./cmd/ged

test:
    go test ./...

readme:
    ./scripts/generate-readme.sh > README.md

man: readme
    npx marked-man README.md > ged.1

docs: readme man

clean:
    rm -f ged ged.1
