# Description

This project is designed to manage and process data related to pick-up points (ПВЗ) and orders. It provides a back-end
system that manages information creation, retrieval, updating, and deletion, as well as management functions.

## Technologies used:

- Apache Kafka: Employs this messaging system for message handling.
- PostgreSQL: Acts as the database to store pick-up point and order details.
- Docker: Containers are used to encapsulate the application and its dependencies, ensuring consistency across
  environments.
- Docker Compose: Manages multi-container Docker applications, simplifying configurations and deployment of services
  like Kafka and PostgreSQL.

## TESTS:

### `make test-module`

> Runs unit testing

### `make test-integration`

> Runs integration testing

### ! Before running tests !

Make sure that both DB are created (including test DB):

`make migration-up`

`make test-migration-up`

## COMMANDS:

### http

> Launches http server

Implemented methods for `/pick-up-point`: `POST`, `GET`

Implemented methods for `/pick-up-point/<ID>`: `GET`,`DELETE`, `PUT`

Examples of using: [CURL examples](#curl-examples)

### cr-take

> Accepts and writes an order from the courier into a file

Required flags: `-cid`, `-oid`, `-at`, `-w`, `-pr`, `-pg`

`-pg` accepts: `package`, `carton`, `film`

**Cannot accept order twice or with a negative available time**

Example of using: `cr-take -cid=1 -oid=1 -at=48h -pr=100 -w=0.5 -pg=carton`

### cr-return

> Returns an order to the courier (deletes the order from file)

Required flag: `-oid`

**Can return orders only with an expiration date in the past and those that have not been given**

Example of using: `cr-return -oid=1`

### cl-give

> Gives one or more orders to the client

Required flags: `-cid`, `-oids`

**Only orders accepted from the courier (present in the file) with expiration dates earlier than the current date can be
given. All order IDs must belong to the same client.**

Example of using: `cl-give -cid=1 -oids=1,3,7`

### cl-orders

> Get orders

Required flag: `-cid`
Optional flags: `-n`, `-ouo`

Example of using: `cl-orders -cid=1 [-n=5] [-ouo]`

### cl-refund

> Accepts a return from the client (move an order from available orders to refunded orders)

Required flags: `-cid`, `-oid`

**Orders can be returned within 2 days from the given date.**

Example of using: `cl-refund -cid=1 -oid=1`

### refund-list

> Provides the list of refunds in a paginated manner

Optional flag: `-p`

Example of using: `refund-list [-p=3]`

### interactive

> Launches an interactive mode that has two commands: `write` and `read`

Commands:

`read` - shows all pick-up points

`write Name,Address,Contact Info` - takes 3 comma-separated arguments (Name, Address & Contact Information) of pick-up
point and writes to file

Example of using: `write Pick-up point #1, Tomorrow Avenue, +78005553535`

## FLAGS

````
-cid int
    Client ID

-at duration
    Available time for picking up the order (e.g. -at=48h)
    
-n int
    How many last orders the user will receive (default -1)

-oid int
    Order ID

-oids string
    Comma-separated slice of ordersID (e.g. -oids=1,3,7)

-ouo
    Get only user orders

-p int
    Page number starting with 1 (default 1)`
    
-pr float
    Price
    
-pg string
    Packaging

-w float
    Weight
````

## CURL examples

### Failed authentication

```
> curl -X DELETE http://localhost:9000/pick-up-point/1 \
-u maksim:makarov
Authentication failed
```

### `/pick-up-point` POST method

```
> curl -o - -u ildus:erbaev -X POST http://localhost:9000/pick-up-point \
-d '{"name":"Pick-up Point A","address":"123 Main St","contact":"8-800-555-35-35"}'

{"ID":1,"name":"Pick-up Point A","address":"123 Main St","contact":"8-800-555-35-35"}
```

### `/pick-up-point` GET method

```
> curl -o - -u ildus:erbaev -X GET http://localhost:9000/pick-up-point 
[{"ID":1,"Name":"Pick-up Point A","Address":"123 Main St","Contact":"123-456-7890"},{"ID":2,"Name":"Pick-up Point B","Address":"123 Main St","Contact":"123-456-7890"},{"ID":3,"Name":"Pick-up Point C","Address":"123 Main St","Contact":"123-456-7890"},{"ID":4,"Name":"Pick-up Point D","Address":"123 Main St","Contact":"8-800-555-35-35"}]
```

### `/pick-up-point/<ID>` PUT method

```
> curl -o - -u ildus:erbaev -X PUT http://localhost:9000/pick-up-point/1 \
-d '{"name":"Pick-up Point ABC","address":"123 Main St","contact":"123-456-7890"}'

{"ID":1,"name":"Pick-up Point ABC","address":"123 Main St","contact":"123-456-7890"}
```

```
> curl -o - -u ildus:erbaev -X PUT http://localhost:9000/pick-up-point/1 \
-d '{"name":NOTSTRING,"address":"123 Main St","contact":"123-456-7890"}' 
error occured: invalid character 'N' looking for beginning of value
```

### `/pick-up-point/<ID>` DELETE method

```
> curl -o - -u ildus:erbaev -X DELETE http://localhost:9000/pick-up-point/2
```

```
curl -X DELETE http://localhost:9000/pick-up-point/9999 \
-u ildus:erbaev  
error occured: not found
```

## Test cases

### cr-take

Successful accept:

````
> cr-take -oid=1 -cid=1 -at=2h -pr=100 -w=6 -pg=carton
Courier's order accepted successfully!
````

Negative available time:

````
> cr-take -oid=12 -cid=46 -at=-48h -pr=100 -w=6 -pg=carton
available time is not given or incorrect
````

Duplicate order ID:

````
> cr-take -oid=12 -cid=46 -at=48h -pr=100 -w=6 -pg=carton
Courier's order accepted successfully!
> cr-take -oid=12 -cid=46 -at=48h -pr=100 -w=6 -pg=carton
order has been already accepted
````

Order weight exceeds limit for packaging:

```
> cr-take -oid=1 -cid=3 -at=-2h -pr=100 -w=11 -pg=carton
order weight exceeds limit for carton
```

Invalid packaging type:

```
> cr-take -oid=1 -cid=3 -at=-2h -pr=100 -w=11 -pg=nothing
invalid packaging type
```

### cr-return

Successful return:

````
> cr-take -oid=100 -cid=1 -at=0.001s
Courier's order accepted successfully!
> cr-return -oid=100
The order was returned to the courier successfully!
````

Order not found:

````
> cr-return -oid=1234
the order was not found
````

Expiration date:

````
> cr-take -oid=40 -cid=1 -at=48h
Courier's order accepted successfully!
> cr-return -oid=40
the order was given or the expiration date is not over yet
````

Order already given:

````
> cr-take -oid=40 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oids=40 -cid=1
Order was given to the client successfully!
> cr-return -oid=40
the order was given or the expiration date is not over yet
````

### cl-give

Successful return:

````
> cr-take -oid=100 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oid=100
Order was given to the client successfully!
````

````
> cr-take -oid=101 -cid=1 -at=48h
Courier's order accepted successfully!
> cr-take -oid=102 -cid=1 -at=48h
Courier's order accepted successfully!
> cr-take -oid=103 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oids=101,102,103 -cid=1
Orders were given to the client successfully!
````

Order expired:

````
> cr-take -oid=100 -cid=1 -at=0.001s
Courier's order accepted successfully!
> cl-give -oids=100 -cid=1
order expired
````

Order has been already accepted:

````
> cr-take -oid=100 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oid=100
Order was given to the client successfully!
> cr-take -oid=100 -cid=1 -at=48h
order has been already accepted
````

````
> cr-take -oid=101 -cid=1 -at=48h
Courier's order accepted successfully!
> cr-take -oid=102 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oid=102
Order was given to the client successfully!
> cl-give -oid=101,102
order has been already accepted
````

### cl-orders

Successful execution:

````
> cl-orders -cid=1
{ID:5 ClientID:3 ExpirationDate:2024-03-29 02:09:27.434786 +0400 +04 Weight:9 Price:105 Packaging:package}
{ID:4 ClientID:3 ExpirationDate:2024-03-29 02:05:14.222981 +0400 +04 Weight:11 Price:120 Packaging:carton}
{ID:2 ClientID:1 ExpirationDate:2024-03-29 01:48:44.650589 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:1 ClientID:1 ExpirationDate:2024-03-29 01:47:06.383251 +0400 +04 Weight:0.5 Price:101 Packaging:film}
````

````
> cl-orders -cid=1 -ouo
{ID:2 ClientID:1 ExpirationDate:2024-03-29 01:48:44.650589 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:1 ClientID:1 ExpirationDate:2024-03-29 01:47:06.383251 +0400 +04 Weight:0.5 Price:101 Packaging:film}
````

````
> cl-orders -cid=1 -ouo -n=1
{ID:2 ClientID:1 ExpirationDate:2024-03-29 01:48:44.650589 +0400 +04 Weight:6 Price:120 Packaging:carton}
````

Empty list:

````
cl-orders -cid=137 -ouo
List of orders is empty
````

### cl-refund

Successful refund:

````
> cr-take -oid=200 -cid=1 -at=48h
Courier's order accepted successfully!
> cl-give -oids=200 -cid=1
Order was given to the client successfully!
> cl-refund -oid=200 -cid=1
Order refunded successfully!
````

Refund time exceeded:

```
> cl-refund -oid=201 -cid=1
it has been more than 2 days since it was given or order was not given
```

Order does not belong to client:

```
> cr-take -oid=202 -cid=2 -at=48h
Courier's order accepted successfully!
> cl-give -oids=202 -cid=2
Order was given to the client successfully!
> cl-refund -oid=202 -cid=1
order does not belong to the client
```

### refund-list

Successful execution:

```
> refund-list
                        Page number: 1
{ID:1 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:2 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:3 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:4 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:5 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:6 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:7 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:8 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:9 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:10 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
                        Page number: 1
```

```
> refund-list -p=2
                        Page number: 2
{ID:11 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:12 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:13 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:14 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:15 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:16 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:17 ClientID:2 ExpirationDate:2024-04-03 18:28:27.993424 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:18 ClientID:223 ExpirationDate:2024-03-05 22:23:53.936445 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:16 ClientID:223 ExpirationDate:2024-03-05 22:23:53.936445 +0400 +04 Weight:6 Price:120 Packaging:carton}
{ID:200 ClientID:1 ExpirationDate:2024-03-07 11:27:42.434234 +0400 +04 Weight:6 Price:120 Packaging:carton}
                        Page number: 2

```

Page does not exist:

```
> refund-list -p=5
page does not exists
```

### interactive

Successful execution:

```
> go run main.go interactive
> read
{ID:1 Name:Pick-up point #1 Address:Lenin street Contact:+78005553535}
{ID:2 Name:Pick-up point #2 Address:Stalin street Contact:+78005554536}
{ID:3 Name:Pick-up point #3 Address:Gorbachev street Contact:+78505253535}
```

```
> go run main.go interactive
> write Pick-up point #4, Pobeda street, +74004654353
Pick-up point has been written successfully!
```

Empty list:

```
> read
List of points is empty
```

No arguments:

```
> write
write command requires arguments
```

Less or more than 3 arguments:

```
> write Pick-up point #5
expected 3 arguments for write command: Name, Address, Contact info
```

```
> write Pick-up point #6, Pushkin steet, +7450044533, Marina
expected 3 arguments for write command: Name, Address, Contact info
```