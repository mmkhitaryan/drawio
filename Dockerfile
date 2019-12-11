FROM node:8 as front
WORKDIR /app
COPY . .
WORKDIR /app/web
RUN yarn install
RUN yarn run build

FROM golang:latest
WORKDIR /root/
COPY --from=front /app .
RUN go get -d -v ./...
RUN go build -o main .
CMD ["./main"]
EXPOSE 80