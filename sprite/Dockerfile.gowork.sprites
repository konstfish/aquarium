FROM golang:1.22-alpine as builder

ARG COMPONENT

WORKDIR /app

COPY $COMPONENT/go.mod $COMPONENT/go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -v -o main $COMPONENT/main/main.go

FROM alpine:latest

ENV GIN_MODE=release

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /app/main .

ADD sprite/sprites /sprites

CMD ["./main"]