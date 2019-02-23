//handle communication with Database
package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type device struct {
	ID     string `json:"id"`
	Model  string `json:"deviceModel"`
	Name   string `json:"name"`
	Note   string `json:"note"`
	Serial string `json:"serial"`
}

func (db Database) postToDB(dev device) error {
	device := &dynamodb.PutItemInput{
		TableName: aws.String("Devices"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(dev.ID),
			},
			"DeviceModel": {
				S: aws.String(dev.Model),
			},
			"Name": {
				S: aws.String(dev.Name),
			},
			"Note": {
				S: aws.String(dev.Note),
			},
			"Serial": {
				S: aws.String(dev.Serial),
			},
		},
	}
	//post Device to database
	_, err := db.DynamoDB.PutItem(device)
	return err

}

func (db Database) getFromDB(id string) (device, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Devices"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				// add /devices/ string before id value
				S: aws.String("/devices/" + id),
			},
		},
	}

	result, err := db.DynamoDB.GetItem(input)
	if err != nil {
		return device{}, err
	}
	if result.Item == nil {
		return device{}, nil
	}

	device := new(device)
	dynamodbattribute.UnmarshalMap(result.Item, device)

	return *device, nil
}
