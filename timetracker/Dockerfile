FROM golang
WORKDIR /app
RUN go get -d -v github.com/lib/pq
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/mami-w/playground-go/timetracker/
# Add the source code:
ADD . $SRC_DIR
COPY tracker.sh /app
# Build it:
RUN cd $SRC_DIR; go build -o timetracker ./main/main.go; cp timetracker /app/
ENTRYPOINT ["./tracker.sh"]

# build with docker build -t mami-w/myimage-go .
# run with docker run -p 8000:8000 mami-w/myimage-go:latest

#FROM golang:1.8.3 as builder
#WORKDIR /go/src/github.com/flaviocopes/findlinks
#RUN go get -d -v golang.org/x/net/html
#COPY findlinks.go  .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o findlinks .

#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/flaviocopes/findlinks/findlinks .
#CMD ["./findlinks"]