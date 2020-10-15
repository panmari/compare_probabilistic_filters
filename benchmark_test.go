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
	cuckooVed "github.com/vedhavyas/cuckoo-filter"
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

func BenchmarkInsert(b *testing.B) {
	b.Run("BloomFilter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
			for _, w := range words {
				f.Add(bloomHash(w))
			}
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := bbloom.New(float64(numWords), 0.002)
			for _, w := range words {
				f.Add([]byte(w))
			}
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckoo.NewFilter(uint(numWords))
			for _, w := range words {
				f.Insert([]byte(w))
			}
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckooV2.NewFilter(uint(numWords))
			for _, w := range words {
				f.Insert([]byte(w))
			}
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckooVed.NewFilter(uint32(numWords))
			for _, w := range words {
				f.Insert([]byte(w))
			}
		}
	})
}

func BenchmarkContainsTrue(b *testing.B) {
	b.Run("BloomFilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words {
				f.Contains(bloomHash(w))
			}
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words {
			f.Add([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words {
				f.Has([]byte(w))
			}
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
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
	})
}

func BenchmarkContainsTruePanmariCuckoo(b *testing.B) {
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

func BenchmarkContainsTrueVedhavyasCuckoo(b *testing.B) {
	f := cuckooVed.NewFilter(uint32(numWords))
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
	f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
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
	f := bbloom.New(float64(numWords), 0.002)
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

func BenchmarkContainsFalseSeiflotfyCuckoo(b *testing.B) {
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

func BenchmarkContainsFalsePanmariCuckoo(b *testing.B) {
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

func BenchmarkContainsFalseVedhavyasCuckoo(b *testing.B) {
	f := cuckooVed.NewFilter(uint32(numWords))
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
