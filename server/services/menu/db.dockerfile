FROM postgres:17-alpine
COPY db/up.sql /docker-entrypoint-initdb.d/1.sql
CMD [ "postgres" ]