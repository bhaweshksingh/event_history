# Event Historization service


## Pre requisites

- Docker
- Golang v1.14+
 
 
## Running App 

1. Export DB_PASSWORD
`export DB_PASSWORD="S3cretP@ssw0rd"` 

2. Bring up the mysql container using:

`make infra-local`

3. Run the migrations on the local DB 
 
`make setup`

4. Build and run the app container.

`make app`

5. Inspect logs using docker 

`docker logs event-history-go -f`

## Verifying the Functionality

Create key Request
```shell script

curl -X POST 'http://localhost:8080' \
--header 'Content-Type: application/json' \
--data-raw '{
    "key": "name",
    "user_id": "user1",
    "value": "john"
}'
```

Update key Request
```shell script

curl -X PUT 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "key": "name",
    "user_id": "user1",
    "value": "Sam"
}'
```

GET historised answer Req
```shell script
curl -X GET 'http://localhost:8080/user1/name'
```


GET latest answer Req
```shell script
curl -X GET 'http://localhost:8080/latest/user1/name'
```


Delete key Req
```shell script
curl -X DELETE 'http://localhost:8080/user1/name'
```


