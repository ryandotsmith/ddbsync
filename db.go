package ddbsync

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/zencoder/ddbsync/models"
)

const DEFAULT_LOCKS_TABLE_NAME string = "Locks"

type database struct {
	client AWSDynamoer
}

var region string = os.Getenv("DDBSYNC_DYNAMODB_REGION")
var endpoint string = os.Getenv("DDBSYNC_DYNAMODB_ENDPOINT")

func disableSSL() bool {
	b, err := strconv.ParseBool(os.Getenv("DDBSYNC_DYNAMODB_DISABLE_SSL"))
	if err != nil {
		return false
	}
	return b
}

var db DBer = &database{
	client: dynamodb.New(&aws.Config{
		Endpoint:   endpoint,
		Region:     region,
		DisableSSL: disableSSL(),
	}),
}

var _ DBer = (*database)(nil) // Forces compile time checking of the interface

var _ AWSDynamoer = (*dynamodb.DynamoDB)(nil) // Forces compile time checking of the interface

type DBer interface {
	Put(string, int64) error
	Get(string) (*models.Item, error)
	Delete(string) error
}

type AWSDynamoer interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

func locksTableName() string {
	tableName := DEFAULT_LOCKS_TABLE_NAME
	env := os.Getenv("DDBSYNC_LOCKS_TABLE_NAME")
	if env != "" {
		tableName = env
	}
	return tableName
}

func (db *database) Put(name string, created int64) error {
	i := map[string]*dynamodb.AttributeValue{
		"Name": &dynamodb.AttributeValue{
			S: aws.String(name),
		},
		"Created": &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(created, 10)),
		},
	}

	e := map[string]*dynamodb.ExpectedAttributeValue{
		"Name": &dynamodb.ExpectedAttributeValue{
			Exists: aws.Boolean(false),
		},
	}

	pit := &dynamodb.PutItemInput{
		TableName: aws.String(locksTableName()),
		Item:      &i,
		Expected:  &e,
	}
	_, err := db.client.PutItem(pit)
	if err != nil {
		return err
	}

	return nil
}

func (db *database) Get(name string) (*models.Item, error) {
	kc := map[string]*dynamodb.Condition{
		"Name": &dynamodb.Condition{
			AttributeValueList: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{
					S: aws.String(name),
				},
			},
			ComparisonOperator: aws.String("EQ"),
		},
	}
	qi := &dynamodb.QueryInput{
		TableName:       aws.String(locksTableName()),
		ConsistentRead:  aws.Boolean(true),
		Select:          aws.String("SPECIFIC_ATTRIBUTES"),
		AttributesToGet: []*string{aws.String("Name"), aws.String("Created")},
		KeyConditions:   &kc,
	}

	qo, err := db.client.Query(qi)
	if err != nil {
		return nil, err
	}

	// Make sure that no or 1 item is returned from DynamoDB
	if qo.Count != nil {
		if *qo.Count == 0 {
			return nil, errors.New(fmt.Sprintf("No item for Name, %s", name))
		} else if *qo.Count > 1 {
			return nil, errors.New(fmt.Sprintf("Expected only 1 item returned from Dynamo, got %d", *qo.Count))
		}
	} else {
		return nil, errors.New("Count not returned")
	}

	if len(qo.Items) < 1 || qo.Items[0] == nil {
		return nil, errors.New("No item returned, count is invalid.")
	}

	n := ""
	c := int64(0)
	for index, element := range *qo.Items[0] {
		if index == "Name" {
			n = *element.S
		}
		if index == "Created" {
			c, _ = strconv.ParseInt(*element.N, 10, 0)
		}
	}
	if n == "" || c == 0 {
		return nil, errors.New("The Name and Created keys were not found in the Dynamo result")
	}
	i := &models.Item{n, c}
	return i, nil
}

func (db *database) Delete(name string) error {
	k := map[string]*dynamodb.AttributeValue{
		"Name": &dynamodb.AttributeValue{
			S: aws.String(name),
		},
	}
	dii := &dynamodb.DeleteItemInput{
		TableName: aws.String(locksTableName()),
		Key:       &k,
	}
	_, err := db.client.DeleteItem(dii)
	if err != nil {
		return err
	}

	return nil
}
