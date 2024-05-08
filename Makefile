build:
	@go build -C src/ -o "../manager"

clean:
	@rm manager
