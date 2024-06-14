.PHONY: run build clean

run:
	CompileDaemon --build="go build -o hex cmd/hex/main.go" --command=./hex

build:
	go build -o hex cmd/hex/main.go

clean:
	rm -f hex