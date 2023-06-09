FROM golang:alpine as builder
WORKDIR /event-history
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o event-history cmd/*.go

FROM scratch
COPY --from=builder /event-history/event-history .
ENTRYPOINT ["./event-history","http-serve"]
