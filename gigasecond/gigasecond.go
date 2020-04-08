package gigasecond

import (
        "time"
        "fmt"
)

func AddGigasecond(t time.Time) time.Time {
	t = t.Add(time.Second * 1000000000)
	fmt.Printf("Final time = %v\n", t)
	return t
}
