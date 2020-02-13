FROM golang:alpine AS intermediate

RUN apk update && \
    apk add --no-cache git make

RUN adduser -D -g '' schooner

WORKDIR $GOPATH/src/

COPY . .

RUN go mod download
RUN go mod verify
RUN make build
RUN make build-healthchecker

FROM scratch

ENV PORT=8000

COPY --from=intermediate /go/src/bin/schooner /go/bin/schooner
COPY --from=intermediate /go/src/bin/healthchecker /go/bin/healthchecker
COPY --from=intermediate /etc/passwd /etc/passwd

USER schooner

WORKDIR /go/bin

HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD ["/go/bin/healthchecker"]

CMD ["/go/bin/schooner"]