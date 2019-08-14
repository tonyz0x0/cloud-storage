# cloud-storage
A cloud storage based on Go

## Architecture

The cloud storage system is organized in microservice architecture. It is combined with multi parts, including Transfer, Upload, Download, Account, API Gateway and Database Proxy. Services are communicated through grpc provided by **go-micro**.

## Dependencies

You should install all dependencies in proper GOPATH before running the program.

### Database

Make sure that **MySQL** and **Redis** are installed in local machine.

### Message Queue

The system requires message queue middleware **RabbitMQ** for asynchronous transmission.

### OSS

Files are stored in both local machine and **Alicloud OSS**

## How to run

```shell
$ cd $GOPATH/filestore-server

$ ./service/start-all.sh
```