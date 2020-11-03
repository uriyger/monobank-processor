FROM alpine:3.12


# Add the user UID:1000, GID:1000
RUN mkdir -p /app && \
    chown -R nobody:nobody /app

COPY app_linux /app/app
RUN apk add libcap && \
    apk --no-cache add tzdata && \
    setcap 'cap_net_bind_service=+ep' /app/app

EXPOSE ${HTTP_PORT}
EXPOSE 8080

WORKDIR /app
USER nobody

ENTRYPOINT ["/app/app"]