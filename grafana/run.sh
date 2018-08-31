docker stop grafana-test
docker rm grafana-test

berlioz build
docker tag berlioz-main-grafana berliozcloud/grafana
docker run --name grafana-test -p 9999:3000 -d berlioz-main-grafana
