package ddbsync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bmizerany/aws4"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	ErrNotFound = errors.New("Not Found")
)

const (
	OpPutItem    = "DynamoDB_20111205.PutItem"
	OpGetItem    = "DynamoDB_20111205.GetItem"
	OpDeleteItem = "DynamoDB_20111205.DeleteItem"
)

type Item struct {
	Name    string
	Created int64
}

type responseError struct {
	resp *http.Response
}

type B struct {
	B bool
}

type S struct {
	S string
}

type N struct {
	N int64 `json:",string"`
}

type U struct {
	Action string
	Value  interface{}
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
			Name    S
			Created N
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

	resp, err := db.do(OpPutItem, t)
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

func (db *database) get(name string) (*Item, error) {
	type T struct {
		TableName      string
		ConsistentRead bool
		Key            struct {
			HashKeyElement S
		}
		AttributesToGet []string
	}

	t := new(T)
	t.TableName = "Locks"
	t.ConsistentRead = true
	t.Key.HashKeyElement.S = name
	t.AttributesToGet = []string{"Name", "Created"}

	resp, err := db.do(OpGetItem, t)
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
			Name    S
			Created N
		}
	}
	r := new(R)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, err
	}

	if r.Item.Name.S == "" {
		return nil, nil
	}
	return &Item{r.Item.Name.S, r.Item.Created.N}, nil
}

func (db *database) delete(name string) error {
	type T struct {
		TableName string
		Key       struct {
			HashKeyElement S
		}
	}

	t := new(T)
	t.TableName = "Locks"
	t.Key.HashKeyElement.S = name

	resp, err := db.do(OpDeleteItem, t)
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
