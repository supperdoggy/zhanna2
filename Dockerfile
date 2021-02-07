FROM golang:1.16rc1-buster
COPY flowers/ /usr/src/myapp
EXPOSE 80
