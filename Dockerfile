FROM golang:1.13-alpine AS build 
ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN apk add git make openssl

WORKDIR /go/src/apps-schedule-webhook
ADD . .
RUN mkdir ssl

RUN make build

FROM scratch
WORKDIR /app
COPY --from=build /go/src/apps-schedule-webhook/schedulemutator .
COPY --from=build /go/src/apps-schedule-webhook/ssl ssl
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app/schedulemutator"]
