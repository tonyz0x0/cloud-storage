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