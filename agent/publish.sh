berlioz build --nocache
# docker login
docker tag berlioz-main-agent berliozcloud/agent
docker push berliozcloud/agent
