FROM berliozcloud/golang-1.10 as builder

ADD src/ /go/src/app
RUN dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .
RUN ls -la /app/main




FROM grafana/grafana

COPY --from=builder /app/main /bin/berlioz_configurator
RUN ls -la /bin/berlioz_configurator

ADD ./provisioning /etc/grafana/provisioning
ADD ./config.ini /etc/grafana/config.ini
RUN mkdir -p /var/lib/grafana/dashboards

USER root
RUN chown -R grafana /etc/grafana/provisioning

USER grafana
ENTRYPOINT ["/bin/berlioz_configurator"]