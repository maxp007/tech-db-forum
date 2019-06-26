﻿FROM ubuntu:18.04

ENV PGSQLVER 10
ENV DEBIAN_FRONTEND 'noninteractive'

RUN echo 'Europe/Moscow' > '/etc/timezone'

RUN apt-get -o Acquire::Check-Valid-Until=false update
RUN apt install -y gcc git wget
RUN apt install -y postgresql-$PGSQLVER

RUN wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
RUN tar -xvf go1.12.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /server
COPY . .

RUN cd /server
RUN go get -u
ENV PORT 5000
EXPOSE $PORT

USER postgres

RUN /etc/init.d/postgresql start &&\
	psql --echo-all --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
	createdb -O docker docker &&\
	psql -d docker -f database/dump.sql &&\
	/etc/init.d/postgresql stop

RUN rm -rf /etc/postgresql/$PGSQLVER/main/pg_hba.conf
RUN echo "local   all             postgres                                peer\n\
local   all             docker                                md5\n\
host    all             all             127.0.0.1/32            md5\n\
host all  all    0.0.0.0/0  md5" >>\
    /etc/postgresql/$PGSQLVER/main/pg_hba.conf

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGSQLVER/main/pg_hba.conf &&\
	echo "listen_addresses='*'" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "fsync = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "synchronous_commit = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "shared_buffers = 256MB" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "work_mem = 51MB" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "maintenance_work_mem = 256MB" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "wal_sync_method = fdatasync" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "commit_delay = 55" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "commit_siblings = 8" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "random_page_cost = 2.0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "cpu_tuple_cost = 0.001" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "cpu_index_tuple_cost = 0.0005" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "autovacuum = on" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "wal_level = minimal" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "max_wal_senders = 0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf

EXPOSE 5432

USER root

CMD service postgresql start && go run main.go



USER postgres

RUN /etc/init.d/postgresql start &&\
	psql --echo-all --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
	createdb -O docker docker &&\
	psql -d docker -f main_microservice/database/sql.sql &&\
	/etc/init.d/postgresql stop

