#Backend Building Stage
FROM golang:alpine3.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd

#Running stage
FROM golang:alpine3.23
WORKDIR /app
COPY --from=builder /app/server .

#Create folder for sqlite
RUN mkdir -p /app/data

#RUN
EXPOSE 8080
CMD ["./server"]
