FROM golang:alpine
ADD . /go/src/leogregianin/brcep
RUN cd /go/src/leogregianin/brcep && go install
CMD ["/go/bin/brcep"]
EXPOSE 8000
