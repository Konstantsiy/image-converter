FROM postgres:latest

COPY ./scripts/create_db.sql /docker-entrypoint-initdb.d/