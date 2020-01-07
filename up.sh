#!/bin/bash

cockroach start-single-node \
--insecure \
--listen-addr=localhost:26257 \
--http-addr=localhost:8080 \
--background

sleep 10s

cockroach sql --insecure \
--user=root \
--host=localhost \
--port=26257 \
--database=postgres < /go/src/apiDomainInfo/sql/statements.sql

sleep 10s

go run /go/src/apiDomainInfo/main.go
