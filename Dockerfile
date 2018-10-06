FROM golang:1.11

RUN useradd -u 10001 myapp

RUN mkdir -p /go/src/github.com/skhvan1111/go-first
ADD . /go/src/github.com/skhvan1111/go-first
WORKDIR /go/src/github.com/skhvan1111/go-first

RUN CGO_ENABLED=0 go build -a -installsuffix cgo \
    -o bin/go-first github.com/skhvan1111/go-first/cmd/go-first

FROM scratch
ENV PORT 8080
ENV DIAG_PORT 8585

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=0 /etc/passwd /etc/passwd
USER myapp

COPY --from=0 /go/src/github.com/skhvan1111/go-first/bin/go-first /go-first
EXPOSE $PORT
EXPOSE $DIAG_PORT

CMD ["/go-first"]