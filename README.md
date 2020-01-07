# apiDomainInfo

Api-Rest desarrollada en golang, diseñada para consumir un api sobre consultas de informacion de servidores asociados a un dominio, procesar datos y guardarlos en una base de datos implementada con el motor cockroachdb.

## Build Setup

``` bash
# run with environment variable port (listen port for api)
PORT=500 go run main.go

```
## Acceder a la api
>La aplicacion puede ser usada de manera local, corriendo en la terminal los comandos mencionados anteriormente.

>Sin embargo la app se desplego con los servcios hosting de heroku, se implemento en un contenedor, junto con las base de datos, es decir, las base de datos de cockroachdb fue embebida dentro del contenedor junto a la api-rest.

**API Endpoint:** `getServerInfo`

>Retorna la informacion asociada a el domino que se especifique en el parametro en un objeto json. La aplicacion retorna la informacion que se obtiene en un solo llamado, puede no estar completa, si se hacen llamados posteriormente, la aplicacion ira refrescando la informacion

Parametros:

* **ipAddres** - hostname; requerido.

Ejemplos:

* https://apidomaininfo.herokuapp.com/v1/servers/getServerInfo/truora.com
* https://apidomaininfo.herokuapp.com/v1/servers/getServerInfo/google.com


**API Endpoint:** `getServersRecord`

>Retorna las servidores que han sido consultados al momento en la aplicacion

Parametros:

> none

Ejemplo:

* https://apidomaininfo.herokuapp.com/v1/servers/getServersRecord
