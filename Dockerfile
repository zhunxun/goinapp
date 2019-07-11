FROM golang:1.12-alpine AS TEST

RUN apk update && \
    apk add --no-cache git gcc && \
    mkdir -p /goinapp

WORKDIR /goinapp

ADD go.mod .
#ADD go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go test -сover -race -v ./... && go build -o goinapp



FROM golang:1.12-alpine AS BUILD

RUN apk update && \
    apk add --no-cache git && \
    mkdir -p /goinapp

WORKDIR /goinapp

ADD go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goinapp



FROM scratch

COPY --from=BUILD /goinapp/goinapp goinapp

EXPOSE 8000
ENTRYPOINT ["./goinapp"]