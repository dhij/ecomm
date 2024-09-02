# Ecomm Microservice

This is an ecomm microservice that consists of ecomm-api, ecomm-grpc, and an ecomm-notification microservices.

This repo is a WIP and consists of a series of tutorials available on Youtube:

- [x] [Ep0 Golang Microservice Course Overview](https://youtu.be/F3AAs5HvBR8)
- [x] [Ep1 Database Setup (Docker & golang-migrate)](https://youtu.be/EFNFtLRcJvw)
- [x] [Ep2 Query MySQL & sqlmock to Write Unit Tests](https://youtu.be/8Sp1lTXtzrQ)
- [x] [Ep3 Query Multiple Tables with One-To-Many in a Database Transaction with sqlx](https://youtu.be/qub5-VSplRI)
- [x] [Ep4 HTTP RESTful API Routing with go-chi](https://youtu.be/v0E6JkBry7I)
- [x] [Ep5 JWT Authentication and Refresh Token](https://youtu.be/HtsEaKuYY2o)
- [x] [Ep6 Admin & Authorization Middlewares](https://youtu.be/ygwF2gtjv18)
- [x] [Ep7 Protobuf, gRPC Client & Server Setup](https://youtu.be/D1a7ny_imUw)
- [x] [Ep8 Enqueue Notifications into a Stateful Database Queue](https://youtu.be/2wiLDfktPzA)
- [x] [Ep9 Send Email Notifications in Goroutines](https://www.youtube.com/watch?v=_Qe3YiZMd7Y)
- [ ] Ep10 Golang Microservice Local Development Setup - TBD
- [ ] Ep11 Kubernetes Basics - TBD

## Ecomm Architecture

![Ecomm Architecture](/assets/ecomm_architecture.jpg)

## Run the gRPC and API Microservices

### Database Setup
```
docker pull mysql:8.4
docker run --name ecomm-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4
```
Pull and run the mysql image in a detached mode (`-d`) with port-forwarding and an environment variable.
```
docker exec -i ecomm-mysql mysql -uroot -ppassword <<< "CREATE DATABASE ecomm;"
```
Make the `CREATE DATABASE` query in the mysql container with `docker exec`.
```
docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3306)/ecomm" up
```
Apply the migration files in `db/migrations` by running the [golang-migrate](https://github.com/golang-migrate/migrate) image with the `up` command.

### Run the Go Apps
```
go run cmd/ecomm-grpc/main.go
go run cmd/ecomm-api/main.go
```

## How Notification Queue Works

![Notification Queue 1](/assets/ecomm-notification-1.jpg)

When an admin creates an order, a new notification event for the "pending" status of an order is enqueued into the `notification_events_queue` table. A notification event is also enqueued for every update of the order status, since we want to send an email notification for every update.the

![Notification Queue 2](/assets/ecomm-notification-2.jpg)

We have a separate `notification_states` table to maintain the state of the notification events even after the event has been deleted/dequeued from the `notification_events_queue`. The state of an event represents whether the notification email is "not sent", "sent", or "failed". Please note that the state of the notification event is what makes this ecomm notification database queue stateful and the "status" of an order has nothing to do with the statefulness.

![Notification Queue 3](/assets/ecomm-notification-3.jpg)

The ecomm-notification microservice will list these notification events from the database queue (ordered in the `created_at` timestamp of the events) and attempt to send an email notification for each of them. An attempt to send a notification email may succeed or fail. We have three possible scenarios for this as a failure case can split up into two different cases.

1. Upon success, we delete/dequeue the event from the notifications queue and update the record in the notification states table to `sent`.
2. Upon failure:
    - If the number of attempts is less than the maximum number of attempts we have defined (ex. 3), we update the attempt count of the notification event in the database queue. The notification event will remain in the queue and processed in the following round.
    - If the number of attempts has exceeded the maximum number of attempts allowed, we delete/dequeue the event from the notifications queue and update the record in the notification states table to `failed`.