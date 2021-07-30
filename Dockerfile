FROM golang:1.17-rc

RUN apt update
RUN apt-get update

ENTRYPOINT [ "/bin/bash" ]