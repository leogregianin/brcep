FROM golang:alpine
ADD . /go/src/leogregianin/brcep
RUN go install leogregianin/brcep
CMD ["/go/bin/brcep"]
EXPOSE 8000
