package quantity_test

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"

	"github.com/chipaca/quantity"
	"gopkg.in/cheggaaa/pb.v1"
)

var wg sync.WaitGroup
var stripeWidth = uint32(runtime.NumCPU())

func testnWidthStripe(pbi *pb.ProgressBar, w int, n uint32, t *testing.T) {
outer:
	for {
		for j := 0; j < 10000; j++ {
			s := quantity.FormatAmount(uint64(n), w)
			if len(s) != w {
				t.Fatalf("formatting %v, expecting something of length %d but got %q", n, w, s)
			}
			m := n + stripeWidth
			if m < n {
				break outer
			}
			n = m
		}
		pbi.Set64(int64(n))
	}
	wg.Done()
	pbi.Finish()
}

func TestFormatAmountAlwaysFixedWidth(t *testing.T) {
	pbs := make([]*pb.ProgressBar, int(stripeWidth))
	for i := range pbs {
		pbs[i] = pb.New64(math.MaxUint32).Prefix(fmt.Sprint(i))
		pbs[i].ShowSpeed = true
		pbs[i].ShowPercent = false
	}
	pool, err := pb.StartPool(pbs...)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i := uint32(0); i < stripeWidth; i++ {
		w := 5
		go testnWidthStripe(pbs[i], w, i, t)
		wg.Add(1)
	}

	wg.Wait()
	pool.Stop()
}
