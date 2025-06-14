FROM postgres:latest

ENV POSTGRES_DB my_database

WORKDIR /docker-entrypoint-initdb.d

ADD psql_dump.sql .

EXPOSE 5432