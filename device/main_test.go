//We used testify that makes it easy to mock
package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"

	"testing"
)

// mockedObject is a mocked object that implements an interface
// that describes an object that the code I am testing relies on.
type mockedObject struct {
	mock.Mock
}

//It's a mock object to test postDevice
func (m *mockedObject) postToDB(d device) error {
	args := m.Called(d)
	return args.Error(0)
}

func TestPostDevice(t *testing.T) {

	testObj := new(mockedObject)
	handler := Handler{testObj}

	//what the mocked postToDB has to return by specified input
	d := device{"id1", "model1", "sensor1", "testing a sensor1", "serial1"}
	testObj.On("postToDB", d).Return(nil)

	js, _ := json.Marshal(d)
	req := events.APIGatewayProxyRequest{Body: string(js)}
	// if postDevice pass the test by specified request
	assert.Equal(t, handler.postDevice(req), events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string("HTTP 201 Created"),
	})

	//serial is empty so postToDB has to return error
	d = device{"idSerialEmpty", "model", "sensor", "testing a sensor", ""}
	testObj.On("postToDB", d).Return(errors.New("empty fields"))

	js, _ = json.Marshal(d)
	req = events.APIGatewayProxyRequest{Body: string(js)}
	assert.Equal(t, handler.postDevice(req), events.APIGatewayProxyResponse{

		StatusCode: http.StatusBadRequest, // error 400
		Body:       string("Empty Field is not valid,Check the following: Serial, "),
	})

	d = device{"idServerError", "model", "sensor", "testing a sensor", "serial"}
	testObj.On("postToDB", d).Return(errors.New("internal server error"))

	js, _ = json.Marshal(d)
	req = events.APIGatewayProxyRequest{Body: string(js)}

	assert.Equal(t, handler.postDevice(req), events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError, //error 500
		Body:       http.StatusText(http.StatusInternalServerError),
	})

}

//It's a mock object to test GetDevice
func (m *mockedObject) getFromDB(value string) (device, error) {
	args := m.Called(value)
	return args.Get(0).(device), args.Error(1)
}

func TestGetDevice(t *testing.T) {
	testObj := new(mockedObject)

	d := device{"id4", "id1", "sensor1", "testing a sensor1", "serial1"}

	//what the mocked getFromDB has to return by specified input
	testObj.On("getFromDB", "id1").Return(d, nil)
	testObj.On("getFromDB", "id doesn't existed").Return(device{}, nil)
	testObj.On("getFromDB", "").Return(device{}, errors.New("error 500"))

	handler := Handler{testObj}

	req := events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "id1"}}
	js, _ := json.Marshal(d)
	assert.Equal(t, handler.getDevice(req), events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, //HTTP 200
		Body:       string(js),
	})

	req = events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "id doesn't existed"}}
	assert.Equal(t, handler.getDevice(req), events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound, //HTTP 404
		Body:       http.StatusText(http.StatusNotFound),
	})

	req = events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": ""}}
	assert.Equal(t, handler.getDevice(req), events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError, //HTTP 500
		Body:       http.StatusText(http.StatusInternalServerError),
	})

}

// test httpSpecifier
func TestHttpSpecifier(t *testing.T) {

	//-----------------------------------------------
	req := events.APIGatewayProxyRequest{HTTPMethod: "GET"}
	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}

	expectedValue, _ := httpSpecifier(req)
	if expectedValue.StatusCode == http.StatusInternalServerError {
		assert.Equal(t, expectedValue.StatusCode, res.StatusCode)
	} else {
		res := events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       http.StatusText(http.StatusNotFound),
		}
		assert.Equal(t, expectedValue.StatusCode, res.StatusCode)
	}

	//-----------------------------------------------
	req = events.APIGatewayProxyRequest{HTTPMethod: "POST"}
	res = events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       http.StatusText(http.StatusBadRequest),
	}

	expectedValue, _ = httpSpecifier(req)
	assert.Equal(t, expectedValue.StatusCode, res.StatusCode)

	//-----------------------------------------------
	req = events.APIGatewayProxyRequest{HTTPMethod: "ERROR"}
	res = events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}

	expectedValue, _ = httpSpecifier(req)
	assert.Equal(t, expectedValue.StatusCode, res.StatusCode)
}
