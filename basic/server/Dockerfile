FROM golang
COPY main.go /
RUN go build -o /main /main.go
EXPOSE 8000
ENTRYPOINT ["/main"]

# build with docker build -t myimage-go .
# run with docker run -p 8000:8000 myimage-go:latest