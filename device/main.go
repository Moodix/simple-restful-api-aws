package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"net/http"
)

type Communicator interface {
	getFromDB(string) (device, error)
	postToDB(bk device) error
}

type Database struct {
	DynamoDB dynamodbiface.DynamoDBAPI
}

type Handler struct {
	communicator Communicator
}

//setup the Database
var db = Database{dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))}

var handler = Handler{db}

//When a http POST is sent to server, Lambda triggers this function
func (a Handler) postDevice(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	receivedDevice := new(device)
	// convert http post request body to a device
	err := json.Unmarshal([]byte(req.Body), receivedDevice)

	//check payload fields
	emptyField := ""
	if receivedDevice.ID == "" {
		emptyField = emptyField + "ID, "
	}
	if receivedDevice.Model == "" {
		emptyField = emptyField + "Model, "
	}
	if receivedDevice.Name == "" {
		emptyField = emptyField + "Name, "
	}
	if receivedDevice.Note == "" {
		emptyField = emptyField + "Note, "
	}
	if receivedDevice.Serial == "" {
		emptyField = emptyField + "Serial, "
	}
	if emptyField != "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest, // error 400
			Body:       string("Empty Field is not valid,Check the following: " + emptyField),
		}
	}

	//post device to database
	err = a.communicator.postToDB(*receivedDevice)

	if err != nil {
		//internal server error occurred

		//log.New(os.Stderr, "ERROR ", log.Llongfile).Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError, //error 500
			Body:       http.StatusText(http.StatusInternalServerError),
		}
	}

	//Device has been posted. everything is ok!
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated, //Http 201
		Body:       string("HTTP 201 Created"),
	}
}

//When a http GET is sent to server, Lambda triggers this function
func (a Handler) getDevice(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	//get the value of path parameter that determines id of device
	id := req.PathParameters["id"]

	//get device from database
	device, err := a.communicator.getFromDB(id)

	if err != nil {
		//internal server error occurred
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError, //HTTP 500
			Body:       http.StatusText(http.StatusInternalServerError),
		}
	}

	if device.ID == "" {
		//there is'nt such a device in the database
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound, //HTTP 404
			Body:       http.StatusText(http.StatusNotFound),
		}
	}

	//everything is ok & we got the device from database!
	js, _ := json.Marshal(device)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, //HTTP 200
		Body:       string(js),
	}
}

func main() {
	lambda.Start(httpSpecifier)
}

//We implement service architecture so the lambda function is responsible for both both Http GET &  http Post
func httpSpecifier(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return handler.getDevice(req), nil
	case "POST":
		return handler.postDevice(req), nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       http.StatusText(http.StatusMethodNotAllowed),
		}, nil
	}
}
