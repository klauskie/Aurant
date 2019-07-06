FROM golang:1.9
WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -d -v github.com/gorilla/mux
RUN go get -d -v github.com/go-sql-driver/mysql
CMD ["go","run","main.go"]
