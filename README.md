# Comparing probabilistic set membership datastructures

NOTE: There's [a writeup](https://panmari.github.io/2020/10/09/probabilistic-filter-golang.html) using the results generated from this repo.

Most notable contenders in this category are

* Bloom filters
* Cuckoo filters

Both have various implementations readily available.

## Comparison of false positive rate

False positive rate is a function of the memory in use. Some implementations
offer to configure a target false positive rate (e.g. `bloomfilter.NewOptimal(maxN uint64, p float64)`).

Results for `wordListMultiplier = 250`:

```bash
steakknife/bloomfilter: size=25796274, mem=58.958 MB, insertFailed=0, fn=0, fp=10, fp_rate=0.000078
AndreasBriese/bbloom: size=25796274, mem=63.999 MB, insertFailed=0, fn=0, fp=9, fp_rate=0.000070
seiflotfy/cuckoofilter: size=25796274, mem=31.999 MB, insertFailed=0, fn=0, fp=2954, fp_rate=0.022998
panmari/cuckoofilter: size=25796274, mem=63.999 MB, insertFailed=0, fn=0, fp=10, fp_rate=0.000078
panmari/cuckoofilter/low: size=25796274, mem=31.999 MB, insertFailed=0, fn=0, fp=3053, fp_rate=0.023768
vedhavyas/cuckoo-filter: size=25796274, mem=191.999 MB, insertFailed=0, fn=0, fp=783, fp_rate=0.006096
linvon/cuckoo-filter/single: size=25796274, mem=63.999 MB, insertFailed=0, fn=0, fp=7, fp_rate=0.000054
linvon/cuckoo-filter/packed: size=25796274, mem=48.140 MB, insertFailed=0, fn=0, fp=92, fp_rate=0.000716
```

## Runtime performance

All benchmarks are for filters with 500 elements.

```
Filters/size=500/Insert/Bloomfilter-4                348ns ± 5%                                                                                                                      
Filters/size=500/Insert/BBloom-4                    54.2ns ± 1%                                                                                                                      
Filters/size=500/Insert/SeiflotfyCuckoo-4           67.6ns ± 1%                                                                                                                      
Filters/size=500/Insert/PanmariCuckoo-4             77.0ns ± 6%                                                                                                                      
Filters/size=500/Insert/VedhavyasCuckoo-4            362ns ± 5%                                                                                                                      
Filters/size=500/Insert/LinCuckoo-4                 77.2ns ± 8%                                                                                                                      
Filters/size=500/ContainsTrue/Bloomfilter-4          393ns ±21%                                                                                                                      
Filters/size=500/ContainsTrue/BBloom-4              47.9ns ± 1%                                                                                                                      
Filters/size=500/ContainsTrue/SeiflotfyCuckoo-4     26.8ns ± 1%                                                                                                                      
Filters/size=500/ContainsTrue/PanmariCuckoo-4       74.5ns ± 8%                                                                                                                      
Filters/size=500/ContainsTrue/VedhavyasCuckoo-4      276ns ± 4%                                                                                                                      
Filters/size=500/ContainsTrue/LinCuckoo-4           50.3ns ± 1%                                                                                                                      
Filters/size=500/ContainsFalse/Bloomfilter-4         359ns ±13%                                                                                                                      
Filters/size=500/ContainsFalse/BBloom-4             48.7ns ± 2%                                                                                                                      
Filters/size=500/ContainsFalse/SeiflotfyCuckoo-4    31.7ns ± 5%                                                                                                                      
Filters/size=500/ContainsFalse/PanmariCuckoo-4      74.8ns ± 8%                                                                                                                      
Filters/size=500/ContainsFalse/VedhavyasCuckoo-4     286ns ± 1%                                                                                                                      
Filters/size=500/ContainsFalse/LinCuckoo-4          72.0ns ± 6%                                                                                                                      
Filters/size=500/containsMixed/Bloomfilter-4         370ns ±16%                                                                                                                      
Filters/size=500/containsMixed/BBloom-4             51.9ns ± 8%                                                                                                                      
Filters/size=500/containsMixed/SeiflotfyCuckoo-4    31.1ns ± 8%                                                                                                                      
Filters/size=500/containsMixed/PanmariCuckoo-4      77.0ns ± 8%                                                                                                                      
Filters/size=500/containsMixed/VedhavyasCuckoo-4     284ns ± 4%                                                                                                                      
Filters/size=500/containsMixed/LinCuckoo-4          63.4ns ± 4%                                                                                                                      
```
