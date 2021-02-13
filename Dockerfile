FROM golang:1.15.8-buster

WORKDIR /app
COPY . ./

RUN make deps
RUN make lint
RUN make build

EXPOSE 8080

CMD ["/app/bin/svc"]
