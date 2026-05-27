FROM golang:1.26.3-alpine
COPY . .
RUN go build -o hobby-lobby .
ENTRYPOINT ["./hobby-lobby"]
