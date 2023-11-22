package fixedgopool

import (
	"fmt"
	"testing"
	"time"
)

func TestGoPool_AddJob(t *testing.T) {

	p := NewPool(
		WithWorkNum(10),
		WithJobQueueSize(200),
	)

	for i := 1; i <= 2000; i++ {
		tmp := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Println("job:", tmp)
		})
	}

	p.WaitJobs()
	p.Close()
}
