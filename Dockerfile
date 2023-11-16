FROM golang:latest AS compiling
RUN mkdir -p /go/src/NewsAgg
WORKDIR /go/src/NewsAgg
ADD . .
WORKDIR /go/src/NewsAgg/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=compiling /go/src/NewsAgg/cmd/app .
COPY --from=compiling /go/src/NewsAgg/cmd/config.json .
#db connection config from ENV
ARG dbhost=192.168.1.35:5432/news
ENV dbhost="${dbhost}"
ARG dbpass=to_be_redefined_at_conrainer_start
ENV dbpass="${dbpass}"

CMD ["./app"]
EXPOSE 8080