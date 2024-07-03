FROM golang:1.22-alpine as builder
WORKDIR /app

COPY . .

# @note: go mod tidy is not recommended for live projects. use: go mod download
RUN go mod tidy && go mod verify
RUN CGO_ENABLED=0 go build -o service cmd/service/main.go

FROM gcr.io/distroless/static-debian12 as runner
WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/service .
COPY --from=builder --chown=nonroot:nonroot /app/.env .
COPY --from=builder --chown=nonroot:nonroot /app/cert/* ./cert/

EXPOSE 3000

ENTRYPOINT ["./service"]