# Comparing probabilistic set membership datastructures

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
BenchmarkInsertBloomFilter-4                8884            134940 ns/op
BenchmarkInsertBBloom-4                    37064             31784 ns/op
BenchmarkInsertCuckoo-4                     5706            200611 ns/op
BenchmarkInsertCuckooV2-4                  25261             47125 ns/op
```

### Lookup for a contained item

```
BenchmarkContainsTrueBloom-4                8960            125152 ns/op
BenchmarkContainsTrueBBloom-4              38833             31392 ns/op
BenchmarkContainsTrueCuckoo-4              50750             23198 ns/op
BenchmarkContainsTrueCuckooV2-4            34770             34430 ns/op
```

### Lookup for a missing item

```
BenchmarkContainsFalseBloom-4               9186            127424 ns/op
BenchmarkContainsFalseBBloom-4             41846             27788 ns/op
BenchmarkContainsFalseCuckoo-4             50150             23507 ns/op
BenchmarkContainsFalseCuckooV2-4           35347             33960 ns/op
```