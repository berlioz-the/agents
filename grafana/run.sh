docker stop grafana-test
docker rm grafana-test

berlioz build
docker tag berlioz-main-grafana berliozcloud/grafana
docker run --name grafana-test -p 9999:3000 -d \
    --network="berlioz" \
    -e "BERLIOZ_TASK_ID=4c03085a-7986-46b7-a9e3-903a88d0b01f" \
    -e "BERLIOZ_ADDRESS=172.16.0.3" \
    -e "BERLIOZ_IDENTITY=1" \
    -e "BERLIOZ_LISTEN_PORT_DEFAULT=3000" \
    -e "BERLIOZ_PROVIDED_PORT_DEFAULT=3000" \
    -e "BERLIOZ_AGENT_PATH=ws://agent1.main.berlioz:55555/4c03085a-7986-46b7-a9e3-903a88d0b01f" \
    -e "BERLIOZ_LISTEN_ADDRESS=0.0.0.0" \
    -e "BERLIOZ_INFRA=local" \
    -e "BERLIOZ_REGION=earth-local" \
    -e "BERLIOZ_INSTANCE_ID=local-1234" \
    -e "BERLIOZ_CLUSTER=sprt" \
    -e "BERLIOZ_SECTOR=main" \
    -e "BERLIOZ_SERVICE=grafana" \
    berlioz-main-grafana
