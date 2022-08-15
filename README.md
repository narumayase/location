# location
REST API to know, based on a point (latitude, longitude), the closest branch office to it.

## Built on 🛠

* [Golang](https://golang.org/)
* [SQLite](https://www.sqlite.org/)
* [Docker](https://www.docker.com/)
* [Swagger](https://swagger.io/)

## Starting 🚀

### Requirements

Necessary tools for the local execution of the service:

- Go go1.15.6+.
- Docker 19.03.6+.

### Example running locally

* Set environment variable LOCATION_ENVIRONMENT=local. 

```
$ export LOCATION_ENVIRONMENT=local
```

* Run: 

```
$ go build
$ go run main.go
```

Response:

```
[GIN-debug] GET    /branch-offices/branch-office/:id --> location/pkg/routes.(*handler).Get-fm (3 handlers)
[GIN-debug] POST   /branch-offices/branch-office --> location/pkg/routes.(*handler).Create-fm (3 handlers)
[GIN-debug] GET    /branch-offices/nearest   --> location/pkg/routes.(*handler).Nearest-fm (3 handlers)
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
```

* To run tests:

```
$ go test ./...

?   	location	[no test files]
?   	location/api	[no test files]
ok  	location/pkg/cmd	0.008s
?   	location/pkg/cmd/mocks	[no test files]
?   	location/pkg/config	[no test files]
?   	location/pkg/db	[no test files]
?   	location/pkg/db/mocks	[no test files]
?   	location/pkg/db/model	[no test files]
ok  	location/pkg/routes	0.012s
```

* To check coverage:

```
$ go test -coverprofile=coverage.out ./...

?   	location	[no test files]
?   	location/api	[no test files]
ok  	location/pkg/cmd	0.008s	coverage: 90.9% of statements
?   	location/pkg/cmd/mocks	[no test files]
?   	location/pkg/config	[no test files]
?   	location/pkg/db	[no test files]
?   	location/pkg/db/mocks	[no test files]
?   	location/pkg/db/model	[no test files]
ok  	location/pkg/routes	0.011s	coverage: 77.1% of statements
```

* Detailed coverage:

```
$ go tool cover -func=coverage.out

location/pkg/cmd/cmd.go:22:		Build		100.0%
location/pkg/cmd/cmd.go:26:		Nearest		100.0%
location/pkg/cmd/cmd.go:52:		Get		0.0%
location/pkg/cmd/cmd.go:60:		Create		100.0%
location/pkg/cmd/cmd.go:70:		all		100.0%
location/pkg/cmd/cmd.go:74:		Find		100.0%
location/pkg/cmd/translator.go:8:	toModel		100.0%
location/pkg/cmd/translator.go:17:	toApi		100.0%
location/pkg/cmd/translator.go:26:	toApis		100.0%
location/pkg/routes/route.go:22:	AddHandler	100.0%
location/pkg/routes/route.go:33:	Get		0.0%
location/pkg/routes/route.go:49:	Create		70.0%
location/pkg/routes/route.go:65:	create		100.0%
location/pkg/routes/route.go:75:	alreadyExists	100.0%
location/pkg/routes/route.go:82:	Nearest		100.0%
total:					(statements)	82.7%
```

### Example using docker image:

```
$ sudo docker image build -t location:0.0.1 .
$ sudo docker run --network host location:0.0.1
```

### API 

More detailed Swagger documentation in the file [swagger.yaml](https://github.com/narumayase/location/blob/master/swagger.yaml).

This API is used to register branch offices, to get the branch by id and also to find the closest branch office to a given point. This point must be expressed in latitude and longitude in simple decimal coordinates with the following format:

Buenos Aires location example:

```
Latitude:-34.6083
Longitude:-58.3712
```

* Branch office creation

Allows you to create a branch with its latitude, longitude and address.

```
request:

POST http://localhost:8080/branch-offices/branch-office

{
    "latitude":-34.6083,
    "longitude":-58.3712,
    "address":"just a random address 123"
}

response:

200 OK

{
    "id":1,
    "latitude":-34.6083,
    "longitude":-58.3712,
    "address":"just a random address 123"
}
```

Example:

```
$ curl -v -d '{"longitude":-58.45678, "latitude":-34.12345, "address":"buenos aires 1234"}' POST http://localhost:8080/branch-offices/branch-office
```

Response:

```
* Rebuilt URL to: POST/
* Could not resolve host: POST
* Closing connection 0
curl: (6) Could not resolve host: POST
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#1)
> POST /branch-offices/branch-office HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Length: 76
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 76 out of 76 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Wed, 06 Jan 2021 16:28:12 GMT
< Content-Length: 81
< 
* Connection #1 to host localhost left intact
{"id":1,"longitude":-58.45678,"latitude":-34.12345,"address":"buenos aires 1234"}
```

* Get branch office by ID

It allows to obtain the data of the branch knowing its id.

```
request:

GET http://localhost:8080/branch-offices/branch-office/1

response:

200 OK
{
    "latitude":-34.6083,
    "longitude":-58.3712,
    "address":"just a random address 123",
    "id":1
}
```

Example:

```
$ curl -v GET http://localhost:8080/branch-offices/branch-office/1
```

Response:

```
* Rebuilt URL to: GET/
* Could not resolve host: GET
* Closing connection 0
curl: (6) Could not resolve host: GET
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#1)
> GET /branch-offices/branch-office/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Wed, 06 Jan 2021 16:29:25 GMT
< Content-Length: 81
< 
* Connection #1 to host localhost left intact
{"id":1,"longitude":-58.45678,"latitude":-34.12345,"address":"buenos aires 1234"}
```

* Find the nearest branch office

Allows you to search for the closest branch to the chosen point (latitude and longitude).

```
request:

GET http://localhost:8080/branch-offices/nearest?latitude=-34.6083&longitude=-58.3712

response:

200 OK
{
    "latitude":-34.6083,
    "longitude":-58.3712,
    "address":"just a random address 123",
    "id":1
}
```

Example:

```
$ curl -v GET 'http://localhost:8080/branch-offices/nearest?longitude=-58.45678&latitude=-34.12345'
```

Response:

```
* Rebuilt URL to: GET/
* Could not resolve host: GET
* Closing connection 0
curl: (6) Could not resolve host: GET
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#1)
> GET /branch-offices/nearest?longitude=-58.45678&latitude=-34.12345 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Wed, 06 Jan 2021 16:30:31 GMT
< Content-Length: 81
< 
* Connection #1 to host localhost left intact
{"id":1,"longitude":-58.45678,"latitude":-34.12345,"address":"buenos aires 1234"}
```

## Roadmap

- [x] Project Structure
- [x] Think algorithm
- [x] Readme
- [x] Algorithm finished
- [x] Latitude and longitude data types
- [x] Unit tests
- [x] Logs
- [x] Dockerize
- [x] Swagger Documentation
