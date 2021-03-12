package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Error string

func (err Error) Error() string { return string(err) }

func sliceHelper(reply interface{}, err error, name string, makeSlice func(int), assign func(int, interface{}) error) error {
	if err != nil {
		return err
	}
	switch reply := reply.(type) {
	case []interface{}:
		makeSlice(len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			if err := assign(i, reply[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return redis.Nil
	case Error:
		return reply
	}

	return fmt.Errorf("cache: unexpected type for %s, got type %T", name, reply)
}

func strings(reply interface{}, err error) ([]string, error) {
	var result []string
	err = sliceHelper(reply, err, "strings",
		func(n int) { result = make([]string, n) }, func(i int, v interface{}) error {
			switch v := v.(type) {
			case string:
				result[i] = v
				return nil
			case []byte:
				result[i] = string(v)
				return nil
			default:
				return fmt.Errorf("cache: unexpected element type for strings, got type %T", v)
			}
		})

	if err != nil {
		if len(result) <= 0 {
			err = redis.Nil
		}
	}

	return result, err
}
