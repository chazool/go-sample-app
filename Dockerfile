FROM golang:1.20.1-bullseye as build 


ENV GOLANG_VERSION="1.20"

RUN apt-get clean
RUN apt-get update && apt-get install -y ca-certificates openssl git tzdata
ARG cert_location=/usr/local/share/ca-certificates

RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.ctr
RUN openssl s_client -showcerts -connect gitlab.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gitlab.ctr
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/proxy.golang.ctr
RUN openssl s_client -showcerts -connect gopkg.in:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gopkg.ctr
RUN openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/storage.googleapis.ctr
RUN openssl s_client -showcerts -connect sum.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/sum.golang.ctr

RUN update-ca-certificates
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/russross/blackfriday-tool@latest

COPY /common  /common

WORKDIR /common
RUN make build

FROM gcr.io/distroless/static-debian11
COPY --from=build /common/build /

ARG DD_API_KEY
ARG DD_SITE
ARG DD_LOG_ENABLE
ARG DD_SERVICE
ARG DD_ENV
ARG DD_VERSION
ARG TRACER_SINK

ENV DD_API_KEY=${DD_API_KEY}
ENV DD_SITE=${DD_SITE}
ENV DD_LOG_ENABLE=${DD_LOG_ENABLE}
ENV DD_SERVICE=${DD_SERVICE}
ENV DD_ENV=${DD_ENV}
ENV DD_VERSION=${DD_VERSION}
ENV TRACER_SINK=${TRACER_SINK}

COPY --from=datadog/serverless-init /datadog-init:beta3 /app/datadog-init
ENTRYPOINT [ "/app/datadog-init" ]

CMD ["/commonservice"]

