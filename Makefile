build:
	go build main.go -i configz

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o configz ./main.go
