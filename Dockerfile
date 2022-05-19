# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /forum
COPY . .
RUN  go mod download

RUN go build -o /app .

EXPOSE 8080

CMD [ "/app" ]