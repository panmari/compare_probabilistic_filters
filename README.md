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
InsertBloomFilter-4              165µs ± 2%
InsertBBloom-4                  30.9µs ± 0%
InsertSeiflotfyCuckoo-4         43.0µs ± 0%
InsertPanmariCuckoo-4           18.9µs ± 1%
InsertVedhavyasCuckoo-4          176µs ± 0%
```

### Lookup for a contained item

```
ContainsTrueBloom-4              150µs ± 1%
ContainsTrueBBloom-4            29.0µs ± 1%
ContainsTrueSeiflotfyCuckoo-4   19.9µs ± 2%
ContainsTruePanmariCuckoo-4     16.7µs ± 0%
ContainsTrueVedhavyasCuckoo-4    143µs ± 3%
```

### Lookup for a missing item

```
ContainsFalseBloom-4             152µs ± 1%
ContainsFalseBBloom-4           26.8µs ± 0%
ContainsFalseSeiflotfyCuckoo-4  21.0µs ± 0%
ContainsFalsePanmariCuckoo-4    24.7µs ± 0%
ContainsFalseVedhavyasCuckoo-4   148µs ± 2%
```