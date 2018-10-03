berlioz build --nocache
# docker login
docker tag berlioz-main-redis berliozcloud/redis
docker push berliozcloud/redis
