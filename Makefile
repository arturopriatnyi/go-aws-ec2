all: build run

build:
	docker build . -t go-aws-ec2

run:
	docker run -p 10000:10000 go-aws-ec2
