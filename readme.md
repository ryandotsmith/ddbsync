# ddbsync

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/zencoder/ddbsync)
[![Build Status](https://img.shields.io/travis/zencoder/ddbsync.svg)](https://travis-ci.org/zencoder/ddbsync)
[![Coverage Status](https://coveralls.io/repos/zencoder/ddbsync/badge.svg?branch=master)](https://coveralls.io/r/zencoder/ddbsync?branch=master)

DynamoDB/sync

This package is designed to emulate the behaviour of `pkg/sync` on top of Amazon's DynamoDB. If you need a distributed locking mechanism, consider using this package and DynamoDB before standing up paxos or Zookeeper.

## Usage

Create a DynamoDB table named *Locks*.

```bash
$ export AWS_ACCESS_KEY=access
$ export AWS_SECRET_KEY=secret
```

```go
// ./main.go

package main

import(
		"time"
		"github.com/zencoder/ddbsync"
)

func main() {
		m := new(ddbsync.Mutex)
		m.Name = "some-name"
		m.TTL = int64(10 * time.Second)
		m.Lock()
		defer m.Unlock()
		// do important work here
		return
}
```

```bash
$ go get github.com/zencoder/ddbsync
$ go run main.go
```

## Related

[ddbsync](https://github.com/ryandotsmith/ddbsync)
[lock-smith](https://github.com/ryandotsmith/lock-smith)
