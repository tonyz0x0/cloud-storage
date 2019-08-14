## Create Databases with Docker

Create two database, one master and one slave

```sh
$ docker run --name storage_db_master_01 -e MYSQL_ROOT_PASSWORD=123456 -d -p 3307:3306 mysql/mysql-server:latest

$ docker run --name storage_db_master_01 -e MYSQL_ROOT_PASSWORD=123456 -d -p 3308:3306 mysql/mysql-server:latest

# CONTAINER ID        IMAGE                       COMMAND                  CREATED              STATUS                        PORTS                               NAMES
# 5418cbfab63f        mysql/mysql-server:latest   "/entrypoint.sh mysq…"   About a minute ago   Up About a minute (healthy)   33060/tcp, 0.0.0.0:3308->3306/tcp   storage_db_slave_01
# fd2fb6dd7373        mysql/mysql-server:latest   "/entrypoint.sh mysq…"   2 minutes ago        Up 2 minutes (healthy)        33060/tcp, 0.0.0.0:3307->3306/tcp   storage_db_master_01
```

## Config Master-Slave

Start two terminals and oo into both the master and slave database dockers and start them

```sh

# master
$ docker exec -it storage_db_master_01 bash
$ mysql -u root -p

# slave
$ docker exec -it storage_db_slave_01 bash
$ mysql -u root -p
```

Change user and password in both databases:

```sh
mysql> ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
```

Find **binlog** file in **master**:

```sh
mysql> show master status;
```

Config **slave** with **master** information:

```sh

CHANGE MASTER TO MASTER_HOST='10.0.0.243', MASTER_USER='root', MASTER_PASSWORD='123456', MASTER_LOG_FILE='binlog.000005', MASTER_LOG_POS=0;

show slave status;

start slave;
```

If you want to stop slave database:

```sh
mysql> stop slave io_thread for channel ''
```

Start two database containers:

```sh
$ docker container start <docker ID>
```

## Ceph Config

```sh

# Monitor
docker run -d --net=host -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /var/log/ceph/:/var/log/ceph/ -e MON_IP=172.20.0.1 -e CEPH_PUBLIC_NETWORK=172.20.0.0/24 --name="ceph-mon" ceph/daemon mon

# mgr
docker run -d --net=host --privileged=true --pid=host --name="ceph-mgr" -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ ceph/daemon mgr

# OSD
docker exec ceph-mon ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring

mkdir -p /data/ceph/osd/vdb

docker run -d --privileged=true --name=ceph-osdvdb --net=host -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /data/ceph/osd/vdb:/var/lib/ceph/osd -e OSD_TYPE=directory -v /etc/localtime:/etc/localtime:ro ceph/daemon osd

# Check status
sudo docker exec ceph-mon ceph -s

# Gateway
docker exec ceph-mon ceph auth get client.bootstrap-rgw -o /var/lib/ceph/bootstrap-rgw/ceph.keyring

docker run -d --net=host --privileged=true --name=ceph-rgw -v /var/lib/ceph/:/var/lib/ceph/ -v /etc/ceph:/etc/ceph -v /etc/localtime:/etc/localtime:ro -e RGW_NAME=rgw0 ceph/daemon rgw

```

## OSS

## RabbitMQ

```sh
# Username: guest, Passowrd: guest
$ mkdir -p /data/rabbitmq
$ sudo chmod 755 /data/rabbitmq/

# Remember to config path /data to file sharing in docker
$ docker run -d --hostname rabbit-node1 --name rabbit-node1 -p 5672:5672 -p 15672:15672 -p 25672:25672 -v /data/rabbitmq:/var/lib/rabbitmq rabbitmq:management

# Start RabbitMQ Container
$ docker container start rabbitmq:management
```

# Protocol Buffers Compiler

# Consul

Install Consul with Docker
```sh
$ docker pull consul:1.4.5

$ docker run --name consul1 -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 consul:1.4.5 agent -server -bootstrap -ui -bind=0.0.0.0 -client=0.0.0.0
```