FROM golang:1.19-alpine AS builder

WORKDIR /build

ADD go.mod .

ENV GO111MODULE=on

COPY . .

# download Go modules and dependencies
RUN go mod download

RUN go build -o crm-worker .

FROM alpine

WORKDIR /build

COPY --from=builder /build/crm-worker /build/crm-worker
COPY .env /build
COPY gcredentials.json /build

CMD ["./crm-worker"]