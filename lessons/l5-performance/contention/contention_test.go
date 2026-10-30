package contention

import (
	"testing"
)

func BenchmarkMixed(b *testing.B) {
	cases := []struct {
		name    string
		readOp  func(*Incr) int64
		writeOp func(*Incr)
	}{
		{name: "mutex", readOp: (*Incr).ReadMutex, writeOp: (*Incr).IncrMutex},
		{name: "rwmutex", readOp: (*Incr).ReadRWMutex, writeOp: (*Incr).IncrRWMutex},
		{name: "atomic", readOp: (*Incr).ReadAtomic, writeOp: (*Incr).IncrAtomic},
	}

	for _, tv := range cases {
		b.Run(tv.name, func(b *testing.B) {
			var i Incr
			b.RunParallel(func(pb *testing.PB) {
				op := 0
				for pb.Next() {
					if op%10 == 0 {
						tv.writeOp(&i)
					} else {
						_ = tv.readOp(&i)
					}
					op++
				}
			})
		})
	}
}
