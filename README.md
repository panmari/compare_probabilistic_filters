# Comparing probabilistic set membership datastructures

NOTE: There's [a writeup](https://panmari.github.io/2020/10/09/probabilistic-filter-golang.html) using the results generated from this repo.

Most notable contenders in this category are

* Bloom filters
* Cuckoo filters

Both have various implementations readily available.

## Comparison of false positive rate

False positive rate is a function of the memory in use. Some implementations
offer to configure a target false positive rate (e.g. `bloomfilter.NewOptimal(maxN uint64, p float64)`).

Results for `wordListMultiplier = 300`:

```bash
steakknife/bloomfilter: mem=64.926 MB fp=20, fp_rate=0.000130
AndreasBriese/bbloom: mem=63.999 MB fp=42, fp_rate=0.000273
seiflotfy/cuckoofilter: mem=31.999 MB fp=4326, fp_rate=0.028164
panmari/cuckoofilter: mem=61.679 MB fp=24, fp_rate=0.000156
```

## Runtime performance

### Insert

Time for constructing a filter with 500 elements.

```
BenchmarkInsertBloomFilter-4                        6692            169984 ns/op
BenchmarkInsertBBloom-4                            37878             30963 ns/op
BenchmarkInsertSeiflotfyCuckoo-4                   27488             43272 ns/op
BenchmarkInsertPanmariCuckoo-4                     61988             18910 ns/op
BenchmarkInsertVedhavyasCuckoo-4                    5964            177837 ns/op
```

### Lookup for a contained item

```
BenchmarkContainsTrueBloom-4                        7820            156934 ns/op
BenchmarkContainsTrueBBloom-4                      39914             29274 ns/op
BenchmarkContainsTrueSeiflotfyCuckoo-4             59832             19753 ns/op
BenchmarkContainsTruePanmariCuckoo-4               49832             23859 ns/op
BenchmarkContainsTrueVedhavyasCuckoo-4              8422            143366 ns/op
```

### Lookup for a missing item

```
BenchmarkContainsFalseBloom-4                       7700            157848 ns/op
BenchmarkContainsFalseBBloom-4                     42843             27156 ns/op
BenchmarkContainsFalseSeiflotfyCuckoo-4            54796             21104 ns/op
BenchmarkContainsFalsePanmariCuckoo-4              44841             26030 ns/op
BenchmarkContainsFalseVedhavyasCuckoo-4             7436            149577 ns/op
```