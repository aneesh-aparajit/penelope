build:
	@go build -o ./build/penelope .

run: build
	@./build/penelope
