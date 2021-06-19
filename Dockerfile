FROM postgres:latest

COPY ./scripts/docker_db.sql /docker-entrypoint-initdb.d/