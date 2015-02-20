package ddbsync

import (
	"errors"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/ryandotsmith/ddbsync/models"
	"strconv"
)

type database struct {
	client AWSDynamoer
}

var db DBer = &database{
	client: dynamodb.New(nil),
}

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

func (db *database) Put(name string, created int64) error {
	i := map[string]dynamodb.AttributeValue{
		"Name": dynamodb.AttributeValue{
			S: aws.String(name),
		},
		"Created": dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(created, 10)),
		},
	}

	e := map[string]dynamodb.ExpectedAttributeValue{
		"Name": dynamodb.ExpectedAttributeValue{
			Exists: aws.Boolean(false),
		},
	}

	pit := &dynamodb.PutItemInput{
		TableName: aws.String("Locks"),
		Item:      i,
		Expected:  e,
	}
	_, err := db.client.PutItem(pit)
	if err != nil {
		return err
	}

	return nil
}

func (db *database) Get(name string) (*models.Item, error) {
	kc := map[string]dynamodb.Condition{
		"Name": dynamodb.Condition{
			AttributeValueList: []dynamodb.AttributeValue{
				dynamodb.AttributeValue{
					S: aws.String(name),
				},
			},
			ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
		},
	}
	qi := &dynamodb.QueryInput{
		TableName:       aws.String("Locks"),
		ConsistentRead:  aws.Boolean(true),
		Select:          aws.String(dynamodb.SelectSpecificAttributes),
		AttributesToGet: []string{"Name", "Created"},
		KeyConditions:   kc,
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

	n := *qo.Items[0]["Name"].S
	c, _ := strconv.ParseInt(*qo.Items[0]["Created"].N, 10, 0)
	i := &models.Item{n, c}
	return i, nil
}

func (db *database) Delete(name string) error {
	k := map[string]dynamodb.AttributeValue{
		"Name": dynamodb.AttributeValue{
			S: aws.String(name),
		},
	}
	dii := &dynamodb.DeleteItemInput{
		TableName: aws.String("Locks"),
		Key:       k,
	}
	_, err := db.client.DeleteItem(dii)
	if err != nil {
		return err
	}

	return nil
}
