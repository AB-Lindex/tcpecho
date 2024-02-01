SOURCES=*.go

build: bin/tcpecho

bin/tcpecho: $(SOURCES) Makefile go.mod
	@mkdir -p bin
	go build -o bin/tcpecho .

run: bin/tcpecho
	bin/tcpecho

runc:
	docker run --rm -it tcpecho:docker

docker:
	-docker rmi tcpecho:docker
	docker build -t tcpecho:docker .

nerdctl:
	-docker rmi tcpecho:nerdctl
	nerdctl build -t tcpecho:nerdctl .

check:
	@echo "Checking...\n"
	gocyclo -over 15 . || echo -n ""
	@echo ""
	golint -min_confidence 0.21 -set_exit_status ./...
	@echo ""
	go mod verify
	@echo "\nAll ok!"

check2:
	@echo ""
	golangci-lint run -E misspell -E depguard -E dupl -E goconst -E gocyclo -E ifshort -E predeclared -E tagliatelle -E errorlint -E godox -D structcheck

release:
	gh release create `cat version.txt`
