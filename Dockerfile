
FROM  golang:1.19



COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build cmd/main.go

CMD [ "/order-portal" ]