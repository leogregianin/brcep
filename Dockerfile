FROM alpine:latest

MAINTAINER Leonardo Gregianin <leogregianin@gmail.com>

WORKDIR "/opt"

ADD .docker_build/brcep /opt/bin/brcep

CMD ["/opt/bin/brcep"]

