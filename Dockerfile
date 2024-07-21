FROM golang:alpine AS builder

RUN mkdir /user && \
  echo 'user:x:504:504:user:/home/user:' > /user/passwd && \
  echo 'user:x:504:user' > /user/group

WORKDIR /go/src

COPY . .

RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates
RUN go get -d -v ./...
RUN test -f config.json || cp config.json.example config.json
RUN test -f dictionary-v2.txt || wget https://tirea.learnnavi.org/dictionarydata/dictionary-v2.txt
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o /fwew-api .

FROM scratch

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /fwew-api /fwew-api
COPY --from=builder /go/src/config.json /config.json
COPY --from=builder /go/src/dictionary-v2.txt /home/user/.fwew/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER user:user

ENTRYPOINT [ "/fwew-api" ]
