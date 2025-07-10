FROM alpine:latest AS base

FROM base AS builder

RUN apk add --no-cache curl tar && \
    curl -sSf https://atlasgo.sh | sh

FROM base AS prod

COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas

WORKDIR /app

COPY ent/migrate/migrations ./migrations
COPY ent/migrate.entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
