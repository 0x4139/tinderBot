FROM alpine:3.5

COPY ./tinder /dist/tinder

RUN apk --update upgrade && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /dist/
ENTRYPOINT ["/dist/tinder"]