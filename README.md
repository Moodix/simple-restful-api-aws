# Simple Restful API with GO on AWS

A simple Restful API on AWS using the following tech stack:

[Serverless Framework](https://serverless.com/)

[Go language](https://golang.org/)

[AWS API Gateway](https://aws.amazon.com/api-gateway/)

[AWS Lambda](https://aws.amazon.com/lambda/)

[AWS DynamoDB](https://aws.amazon.com/dynamodb/)


The API accepts the following JSON requests and produce the corresponding HTTP responses:

### Request 1:
```
HTTP POST
URL: https://<api-gateway-url>/api/devices
Body (application/json):
{
"id": "/devices/id1",
"deviceModel": "/devicemodels/id1",
"name": "Sensor",
"note": "Testing a sensor.",
"serial": "A020000102"
}
```
### Response 1 - Success:
```
HTTP 201 Created
```
### Response 1 - Failure 1:
``` 
HTTP 400 Bad Request
If any of the payload fields are missing. Response body should have a descriptive
error message for the client to be able to detect the problem.
```
### Response 1 - Failure 2:
```
HTTP 500 Internal Server Error
If any exceptional situation occurs on the server side.
```
### Request 2:
```
HTTP GET
URL: https://<api-gateway-url>/api/devices/{id}
Example: GET https://api123.amazonaws.com/api/devices/id1
```
### Response 2 - Success:
```
HTTP 200 OK
Body (application/json):
{
"id": "/devices/id1",
"deviceModel": "/devicemodels/id1",
"name": "Sensor",
"note": "Testing a sensor.",
"serial": "A020000102"
}
```
### Response 2 - Failure 1:
```
HTTP 404 Not Found
If the request id does not exist.
```
### Response 1 - Failure 2:
```
HTTP 500 Internal Server Error
If any exceptional situation occurs on the server side.
```
## Project Architecture
This project is small and simple so I implement Service architecture that means a lambda function can handle different action (responds to Http GET & Http POST). Micro-service architecture is another way that each lambda function is responsible for one action only.

## How to Run

1. Clone the project to `/src `directory Go uses for its workspaces. Use `cd src/simple-restful-api-aws` for navigatating to project folder.

2. Type below to build

```
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o bin/devices simple-restful-api-aws/device
```

3. Make a zip of executable device file

`zip -j bin/devices.zip bin/devices`


4. Deploy to AWS by serverless

`sls deploy`


## How to Test

To post a device use the following command. Make sure to change  `<rest-api-id>` . You can get it from the link shown after deploy

```
curl -i -H "Content-Type: application/json" -X POST https://<rest-api-id>.execute-api.us-east-1.amazonaws.com/api/devices -d '{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"","note":"Testing a sensor.”,"serial":"A020000102"}'
```

To get a device from database you can use this comand:

```
curl -i https://<rest-api-id>.execute-api.us-east-1.amazonaws.com/api/devices/id1
```


## Unit Test

I put tests in `main_test.go` and `dataLayer_test.go` . Total coverage of statement is 97.7% . Separate apart coverage of `main_test.go` is 96.9% and `dataLayer_test.go`  is 100%.

To see the test unit coverage go to /devices folder by cd devices and execute following:

```
go test -coverprofile=cover.out
```

To see in Html:

```
go tool cover -html=cover.out
```
