FROM golang

COPY [".", "/go/src/apiDomainInfo"]

RUN wget -qO- https://binaries.cockroachdb.com/cockroach-v19.2.2.linux-amd64.tgz | tar  xvz

RUN cp -i cockroach-v19.2.2.linux-amd64/cockroach /usr/local/bin/

WORKDIR /go/src/apiDomainInfo

EXPOSE 3000

CMD ["sh", "up.sh"]