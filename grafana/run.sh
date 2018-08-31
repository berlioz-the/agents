docker stop grafana-test
docker rm grafana-test

berlioz build
# docker login
docker tag berlioz-main-grafana berliozcloud/grafana
# docker push berliozcloud/grafana

docker run --name grafana-test -p 9999:3000 -d berlioz-main-grafana
