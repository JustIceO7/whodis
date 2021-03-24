FROM golang:1.16-alpine as doc

WORKDIR /doc

COPY . .

# Generate OpenAPI docs

RUN go get github.com/swaggo/swag/cmd/swag

RUN $GOPATH/bin/swag init -g internal/api/api.go


FROM golang:1.16-alpine AS dev

WORKDIR /bot

COPY --from=doc /doc .

RUN apk add git

RUN GO111MODULE=on go get github.com/cortesi/modd/cmd/modd

COPY go.mod .
COPY go.sum .

RUN go mod download

# Compile
RUN go install github.com/uccnetsoc/veribot/cmd/veribot

CMD ["go", "run", "*.go"]


FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/veribot ./veribot

CMD ["sh", "-c", "veribot"]