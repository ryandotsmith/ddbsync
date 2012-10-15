package ddbsync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/aws4"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	opPutItem    = "DynamoDB_20111205.PutItem"
	opGetItem    = "DynamoDB_20111205.GetItem"
	opDeleteItem = "DynamoDB_20111205.DeleteItem"
)

type item struct {
	Name    string
	Created int64
}

type responseError struct {
	resp *http.Response
}

type s struct {
	S string
}

type n struct {
	N int64 `json:",string"`
}

type database struct {
	keys *aws4.Keys
	s    *aws4.Service
}

var db = &database{
	keys: &aws4.Keys{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	},
	s: &aws4.Service{
		Name:   "dynamodb",
		Region: "us-east-1",
	},
}

func (db *database) put(name string, created int64) error {
	type T struct {
		TableName string
		Item      struct {
			Name    s
			Created n
		}
		Expected struct {
			Name struct {
				Exists bool
			}
		}
	}

	t := new(T)
	t.TableName = "Locks"
	t.Item.Name.S = name
	t.Item.Created.N = created
	t.Expected.Name.Exists = false

	resp, err := db.do(opPutItem, t)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("updateMinute error: %d %q", resp.StatusCode, string(b))
	}

	return nil
}

func (db *database) get(name string) (*item, error) {
	type T struct {
		TableName      string
		ConsistentRead bool
		Key            struct {
			HashKeyElement s
		}
		AttributesToGet []string
	}

	t := new(T)
	t.TableName = "Locks"
	t.ConsistentRead = true
	t.Key.HashKeyElement.S = name
	t.AttributesToGet = []string{"Name", "Created"}

	resp, err := db.do(opGetItem, t)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("get error: %d %q", resp.StatusCode, string(b))
	}

	type R struct {
		Item struct {
			Name    s
			Created n
		}
	}
	r := new(R)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, err
	}

	if r.Item.Name.S == "" {
		return nil, nil
	}
	return &item{r.Item.Name.S, r.Item.Created.N}, nil
}

func (db *database) delete(name string) error {
	type T struct {
		TableName string
		Key       struct {
			HashKeyElement s
		}
	}

	t := new(T)
	t.TableName = "Locks"
	t.Key.HashKeyElement.S = name

	resp, err := db.do(opDeleteItem, t)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("delete error: %d %q", resp.StatusCode, string(b))
	}
	return nil
}

func (db *database) do(op string, v interface{}) (*http.Response, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}

	r, _ := http.NewRequest("POST", "https://dynamodb.us-east-1.amazonaws.com/", b)
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("X-Amz-Target", op)
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")

	err := db.s.Sign(db.keys, r)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}
