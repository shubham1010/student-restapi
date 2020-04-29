FROM golang:alpine as builder

LABEL maintainer = "Satish Sangam"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM mariadb:10.3

WORKDIR /root
ENV MYSQL_ROOT_PASSWORD root
ENV MYSQL_USER satish
ENV MYSQL_PASSWORD #Ss77
ENV MYSQL_DATABASE demo

COPY migrations /root/migrations
COPY databaseConfig /root/databaseConfig
COPY startup.sh /root/startup.sh
RUN chmod +x startup.sh

COPY --from=builder /app/main .
EXPOSE 8080
