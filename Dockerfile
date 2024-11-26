FROM golang:1.23-alpine
EXPOSE 44317

WORKDIR /
COPY . /
RUN go build .
CMD ["/yikes"]
