FROM docker.io/golang:alpine as build

WORKDIR /app/

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY account/account.go account/account.go
ADD main.go .
RUN go build -o /gi-helper

FROM docker.io/alpine

RUN apk add --no-cache curl
COPY --from=build /gi-helper /bin/gi-helper

CMD [ "/bin/gi-helper" ]