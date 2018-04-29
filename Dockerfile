FROM ubuntu:16.04

USER root


ENV PGVER 9.5
RUN apt-get update -q
RUN apt-get install -q -y wget
RUN apt-get install -q -y postgresql-$PGVER


USER postgres

RUN /etc/init.d/postgresql start &&\
    psql -a -f scheme.sql &&\
    /etc/init.d/postgresql stop


RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root

RUN wget "https://dl.google.com/go/go1.10.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz &&\
mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

ENV GOPATH $HOME/go/db_forum
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH &&\
export PATH=$PATH:/usr/local/go/bin

WORKDIR $GOPATH/src
ADD . $GOPATH/src/
EXPOSE 5000
WORKDIR $GOPATH/src/goExample

CMD service postgresql start && go build github.com/pdmitrya/goExample && ./goExample