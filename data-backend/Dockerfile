FROM golang:alpine as buildData

WORKDIR /usr/data-backend

COPY . .

RUN go build

CMD ["./data-backend"]