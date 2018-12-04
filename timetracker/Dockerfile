FROM golang
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/mami-w/timetracker/
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go build -o timetracker ./main/main.go; cp timetracker /app/
ENTRYPOINT ["./timetracker"]

# build with docker build -t mami-w/myimage-go .
# run with docker run -p 8000:8000 mami-w/myimage-go:latest