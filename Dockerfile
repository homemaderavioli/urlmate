FROM golang:1.14

WORKDIR /app

COPY . .

#RUN go get -d -v ./..
#RUN go install -v ./..
RUN go build -o url .

CMD ["/app/url"]