## Running


### CONFIG
cluster_announce_ip
cluster_announce_port
cluster_announce_bus_port

### Ports

Cluent: 6379
Gossip: 16379

### Custom Image

docker build -t berlioz-docker .

docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/internals/redis.conf:/etc/redis/redis.conf --name redis-1 redis redis-server /etc/redis/redis.conf

docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/internals/redis.conf:/etc/redis/redis.conf --name redis-2 redis redis-server /etc/redis/redis.conf

### Configuring Cluster
https://get-reddie.com/blog/redis4-cluster-docker-compose/
https://www.ctolib.com/docs/sfile/redis-doc-cn/cn/topics/cluster-spec.html

ruby redis-trib.rb create --replicas 0  172.17.0.2:6379 172.17.0.3:6379 172.17.0.4:6379

### Kill All

docker kill redis-1; docker rm redis-1
docker kill redis-2; docker rm redis-2
docker kill redis-3; docker rm redis-3
docker kill redis-4; docker rm redis-4


### Redis 4
docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-1 redis redis-server /etc/redis/redis.conf

docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-2 redis redis-server /etc/redis/redis.conf

docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-3 redis redis-server /etc/redis/redis.conf


172.17.0.2
172.17.0.3


redis-cli cluster meet 172.17.0.2 6379



### Redis 5

172.17.0.2
172.17.0.3
172.17.0.4
172.17.0.5

docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-1 redis:5.0-rc redis-server /etc/redis/redis.conf
docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-2 redis:5.0-rc redis-server /etc/redis/redis.conf
docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-3 redis:5.0-rc redis-server /etc/redis/redis.conf
docker run -d -v d:/Repos/berlioz-corp/samples.git/08.Redis/docs/redis.conf:/etc/redis/redis.conf --name redis-4 redis:5.0-rc redis-server /etc/redis/redis.conf


redis-cli cluster create --cluster-replicas 2 172.17.0.2:6379 172.17.0.3:6379 172.17.0.4:6379 172.17.0.5:6379



### Debugging

ruby /var/local/redis/redis-trib.rb create --verbose --replicas 0 10.1.1.11:10050 10.1.1.11:10051 10.1.0.35:10050 

ruby /var/local/redis/redis-trib.rb create --verbose --replicas 0  52.91.138.97:10050 52.91.138.97:10051 54.197.133.192:10050

docker-entrypoint.sh redis-server /etc/redis.conf

docker-entrypoint.sh redis-server /etc/redis.conf --cluster-announce-ip 10.1.1.11 --cluster-announce-port 10050 --cluster-announce-bus-port 10060 