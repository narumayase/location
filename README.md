# location
API REST para saber en base a un punto (lat, long), la sucursal más cercana al mismo.

## Construido con 🛠

* [Golang](https://golang.org/) - Lenguaje de programación.
* [SQLite](https://www.sqlite.org/) - Base de datos.
* [Docker](https://www.docker.com/) - Manejador de contenedores.
* [Swagger](https://swagger.io/) - Documentación de la API.

## Comenzando 🚀

### Pre-requisitos 

Herramientas necesarias para la ejecución local del servicio:

- Go go1.15.6+.
- Docker 19.03.6+.

### Ejemplo ejecutando localmente

* Configurar variable de entorno LOCATION_ENVIRONMENT=local. 

```
$ export LOCATION_ENVIRONMENT=local
```

* Ejecutar: 

```
$ go build
$ go run main.go
```

Respuesta:

```
[GIN-debug] GET    /branch-offices/branch-office/:id --> location/pkg/routes.(*handler).Get-fm (3 handlers)
[GIN-debug] POST   /branch-offices/branch-office --> location/pkg/routes.(*handler).Create-fm (3 handlers)
[GIN-debug] GET    /branch-offices/nearest   --> location/pkg/routes.(*handler).Nearest-fm (3 handlers)
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
```

* Para ejecutar los tests:

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

* Para evaluar coverage:

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

* Coverage de manera más detallada:

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

### Ejemplo usando imagen docker

```
$ sudo docker image build -t location:0.0.1 .
$ sudo docker run --network host location:0.0.1
```

### API 

Documentación Swagger más detallada en el archivo [swagger.yaml](https://github.com/justskythings/location/blob/master/swagger.yaml).

Esta API sirve para dar de alta sucursales, para consultar la sucursal por id y también para encontrar la sucursal más cercana a un punto dado. Este punto debe estar expresado en latitud y longitud en coordenadas decimales simples con el siguiente formato:

Ejemplo de ubicación de Buenos Aires:

```
Latitud:-34.6083
Longitud:-58.3712
```

* Creación de sucursal

Permite crear una sucursal con su latitud, longitud y dirección.

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

Ejemplo:

```
$ curl -v -d '{"longitude":-58.45678, "latitude":-34.12345, "address":"buenos aires 1234"}' POST http://localhost:8080/branch-offices/branch-office
```

Respuesta:

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

* Obtener sucursal por ID

Permite obtener los datos de la sucursal conociendo su id.

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

Ejemplo:

```
$ curl -v GET http://localhost:8080/branch-offices/branch-office/1
```

Respuesta:

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

* Buscar la sucursal más cercana

Permite buscar la sucursal más cercana al punto (latitud y longitud) elegido.

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

Ejemplo:

```
$ curl -v GET 'http://localhost:8080/branch-offices/nearest?longitude=-58.45678&latitude=-34.12345'
```

Respuesta:

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

- [x] Estructura de proyecto
- [x] Pensar algoritmo 
- [x] Readme
- [x] Algoritmo terminado 
- [x] Tipos de datos de latitud y longitud
- [x] Test unitarios 
- [x] Logs
- [x] Dockerizar
- [x] Documentación Swagger
 
