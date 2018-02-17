FROM golang:alpine
ADD . /go/src/zeit/brcep
RUN go install zeit/brcep
CMD ["/go/bin/brcep"]
EXPOSE 8000
