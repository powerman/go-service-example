FROM alpine:3.12

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=5s \
    CMD wget -q -O - http://$HOSTNAME:8000/health-check || exit 1

COPY . .

ENTRYPOINT [ "bin/address-book" ]

CMD [ "serve" ]
