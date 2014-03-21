docker:
	GOOS=linux go build
	docker build -t rochacon/s3stub .
