package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"testing"
)

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	payload map[string]string // Store expected return values
	err     error
}

// Mock PutItem such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	output := new(dynamodb.PutItemOutput)
	output.Attributes = make(map[string]*dynamodb.AttributeValue)

	var err error

	for key, _ := range fd.payload {

		if key == "id" && fd.payload[key] == "id_test" {
			output.SetAttributes(
				map[string]*dynamodb.AttributeValue{
					"id":          &dynamodb.AttributeValue{S: aws.String("id_test")},
					"deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
					"name":        &dynamodb.AttributeValue{S: aws.String("name_test")},
					"note":        &dynamodb.AttributeValue{S: aws.String("note_test")},
					"serial":      &dynamodb.AttributeValue{S: aws.String("serial_test")},
				},
			)
			err = nil

		}

	}
	return output, err
}

func TestPutItemToDB(t *testing.T) {

	d := device{"id_test", "deviceModel_test", "name_test", "note_test", "serial_test"}

	db := new(Database)
	db.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{
			"id":          "id_test",
			"deviceModel": "deviceModel_test",
			"name":        "name_test",
			"note":        "note_test",
			"serial":      "serial_test"},
	}
	db.postToDB(d)

}

// Mock GetItem such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	output := new(dynamodb.GetItemOutput)
	output.Item = make(map[string]*dynamodb.AttributeValue)

	var err error

	for key, _ := range fd.payload {

		if key == "id" && fd.payload[key] == "id_test" {
			output.SetItem(
				map[string]*dynamodb.AttributeValue{
					"id":          &dynamodb.AttributeValue{S: aws.String("id_test")},
					"deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
					"name":        &dynamodb.AttributeValue{S: aws.String("name_test")},
					"note":        &dynamodb.AttributeValue{S: aws.String("note_test")},
					"serial":      &dynamodb.AttributeValue{S: aws.String("serial_test")},
				},
			)
			err = nil

		}

		if key == "id" && fd.payload[key] == "id_test1" {
			output.SetItem(nil)
			err = errors.New("error")

		}

		if key == "id" && fd.payload[key] == "id_test2" {
			output.SetItem(nil)
			err = nil

		}

	}
	return output, err
}

func TestGetItemFromDB(t *testing.T) {

	expectedValue := device{"id_test", "deviceModel_test", "name_test", "note_test", "serial_test"}
	expectedKey := "id_test"

	db := new(Database)
	db.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{
			"id":          "id_test",
			"deviceModel": "deviceModel_test",
			"name":        "name_test",
			"note":        "note_test",
			"serial":      "serial_test"},
	}

	if actualValue, err := db.getFromDB(expectedKey); actualValue != expectedValue {
		t.Errorf("Expected %q but got %q - error %q", expectedValue, actualValue, err)
	}

	expectedValue = device{"id_test1", "deviceModel_test", "name_test", "note_test", "serial_test"}
	expectedKey = "id_test1"

	db = new(Database)
	db.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{
			"id":          "id_test1",
			"deviceModel": "deviceModel_test",
			"name":        "name_test",
			"note":        "note_test",
			"serial":      "serial_test"},
	}

	if actualValue, err := db.getFromDB(expectedKey); err == nil {
		t.Errorf("Expected %q but got %q - error %q", expectedValue, actualValue, err)
	}

	expectedValue = device{"id_test2", "deviceModel_test", "name_test", "note_test", "serial_test"}
	expectedKey = "id_test2"
	db = new(Database)
	db.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{
			"id":          "id_test2",
			"deviceModel": "deviceModel_test",
			"name":        "name_test",
			"note":        "note_test",
			"serial":      "serial_test"},
	}

	if actualValue, err := db.getFromDB(expectedKey); actualValue.ID != "" {
		t.Errorf("Expected %q but got %q - error %q", expectedValue, actualValue, err)
	}

}
