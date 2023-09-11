# MULTI STAGE BUILD

# builder stage
FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


FROM alpine:3.17 as RUNNER
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .

EXPOSE 3000
CMD [ "/app/start.sh","/app/main" ]
