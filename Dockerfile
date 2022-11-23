FROM alpine:3.17

LABEL org.opencontainers.image.source="https://github.com/powerman/go-service-example"

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=5s \
    CMD wget -q -O - http://$HOSTNAME:8000/health-check || exit 1

COPY . .

ENTRYPOINT [ "bin/address-book" ]

CMD [ "serve" ]
