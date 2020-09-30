package main

import (
	"bufio"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"os"

	"github.com/irfansharif/cfilter"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/steakknife/bloomfilter"
)

// Inserts the size of wordlist times this items into the filters.
const wordListMultiplier = 9

func main() {
	words := readWords()

	testCases := []struct {
		name       string
		testFilter func([]string) (int64, int64)
	}{
		{
			"steakknife/bloomfilter",
			testBloomfilter,
		},
		{
			"seiflotfy/cuckoofilter",
			testCuckoofilter,
		},
		// {
		// 	// panic: runtime error: index out of range
		// 	"irfansharif/cfilter",
		// 	testCfilter,
		// },
	}
	for _, tc := range testCases {
		fp, tn := tc.testFilter(words)
		fpRate := float64(fp) / (float64(fp + tn))
		// TODO(panmari): Also print allocated memory for better comparability
		fmt.Printf("%s: fp=%d, fp_rate=%f\n", tc.name, fp, fpRate)
	}

}

func testBloomfilter(words []string) (fp, tn int64) {
	bf, err := bloomfilter.NewOptimal(uint64(len(words)*wordListMultiplier), 0.001)
	if err != nil {
		log.Fatalf("failed creating bloom filter with size %d: %v", len(words), err)
	}

	insert := func(s string) { bf.Add(bloomHash(s)) }
	contains := func(s string) bool { return bf.Contains(bloomHash(s)) }
	return testImplementation(words, insert, contains)
}

func testCuckoofilter(words []string) (fp, tn int64) {
	cf := cuckoo.NewFilter(uint(len(words) * wordListMultiplier))

	insert := func(s string) { cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, insert, contains)
}

func testCfilter(words []string) (fp, tn int64) {
	cf := cfilter.New(cfilter.Size(uint(len(words) * wordListMultiplier)))

	insert := func(s string) { cf.Insert([]byte(s)) }
	contains := func(s string) bool { return cf.Lookup([]byte(s)) }
	return testImplementation(words, insert, contains)
}

func testImplementation(words []string, insert func(string), contains func(string) bool) (fp, tn int64) {
	remaining := make([]string, 0, len(words)/200)
	for i, w1 := range words {
		for j, w2 := range words[0:wordListMultiplier] {
			w := w1 + w2
			if (i+j)%200 == 0 {
				remaining = append(remaining, w)
				continue
			}
			insert(w)
		}
	}

	for _, w := range remaining {
		if contains(w) {
			fp++
		} else {
			tn++
		}
	}
	return fp, tn
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

func bloomHash(s string) hash.Hash64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h
}
