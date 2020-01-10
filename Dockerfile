FROM node:8 as front
RUN mkdir -p /home/go/app
COPY ./web /home/go/app
WORKDIR /home/go/app/web
RUN yarn install
RUN yarn run build

FROM scratch
COPY --from=front /home/go/app/web/dist /home/go/app
COPY ./build/linux/drawio /home/go/app/drawio
WORKDIR /home/go/app
EXPOSE 80
ENTRYPOINT ["/home/go/app/drawio"]