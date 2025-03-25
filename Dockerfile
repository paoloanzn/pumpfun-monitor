FROM golang:1.24.1-alpine

WORKDIR /app

RUN echo "fs.file-max = 1000000" >> /etc/sysctl.conf

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pumpfun-monitor ./src

EXPOSE 8080

CMD ["./pumpfun-monitor", "start", "-mint-workers=5", "-migration-workers=5"]
