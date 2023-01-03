FROM golang:1.19-alpine

WORKDIR ./go-aws-ec2
COPY . .

RUN go build -o ./build/go-aws-ec2 ./cmd/go-aws-ec2/main.go
CMD ["./build/go-aws-ec2"]
