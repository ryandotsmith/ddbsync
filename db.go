package ddbsync

import (
	"errors"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"log"
	"strconv"
)

type item struct {
	Name    string
	Created int64
}

type database struct {
	client *dynamodb.DynamoDB
}

var db DBer = &database{
	client: dynamodb.New(nil),
}

type DBer interface {
	Put(string, int64) error
	Get(string) (*item, error)
	Delete(string) error
}

func (db *database) Put(name string, created int64) error {
	log.Printf("put called. name = %s, created = %d", name, created)
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
	pio, err := db.client.PutItem(pit)
	if err != nil {
		log.Printf("put. Error = %s", err.Error())
		return err
	}

	log.Printf("PutItem finished. name = %s, created = %d, pio = %s", name, created, pio)

	return nil
}

func (db *database) Get(name string) (*item, error) {
	log.Printf("get called. name = %s", name)
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

	//Make sure that no or 1 item is returned from DynamoDB
	if qo.Count != nil {
		if *qo.Count == 0 {
			eStr := fmt.Sprintf("No item for Name, %s", name)
			return nil, errors.New(eStr)
		} else if *qo.Count > 1 {
			eStr := fmt.Sprintf("Expected only 1 item returned from Dynamo, got %d", *qo.Count)
			return nil, errors.New(eStr)
		}
	} else {
		return nil, errors.New("Count not returned")
	}

	n := *qo.Items[0]["Name"].S
	c, _ := strconv.ParseInt(*qo.Items[0]["Created"].N, 10, 0)
	i := &item{n, c}
	log.Println("get. name = %s, i = %s", name, i)
	return i, nil
}

func (db *database) Delete(name string) error {
	log.Printf("delete called. name = %s", name)
	k := map[string]dynamodb.AttributeValue{
		"Name": dynamodb.AttributeValue{
			S: aws.String(name),
		},
	}
	dii := &dynamodb.DeleteItemInput{
		TableName: aws.String("Locks"),
		Key:       k,
	}
	dio, err := db.client.DeleteItem(dii)
	if err != nil {
		return err
	}

	log.Printf("Deleted item. name = %s, dio = %s", name, dio)

	return nil
}
