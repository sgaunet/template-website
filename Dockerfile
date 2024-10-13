FROM alpine:3.18.4 AS builder
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
# COPY resources /
COPY template-website               /opt/webserver/
# COPY static                             /opt/webserver/static/
# COPY templates                          /opt/webserver/templates
WORKDIR /opt/webserver
USER MyUser
EXPOSE 8080
CMD [ "/opt/webserver/template-website" ]
