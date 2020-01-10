FROM node:8 as front
WORKDIR /app
COPY . .
WORKDIR /app/web
RUN yarn install
RUN yarn run build

FROM scratch
WORKDIR /root/
COPY --from=front /app .
COPY ./build/amd64/drawio .

CMD ["./drawio"]
EXPOSE 80