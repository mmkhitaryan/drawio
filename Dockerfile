FROM node:12.20.0-alpine3.11 as front
WORKDIR /home/go/app
COPY ./app .
RUN npm install
RUN npm run build

FROM golang:1.13.4-alpine as builder
RUN apk add --update --no-cache bash tzdata git
RUN go get -u github.com/mmkhitaryan/drawio
WORKDIR $GOPATH/src/github.com/mmkhitaryan/drawio
RUN go install
RUN mkdir -p /home/go/app
RUN mkdir -p /usr/local/go/lib/time
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /home/go/app/drawio main.go

FROM scratch
COPY --from=builder /home/go/app /home/go/app
WORKDIR /home/go/app
COPY --from=builder /home/go/app/drawio /home/go/app/drawio
COPY --from=front /home/go/app/dist /home/go/app/web
EXPOSE 80
ENTRYPOINT ["/home/go/app/drawio"]