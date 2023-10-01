package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"os"
	"runtime"

	"github.com/AndreasBriese/bbloom"
	cuckooLin "github.com/linvon/cuckoo-filter"
	cuckooLk "github.com/livekit/cuckoofilter"
	cuckooV2 "github.com/panmari/cuckoofilter"
	cuckooLocal "github.com/panmari/cuckoofilter_local"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/steakknife/bloomfilter"
	cuckooVed "github.com/vedhavyas/cuckoo-filter"
	"golang.org/x/exp/slices"
)

// Inserts the size of wordlist times this items into the filters.
var (
	wordListMultiplier = flag.Int("word_list_multiplier", 250, "Determines the number of inserted items. Word list length times this multiplier entries are inserted.")
	wordListPath       = flag.String("word_list_path", "/usr/share/dict/words", "Path to list with words")
)

func main() {
	words := readWords()

	testCases := []struct {
		name       string
		testFilter func([][]byte) filterStats
	}{
		{
			"steakknife/bloomfilter",
			testBloomfilter,
		}, {
			"AndreasBriese/bbloom",
			testBbloom,
		}, {
			"seiflotfy/cuckoofilter",
			testCuckoofilter,
		}, {
			"panmari/cuckoofilter",
			testCuckoofilterV2,
		}, {
			"panmari/cuckoofilter/low",
			testCuckoofilterV2Low,
		}, {
			"panmari/cuckoofilter/fastrand",
			testCuckoofilterFastrand,
		}, {
			"livekit/cuckoofilter",
			testCuckoofilterLk,
		}, {
			"vedhavyas/cuckoo-filter",
			testCuckoofilterVed,
		}, {
			"linvon/cuckoo-filter/single",
			testCuckoofilterLinSingle,
		}, {
			"linvon/cuckoo-filter/packed",
			testCuckoofilterLinPacked,
			// }, {
			// 	// panic: runtime error: index out of range
			// 	"irfansharif/cfilter",
			// 	testCfilter,
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

func filterSize(words [][]byte) int {
	return *wordListMultiplier * len(words)
}

func testBloomfilter(words [][]byte) filterStats {
	memBefore := heapAllocs()
	bf, err := bloomfilter.NewOptimal(uint64(filterSize(words)), 0.0001)
	if err != nil {
		log.Fatalf("failed creating bloom filter with size %d: %v", len(words), err)
	}

	insert := func(b []byte) bool { bf.Add(bloomHash(b)); return true }
	contains := func(b []byte) bool { return bf.Contains(bloomHash(b)) }
	return testImplementation(words, memBefore, insert, contains)
}

func testBbloom(words [][]byte) filterStats {
	memBefore := heapAllocs()
	bf := bbloom.New(float64(filterSize(words)), 0.002)

	insert := func(b []byte) bool { bf.Add(b); return true }
	contains := func(b []byte) bool { return bf.Has(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilter(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckoo.NewFilter(uint(filterSize(words)))

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterV2(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooV2.NewFilter(cuckooV2.Config{NumElements: uint(filterSize(words))})

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterV2Low(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooV2.NewFilter(cuckooV2.Config{
		NumElements: uint(filterSize(words)),
		Precision:   cuckooV2.Low,
	})

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterLk(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooLk.NewFilter(uint(filterSize(words)))

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterFastrand(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooLocal.NewFilter(cuckooLocal.Config{
		NumElements: uint(filterSize(words)),
	})

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterVed(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooVed.NewFilter(uint32(filterSize(words)))

	insert := func(b []byte) bool { return cf.Insert(b) }
	contains := func(b []byte) bool { return cf.Lookup(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterLinSingle(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooLin.NewFilter(4, 16, uint(filterSize(words)), cuckooLin.TableTypeSingle)

	insert := func(b []byte) bool { return cf.Add(b) }
	contains := func(b []byte) bool { return cf.Contain(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testCuckoofilterLinPacked(words [][]byte) filterStats {
	memBefore := heapAllocs()
	cf := cuckooLin.NewFilter(4, 13, uint(filterSize(words)), cuckooLin.TableTypePacked)

	insert := func(b []byte) bool { return cf.Add(b) }
	contains := func(b []byte) bool { return cf.Contain(b) }
	return testImplementation(words, memBefore, insert, contains)
}

func testImplementation(words [][]byte, memBefore uint64,
	insert func([]byte) bool, contains func([]byte) bool) (stats filterStats) {
	skip := func(i, j int) bool { return (i+j)%200 == 0 }
	for i, w1 := range words {
		for j, w2 := range words[0:*wordListMultiplier] {
			if !skip(i, j) {
				w := bytes.Join([][]byte{w1, w2}, []byte{})
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
	remaining := make([][]byte, 0, len(words)/200)
	for i, w1 := range words {
		for j, w2 := range words[0:*wordListMultiplier] {
			if !skip(i, j) {
				continue
			}
			w := bytes.Join([][]byte{w1, w2}, []byte{})
			remaining = append(remaining, w)
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

func readWords() [][]byte {
	file, err := os.Open(*wordListPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var words [][]byte
	for scanner.Scan() {
		words = append(words, slices.Clone(scanner.Bytes()))
	}
	return words
}

func heapAllocs() uint64 {
	runtime.GC() // Run GC to clean up unreachable objects
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

func bloomHash(b []byte) hash.Hash64 {
	h := fnv.New64()
	h.Write(b)
	return h
}
