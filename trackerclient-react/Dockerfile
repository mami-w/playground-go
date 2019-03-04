################### stage 1
FROM node:latest as build-stage

WORKDIR /tmp
COPY package.json .
COPY .babelrc .
COPY webpack.config.js .
COPY App/ app/
COPY assets/index.html assets/

RUN npm install -i
#RUN pwd; ls -la
RUN npm run build

#################   stage 2
FROM golang

WORKDIR /app
RUN go get -d -v github.com/gorilla/mux

ENV SRC_DIR=/go/src/github.com/mami-w/playground-go/trackerclient-react/
WORKDIR $SRC_DIR
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/gorilla/mux
COPY server server
COPY main.go .
RUN cd $SRC_DIR; go build -o trackerclient-react ./main.go; cp trackerclient-react /app/

WORKDIR /app
COPY --from=build-stage /tmp/assets/ assets/
COPY trackerclient.sh .
ENTRYPOINT ["./trackerclient.sh"]

#build with 'docker build -t fwaupo/trackerclient-react .'
#run with 'docker run un -p 8006:8006 [--env-file <env file>] [-it] [-d] fwaupo/trackerclient:latest

# Stage 1, based on Nginx, to have only the compiled app, ready for production with Nginx
#FROM nginx:1.15
#COPY --from=build-stage /app/build/ /usr/share/nginx/html
# Copy the default nginx.conf provided by tiangolo/node-frontend
#COPY --from=build-stage /nginx.conf /etc/nginx/conf.d/default.conf


#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
#CMD ["./app"]
