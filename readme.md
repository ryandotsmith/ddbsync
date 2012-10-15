# ddbsync

DynamoDB/sync

This package is designed to emulate the behaviour of `pkg/sync` on top of Amazon's DynamoDB. If you need a distributed locking mechanism, consider using this package and DynamoDB before standing up paxos or Zookeeper.

[GoPkgDoc](http://go.pkgdoc.org/github.com/ryandotsmith/ddbsync)

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
		"github.com/ryandotsmith/ddbsync"
)

func main() {
		m := new(ddbsync.Mutex)
		m.Name = "some-name"
		m.Ttl = 10 * time.Second
		m.Lock()
		defer m.Unlock()
		// do important work here
		return
}
```

```bash
$ go get github.com/ryandotsmith/ddbsync
$ go run main.go
```

## Related

[lock-smith](https://github.com/ryandotsmith/lock-smith)
