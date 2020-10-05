// Benchmark for thread-unsafe interactions with probabilistic filters.
// Note that github.com/steakknife/bloomfilter doesn't allow interacting
// in a thread-unsafe way, leading to higher numbers there.
package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/AndreasBriese/bbloom"
	cuckooV2 "github.com/panmari/cuckoofilter"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/steakknife/bloomfilter"
)

var (
	words      []string
	otherWords []string
	numWords   = 500
)

func init() {
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)
	for i := 0; i < numWords && scanner.Scan(); i++ {
		words = append(words, scanner.Text())
	}
	for i := 0; i < numWords && scanner.Scan(); i++ {
		otherWords = append(otherWords, scanner.Text())
	}
}

func BenchmarkInsertBloomFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.001)
		for _, w := range words {
			f.Add(bloomHash(w))
		}
	}
}

func BenchmarkInsertBBloom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := bbloom.New(float64(numWords), 0.001)
		for _, w := range words {
			f.Add([]byte(w))
		}
	}
}

func BenchmarkInsertCuckoo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words {
			f.Insert([]byte(w))
		}
	}
}

func BenchmarkInsertCuckooV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := cuckooV2.NewFilter(uint(numWords))
		for _, w := range words {
			f.Insert([]byte(w))
		}
	}
}

func BenchmarkContainsTrueBloom(b *testing.B) {
	f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.001)
	for _, w := range words {
		f.Add(bloomHash(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			f.Contains(bloomHash(w))
		}
	}
}

func BenchmarkContainsTrueBBloom(b *testing.B) {
	f := bbloom.New(float64(numWords), 0.001)
	for _, w := range words {
		f.Add([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			f.Has([]byte(w))
		}
	}
}

func BenchmarkContainsTrueCuckoo(b *testing.B) {
	f := cuckoo.NewFilter(uint(numWords))
	for _, w := range words {
		f.Insert([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			f.Lookup([]byte(w))
		}
	}
}

func BenchmarkContainsTrueCuckooV2(b *testing.B) {
	f := cuckooV2.NewFilter(uint(numWords))
	for _, w := range words {
		f.Insert([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			f.Lookup([]byte(w))
		}
	}
}

func BenchmarkContainsFalseBloom(b *testing.B) {
	f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.001)
	for _, w := range words {
		f.Add(bloomHash(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range otherWords {
			f.Contains(bloomHash(w))
		}
	}
}

func BenchmarkContainsFalseBBloom(b *testing.B) {
	f := bbloom.New(float64(numWords), 0.001)
	for _, w := range words {
		f.Add([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range otherWords {
			f.Has([]byte(w))
		}
	}
}

func BenchmarkContainsFalseCuckoo(b *testing.B) {
	f := cuckoo.NewFilter(uint(numWords))
	for _, w := range words {
		f.Insert([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range otherWords {
			f.Lookup([]byte(w))
		}
	}
}

func BenchmarkContainsFalseCuckooV2(b *testing.B) {
	f := cuckooV2.NewFilter(uint(numWords))
	for _, w := range words {
		f.Insert([]byte(w))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range otherWords {
			f.Lookup([]byte(w))
		}
	}
}
