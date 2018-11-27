# LinkedinApiServer
This is a repo for demo HTTP API server.

### Running
```
go build
```

Then running the below command will start the server with default configuration

```
./LinkedinApiServer start
```
To give a port number
```
./LinkedinApiServer start -p <port number>
```
To bypass login information
```
./LinkedinApiServer start -b
```
To stop the server after a definite time
``` 
./LinkedinApiServer start -s <minute>
```
To check the version of the API
```
./LinkedinApiServer version
```
To get a JWT token
```
./LinkedinApiServer gentkn -u <username> -e <expiration time in minute>
```

### CURL commands

##### Read all profiles
```
curl --user admin:admin -s -X GET localhost:8080/in
```
##### Read a profile
```
curl --user admin:admin -s -X GET localhost:8080/in/fahim-abrar
```
##### Delete a profile
```
curl -H "Authorization: Bearer <token>" -s -X DELETE localhost:8080/in/masud-rahman
```
##### Create a profile
```
curl -H "Authorization: Bearer <token>" -s -X POST -H 'Content-Type: application/json' -d '{"id":"kfoozminus","name":"Jannatul Ferdous","company":"AppsCode Inc.","position":"Software Engineer","skill":[{"name":"C++","noOfEndorsement":100},{"name":"C","noOfEndorsement":100}]}' localhost:8080/in
```
##### Update a profile info
```
curl -H "Authorization: Bearer <token>" -s -X PUT -H 'Content-Type: application/json' -d '{"id":"masud-rahman","name":"Masudur Rahman (modified)","company":"AppsCode Inc.","position":"Software Engineer","skill":[{"name":"C","noOfEndorsement":3},{"name":"C++","noOfEndorsement":4}]}' localhost:8080/in/masud-rahman
```
##### Invalid Authorization Header check
```
curl --user fahim:dfsd:d localhost:8080/in
```
##### Expired Token Check
```
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTQzMjU5OTQ2LCJ1c2VyIjoiZmFoaW0ifQ.qRTYLq4en4MMRZdNs3XjhOAOHSrkt_UqZM-xmpnoXIo" -s -X POST -H 'Content-Type: application/json' -d '{"id":"kfoozminus","name":"Jannatul Ferdous","company":"AppsCode Inc.","position":"Software Engineer","skill":[{"name":"C++","noOfEndorsement":100},{"name":"C","noOfEndorsement":100}]}' localhost:8080/in
```
