# ddbsync

DynamoDB/sync

This package is designed to emulate the behaviour of `pkg/sync` on top of Amazon's DynamoDB. If you need a distributed locking mechanism, consider using this package and DynamoDB before standing up paxos of Zookeeper.

## Usage

Create a DynamoDB table named *Locks*.

```bash
$ export AWS_ACCESS_KEY=access
$ export AWS_SECRET_KEY=secret
```

```go
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

## Related

[lock-smith](https://github.com/ryandotsmith/lock-smith)
