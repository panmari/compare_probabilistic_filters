package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"os"
	"runtime"

	"github.com/AndreasBriese/bbloom"
	"github.com/irfansharif/cfilter"
	cuckooLin "github.com/linvon/cuckoo-filter"
	cuckooV2 "github.com/panmari/cuckoofilter"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/steakknife/bloomfilter"
	cuckooVed "github.com/vedhavyas/cuckoo-filter"
)

// Inserts the size of wordlist times this items into the filters.
var (
	wordListMultiplier = flag.Int("word_list_multiplier", 250, "Determines the number of inserted items.")
)

func main() {
	words := readWords()

	testCases := []struct {
		name       string
		testFilter func([]string) filterStats
	}{
		{
			"steakknife/bloomfilter",
			testBloomfilter,
		},
		{
			"AndreasBriese/bbloom",
			testBbloom,
		},
		{
			"seiflotfy/cuckoofilter",
			testCuckoofilter,
		},
		{
			"panmari/cuckoofilter",
			testCuckoofilterV2,
		}, {
			"vedhavyas/cuckoo-filter",
			testCuckoofilterVed,
		}, {
			"linvon/cuckoo-filter",
			testCuckoofilterLin,
		},
		{
			// panic: runtime error: index out of range
			"irfansharif/cfilter",
			testCfilter,
		},
	}
	// for _, size := range []int{10, 50, 100, 150, 200, 250, 300, 350, 400, 450} {
	// *wordListMultiplier = size
	for _, tc := range testCases {
		stats := tc.testFilter(words)
		fpRate := float64(stats.fp) / (float64(stats.fp + stats.tn))
		const megabyte = 1 << 20
		memMB := float64(stats.mem) / float64(megabyte)
		fmt.Printf("%s: size=%d, mem=%.3f MB, insertFailed=%d, fn=%d, fp=%d, fp_rate=%f\n",
			tc.name, filterSize(words), memMB, stats.insertFailed, stats.fn, stats.fp, fpRate)
	}
	// }
}

type filterStats struct {
	insertFailed, tp, fp, tn, fn, mem int64
}

func filterSize(words []string) int {
	return len(words) * (*wordListMultiplier + 1)
}

func testBloomfilter(words []string) filterStats {
	memBefore := heapAllocs()
	bf, err := bloomfilter.NewOptimal(uint64(filterSize(words)), 0.0001)
	if err != nil {
		log.Fatalf("failed creating bloom filter with size %d: %v", len(words), err)
	}

	insert := func(s string) bool { bf.Add(bloomHash(s)); return true }
	contains := func(s string) bool { return bf.Contains(bloomHash(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testBbloom(words []string) filterStats {
	memBefore := heapAllocs()
	bf := bbloom.New(float64(filterSize(words)), 0.002)

	insert := func(s string) bool { bf.Add([]byte(s)); return true }
	contains := func(s string) bool { return bf.Has([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilter(words []string) filterStats {
	memBefore := heapAllocs()
	cf := cuckoo.NewFilter(uint(filterSize(words)))

	insert := func(s string) bool { return cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterV2(words []string) filterStats {
	memBefore := heapAllocs()
	cf := cuckooV2.NewFilter(uint(filterSize(words)))

	insert := func(s string) bool { return cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterVed(words []string) filterStats {
	memBefore := heapAllocs()
	cf := cuckooVed.NewFilter(uint32(filterSize(words)))

	insert := func(s string) bool { return cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCfilter(words []string) filterStats {
	memBefore := heapAllocs()

	cf := cfilter.New(cfilter.Size(uint(filterSize(words))))

	insert := func(s string) bool { return cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterLin(words []string) filterStats {
	memBefore := heapAllocs()

	cf := cuckooLin.NewFilter(4, 13, uint(filterSize(words)), cuckooLin.TableTypePacked)

	insert := func(s string) bool { return cf.Add([]byte(s)) }
	contains := func(s string) bool { return cf.Contain([]byte(s)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testImplementation(words []string, memBefore uint64,
	insert func(string) bool, contains func(string) bool) (stats filterStats) {
	skip := func(i, j int) bool { return (i+j)%200 == 0 }
	for i, w1 := range words {
		insert(w1)
		for j, w2 := range append(words[0:*wordListMultiplier]) {
			if !skip(i, j) {
				w := w1 + w2
				if ok := insert(w); !ok {
					stats.insertFailed++
				}
			}
		}
	}
	memAfter := heapAllocs()
	stats.mem = int64(memAfter - memBefore)

	// Construct non-contained words in a second step in order to not influence
	// memory measurement above.
	remaining := make([]string, 0, len(words)/200)
	for i, w1 := range words {
		for j, w2 := range words[0:*wordListMultiplier] {
			w := w1 + w2
			if skip(i, j) {
				remaining = append(remaining, w)
			}
		}
	}

	for _, w := range remaining {
		if contains(w) {
			stats.fp++
		} else {
			stats.tn++
		}
	}
	for _, w := range words {
		if contains(w) {
			stats.tp++
		} else {
			stats.fn++
		}
	}
	return stats
}

func readWords() []string {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

func heapAllocs() uint64 {
	runtime.GC() // Run GC to clean up unreachable objects
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

func bloomHash(s string) hash.Hash64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h
}
