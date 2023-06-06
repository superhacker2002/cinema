FROM postgres
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 2587
ENV POSTGRES_DB cinema
COPY database/cinemadb.sql /docker-entrypoint-initdb.d/
COPY database/mock_data.sql /docker-entrypoint-initdb.d/
