FROM golang:1.22.3
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscriber .
ENV CONFPATH=/app/data
CMD ["/subscriber"]