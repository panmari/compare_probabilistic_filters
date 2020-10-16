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

const maxNumWords = 50000

func init() {
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)
	for i := 0; i < maxNumWords && scanner.Scan(); i++ {
		words = append(words, scanner.Text())
	}
	for i := 0; i < maxNumWords && scanner.Scan(); i++ {
		otherWords = append(otherWords, scanner.Text())
	}
}

func BenchmarkFilters(b *testing.B) {
	for _, n := range []int{500, 2000, 10000, 50000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			if n > maxNumWords {
				b.Fatalf("Num words too large: %d > %d", n, maxNumWords)
			}
			numWords = n
			b.Run("Insert", insert)
			b.Run("ContainsTrue", containsTrue)
			b.Run("ContainsFalse", containsFalse)
		})
	}
}

func insert(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
			for _, w := range words[:numWords] {
				f.Add(bloomHash(w))
			}
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := bbloom.New(float64(numWords), 0.002)
			for _, w := range words[:numWords] {
				f.Add([]byte(w))
			}
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckoo.NewFilter(uint(numWords))
			for _, w := range words[:numWords] {
				f.Insert([]byte(w))
			}
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckooV2.NewFilter(uint(numWords))
			for _, w := range words[:numWords] {
				f.Insert([]byte(w))
			}
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := cuckooVed.NewFilter(uint32(numWords))
			for _, w := range words[:numWords] {
				f.Insert([]byte(w))
			}
		}
	})
}

func containsTrue(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words[:numWords] {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words[:numWords] {
				f.Contains(bloomHash(w))
			}
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words[:numWords] {
			f.Add([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words[:numWords] {
				f.Has([]byte(w))
			}
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		f := cuckooV2.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		f := cuckooVed.NewFilter(uint32(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range words[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
}

func containsFalse(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words[:numWords] {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range otherWords[:numWords] {
				f.Contains(bloomHash(w))
			}
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words[:numWords] {
			f.Add([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range otherWords[:numWords] {
				f.Has([]byte(w))
			}
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range otherWords[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		f := cuckooV2.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range otherWords[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		f := cuckooVed.NewFilter(uint32(numWords))
		for _, w := range words[:numWords] {
			f.Insert([]byte(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, w := range otherWords[:numWords] {
				f.Lookup([]byte(w))
			}
		}
	})
}
