// Benchmark for thread-unsafe interactions with probabilistic filters.
// Note that github.com/steakknife/bloomfilter doesn't allow interacting
// in a thread-unsafe way, leading to higher numbers there.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/AndreasBriese/bbloom"
	cuckooLin "github.com/linvon/cuckoo-filter"
	cuckooV2 "github.com/panmari/cuckoofilter"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/steakknife/bloomfilter"
	cuckooVed "github.com/vedhavyas/cuckoo-filter"
)

var (
	words      [][]byte // Words contained in filter.
	otherWords [][]byte // Words NOT contained in filter.
	mixedWords [][]byte // Mix of words that are contained/not contained.
	numWords   = 500
)

const maxNumWords = 50000

func initWords() error {
	fd, err := os.Open(*wordListPath)
	if err != nil {
		return fmt.Errorf("Failed reading %q: %v", *wordListPath, err)
	}
	scanner := bufio.NewScanner(fd)
	for i := 0; i < maxNumWords && scanner.Scan(); i++ {
		words = append(words, scanner.Bytes())
	}
	for i := 0; i < maxNumWords && scanner.Scan(); i++ {
		otherWords = append(otherWords, scanner.Bytes())
	}
	if readWords := len(words) + len(otherWords); readWords < 2*maxNumWords {
		return fmt.Errorf("Not enough words in wordlist, want at least %d words, got %d\n", 2*maxNumWords, readWords)
	}
	r := rand.New(rand.NewSource(0))
	wordsIndex := 0
	otherWordsIndex := 0
	for i := 0; i < maxNumWords; i++ {
		if r.Intn(2) == 0 {
			mixedWords = append(mixedWords, words[wordsIndex])
			wordsIndex++
			continue
		}
		mixedWords = append(mixedWords, otherWords[otherWordsIndex])
		otherWordsIndex++
	}
	return nil
}

func BenchmarkFilters(b *testing.B) {
	if err := initWords(); err != nil {
		b.Fatalf("Failed initializing words: %v", err)
	}
	for _, n := range []int{500, 5000, 50000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			if n > maxNumWords {
				b.Fatalf("Num words too large: %d > %d", n, maxNumWords)
			}
			numWords = n
			b.Run("Insert", insert)
			b.Run("ContainsTrue", containsTrue)
			b.Run("ContainsFalse", containsFalse)
			b.Run("containsMixed", containsMixed)
		})
	}
}

func insert(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		for i := 0; i < b.N; {
			f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
			for _, w := range words[:numWords] {
				f.Add(bloomHash(w))
			}
			i += numWords
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		for i := 0; i < b.N; {
			f := bbloom.New(float64(numWords), 0.002)
			for _, w := range words[:numWords] {
				f.Add(w)
			}
			i += numWords
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; {
			f := cuckoo.NewFilter(uint(numWords))
			for _, w := range words[:numWords] {
				f.Insert(w)
			}
			i += numWords
		}
	})
	for _, precision := range []cuckooV2.FilterPrecision{cuckooV2.Low, cuckooV2.Medium} {
		b.Run("PanmariCuckoo/"+fmt.Sprint(precision), func(b *testing.B) {
			for i := 0; i < b.N; {
				f := cuckooV2.NewFilter(cuckooV2.Config{NumElements: uint(numWords), Precision: precision})
				for _, w := range words[:numWords] {
					f.Insert(w)
				}
				i += numWords
			}
		})
	}
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		for i := 0; i < b.N; {
			f := cuckooVed.NewFilter(uint32(numWords))
			for _, w := range words[:numWords] {
				f.Insert(w)
			}
			i += numWords
		}
	})
	b.Run("LinCuckoo", func(b *testing.B) {
		b.Run("single", func(b *testing.B) {
			for i := 0; i < b.N; {
				f := cuckooLin.NewFilter(4, 16, uint(numWords), cuckooLin.TableTypeSingle)
				for _, w := range words[:numWords] {
					f.Add(w)
				}
				i += numWords
			}
		})
		b.Run("packed", func(b *testing.B) {
			for i := 0; i < b.N; {
				f := cuckooLin.NewFilter(4, 13, uint(numWords), cuckooLin.TableTypePacked)
				for _, w := range words[:numWords] {
					f.Add(w)
				}
				i += numWords
			}
		})
	})
}

