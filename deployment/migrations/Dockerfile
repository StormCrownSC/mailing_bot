FROM migrate/migrate:latest

COPY ./migrations /migrations
ENTRYPOINT ["/bin/sh", "-c", "sleep 10 && migrate -path ./migrations -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}\" up"]