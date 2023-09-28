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
steakknife/bloomfilter: size=41925250, mem=95.818 MB, insertFailed=0, fn=167684, fp=22, fp_rate=0.000105
AndreasBriese/bbloom: size=41925250, mem=127.998 MB, insertFailed=0, fn=167698, fp=8, fp_rate=0.000038
seiflotfy/cuckoofilter: size=41925250, mem=63.999 MB, insertFailed=0, fn=164500, fp=4021, fp_rate=0.019184
panmari/cuckoofilter: size=41925250, mem=127.998 MB, insertFailed=0, fn=167694, fp=16, fp_rate=0.000076
panmari/cuckoofilter/low: size=41925250, mem=64.000 MB, insertFailed=0, fn=164501, fp=4055, fp_rate=0.019346
livekit/cuckoofilter: size=41925250, mem=127.998 MB, insertFailed=0, fn=167683, fp=10, fp_rate=0.000048
vedhavyas/cuckoo-filter: size=41925250, mem=384.002 MB, insertFailed=0, fn=166158, fp=1986, fp_rate=0.009475
linvon/cuckoo-filter/single: size=41925250, mem=127.998 MB, insertFailed=0, fn=167690, fp=21, fp_rate=0.000100
linvon/cuckoo-filter/packed: size=41925250, mem=96.139 MB, insertFailed=0, fn=167600, fp=121, fp_rate=0.000577
```

## Runtime performance

All benchmarks are for filters with 50000 elements.

```
goos: linux
goarch: amd64
pkg: github.com/panmari/compare_probabilistic_filters
cpu: 12th Gen Intel(R) Core(TM) i7-1265U
BenchmarkFilters/size=50000/Insert/Bloomfilter-12                        5158182               207.2 ns/op
BenchmarkFilters/size=50000/Insert/BBloom-12                            29068202                39.98 ns/op
BenchmarkFilters/size=50000/Insert/SeiflotfyCuckoo-12                    2814543               424.9 ns/op
BenchmarkFilters/size=50000/Insert/PanmariCuckoo/Low-12                  2798888               439.3 ns/op
BenchmarkFilters/size=50000/Insert/PanmariCuckoo/Medium-12               2520884               449.9 ns/op
BenchmarkFilters/size=50000/Insert/LivekitCuckoo-12                      4466452               229.0 ns/op
BenchmarkFilters/size=50000/Insert/VedhavyasCuckoo-12                    1284637               993.5 ns/op
BenchmarkFilters/size=50000/Insert/LinCuckoo/single-12                 156161126                 7.360 ns/op
BenchmarkFilters/size=50000/Insert/LinCuckoo/packed-12                  94140253                12.31 ns/op
BenchmarkFilters/size=50000/ContainsTrue/Bloomfilter-12                  9112200               130.9 ns/op
BenchmarkFilters/size=50000/ContainsTrue/BBloom-12                      53032592                22.48 ns/op
BenchmarkFilters/size=50000/ContainsTrue/SeiflotfyCuckoo-12             57746794                20.80 ns/op
BenchmarkFilters/size=50000/ContainsTrue/PanmariCuckoo/Low-12           54694581                22.56 ns/op
BenchmarkFilters/size=50000/ContainsTrue/PanmariCuckoo/Medium-12        53476204                23.82 ns/op
BenchmarkFilters/size=50000/ContainsTrue/LivekitCuckoo-12              120490034                 9.982 ns/op
BenchmarkFilters/size=50000/ContainsTrue/VedhavyasCuckoo-12             10348413               122.8 ns/op
BenchmarkFilters/size=50000/ContainsTrue/LinCuckoo/single-12            42423476                28.94 ns/op
BenchmarkFilters/size=50000/ContainsTrue/LinCuckoo/packed-12            46023370                26.81 ns/op
BenchmarkFilters/size=50000/ContainsFalse/Bloomfilter-12                 9116546               134.2 ns/op
BenchmarkFilters/size=50000/ContainsFalse/BBloom-12                     44878254                26.90 ns/op
BenchmarkFilters/size=50000/ContainsFalse/SeiflotfyCuckoo-12            56785233                24.66 ns/op
BenchmarkFilters/size=50000/ContainsFalse/PanmariCuckoo/Low-12          45640098                27.04 ns/op
BenchmarkFilters/size=50000/ContainsFalse/PanmariCuckoo/Medium-12       44269923                25.85 ns/op
BenchmarkFilters/size=50000/ContainsFalse/LivekitCuckoo-12              95376951                13.03 ns/op
BenchmarkFilters/size=50000/ContainsFalse/VedhavyasCuckoo-12             8759331               135.2 ns/op
BenchmarkFilters/size=50000/ContainsFalse/LinCuckoo-12                  36252099                32.54 ns/op
BenchmarkFilters/size=50000/containsMixed/Bloomfilter-12                 9289278               134.9 ns/op
BenchmarkFilters/size=50000/containsMixed/BBloom-12                     46990591                27.80 ns/op
BenchmarkFilters/size=50000/containsMixed/SeiflotfyCuckoo-12            49834666                24.01 ns/op
BenchmarkFilters/size=50000/containsMixed/PanmariCuckoo/Low-12          47381558                26.40 ns/op
BenchmarkFilters/size=50000/containsMixed/PanmariCuckoo/Medium-12       40522986                26.58 ns/op
BenchmarkFilters/size=50000/containsMixed/LivekitCuckoo-12              96900874                12.65 ns/op
BenchmarkFilters/size=50000/containsMixed/VedhavyasCuckoo-12             9281860               132.7 ns/op
BenchmarkFilters/size=50000/containsMixed/LinCuckoo/single-12           35441432                43.59 ns/op
BenchmarkFilters/size=50000/containsMixed/LinCuckoo/packed-12           24818893                66.44 ns/op  
```
