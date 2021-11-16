# build app
FROM golang:alpine AS gobuild
RUN apk update && \
    apk add ca-certificates gcc git make musl-dev
RUN mkdir /home/build
COPY . /home/build
WORKDIR /home/build
RUN go mod tidy && go build -o devicemanager

FROM alpine:latest
ARG PUID=2000
ARG PGID=2000

# Run App
RUN addgroup -g ${PGID} devicemanager && \
    adduser -H -D -u ${PUID} -G devicemanager devicemanager

WORKDIR /opt/devicemanager

COPY --from=gobuild /home/build/devicemanager devicemanager
RUN chown -R ${PUID}:${PGID} /opt/devicemanager

USER devicemanager

CMD ["/opt/devicemanager/devicemanager"]