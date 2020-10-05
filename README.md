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

Time for 500 inserts.

```
BenchmarkInsertBloom-4                      7716            172887 ns/op
BenchmarkInsertCuckoo-4                     5728            209894 ns/op
BenchmarkInsertCuckooV2-4                  24951             51050 ns/op
```

### Lookup for a contained item

```
BenchmarkContainsTrueBloom-4                9163            140708 ns/op
BenchmarkContainsTrueCuckoo-4              48384             23680 ns/op
BenchmarkContainsTrueCuckooV2-4            33484             36829 ns/op
```

### Lookup for a missing item

```
BenchmarkContainsFalseBloom-4               8449            150120 ns/op
BenchmarkContainsFalseCuckoo-4             48962             24183 ns/op
BenchmarkContainsFalseCuckooV2-4           33981             34827 ns/op
```