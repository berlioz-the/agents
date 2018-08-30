FROM berliozcloud/golang-1.10 as builder

ADD src/ /go/src/app
RUN dep ensure -v

# RUN go build -o /app/main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .
RUN ls -la /app/main

FROM prom/prometheus

COPY --from=builder /app/main /bin/berlioz_configurator
RUN ls -la /bin/berlioz_configurator

ADD berlioz_runner.sh /bin/berlioz_runner.sh
RUN ls -la /bin/berlioz_runner.sh

ENTRYPOINT [ "/bin/sh"]
CMD [ "/bin/berlioz_runner.sh"]
# RUN /bin/berlioz_configurator