func containsTrue(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words[:numWords] {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range words[:numWords] {
				f.Contains(bloomHash(w))
			}
			i += numWords
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words[:numWords] {
			f.Add(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range words[:numWords] {
				f.Has(w)
			}
			i += numWords
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range words[:numWords] {
				runtime.KeepAlive(f.Lookup(w))
			}
			i += numWords
		}
	})
	for _, precision := range []cuckooV2.FilterPrecision{cuckooV2.Low, cuckooV2.Medium} {
		b.Run("PanmariCuckoo/"+fmt.Sprint(precision), func(b *testing.B) {
			f := cuckooV2.NewFilter(cuckooV2.Config{NumElements: uint(numWords), Precision: precision})
			for _, w := range words[:numWords] {
				f.Insert(w)
			}
			if err := os.WriteFile("/tmp/"+strings.ReplaceAll(b.Name(), "/", "_"), f.Encode(), 0644); err != nil {
				b.Error(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; {
				for _, w := range words[:numWords] {
					f.Lookup(w)
				}
				i += numWords
			}
		})
	}
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		f := cuckooVed.NewFilter(uint32(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range words[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("LinCuckoo", func(b *testing.B) {
		b.Run("single", func(b *testing.B) {
			f := cuckooLin.NewFilter(4, 16, uint(numWords), cuckooLin.TableTypeSingle)
			for _, w := range words[:numWords] {
				f.Add(w)
			}
			b.ResetTimer()
			for i := 0; i < b.N; {
				for _, w := range words[:numWords] {
					f.Contain(w)
				}
				i += numWords
			}
		})
		b.Run("packed", func(b *testing.B) {
			f := cuckooLin.NewFilter(4, 13, uint(numWords), cuckooLin.TableTypePacked)
			for _, w := range words[:numWords] {
				f.Add(w)
			}
			b.ResetTimer()
			for i := 0; i < b.N; {
				for _, w := range words[:numWords] {
					f.Contain(w)
				}
				i += numWords
			}
		})
	})
}

func containsFalse(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words[:numWords] {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Contains(bloomHash(w))
			}
			i += numWords
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words[:numWords] {
			f.Add(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Has(w)
			}
			i += numWords
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		f := cuckooV2.NewFilter(cuckooV2.Config{NumElements: uint(numWords), Precision: cuckooV2.Low})
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		f := cuckooVed.NewFilter(uint32(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("LinCuckoo", func(b *testing.B) {
		f := cuckooLin.NewFilter(4, 16, uint(numWords), cuckooLin.TableTypeSingle)
		for _, w := range words[:numWords] {
			f.Add(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range otherWords[:numWords] {
				f.Contain(w)
			}
			i += numWords
		}
	})
}

func containsMixed(b *testing.B) {
	b.Run("Bloomfilter", func(b *testing.B) {
		f, _ := bloomfilter.NewOptimal(uint64(numWords), 0.0001)
		for _, w := range words[:numWords] {
			f.Add(bloomHash(w))
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range mixedWords[:numWords] {
				f.Contains(bloomHash(w))
			}
			i += numWords
		}
	})
	b.Run("BBloom", func(b *testing.B) {
		f := bbloom.New(float64(numWords), 0.002)
		for _, w := range words[:numWords] {
			f.Add(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range mixedWords[:numWords] {
				f.Has(w)
			}
			i += numWords
		}
	})
	b.Run("SeiflotfyCuckoo", func(b *testing.B) {
		f := cuckoo.NewFilter(uint(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range mixedWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("PanmariCuckoo", func(b *testing.B) {
		f := cuckooV2.NewFilter(cuckooV2.Config{NumElements: uint(numWords), Precision: cuckooV2.Low})
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range mixedWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("VedhavyasCuckoo", func(b *testing.B) {
		f := cuckooVed.NewFilter(uint32(numWords))
		for _, w := range words[:numWords] {
			f.Insert(w)
		}
		b.ResetTimer()
		for i := 0; i < b.N; {
			for _, w := range mixedWords[:numWords] {
				f.Lookup(w)
			}
			i += numWords
		}
	})
	b.Run("LinCuckoo", func(b *testing.B) {
		b.Run("single", func(b *testing.B) {
			f := cuckooLin.NewFilter(4, 16, uint(numWords), cuckooLin.TableTypeSingle)
			for _, w := range words[:numWords] {
				f.Add(w)
			}
			b.ResetTimer()
			for i := 0; i < b.N; {
				for _, w := range mixedWords[:numWords] {
					f.Contain(w)
				}
				i += numWords
			}
		})
		b.Run("packed", func(b *testing.B) {
			f := cuckooLin.NewFilter(4, 13, uint(numWords), cuckooLin.TableTypePacked)
			for _, w := range words[:numWords] {
				f.Add(w)
			}
			b.ResetTimer()
			for i := 0; i < b.N; {
				for _, w := range mixedWords[:numWords] {
					f.Contain(w)
				}
				i += numWords
			}
		})
	})
}
