FROM golang:alpine
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /app
ADD . .
RUN go build -o main .
EXPOSE 8100
EXPOSE 8200
EXPOSE 8300
CMD ["/app/main"]