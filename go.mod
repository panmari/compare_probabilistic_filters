module github.com/panmari/compare_probabilistic_filters

go 1.18

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96
	github.com/google/go-cmp v0.5.9
	github.com/irfansharif/cfilter v0.1.1
	github.com/linvon/cuckoo-filter v0.4.0
	github.com/livekit/cuckoofilter v1.1.0
	github.com/panmari/cuckoofilter v1.0.4-0.20220116144839-ac182fd3f9f3
	github.com/panmari/cuckoofilter_local v0.0.0
	github.com/seiflotfy/cuckoofilter v0.0.0-20220411075957-e3b120b3f5fb
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/vedhavyas/cuckoo-filter v1.6.2
	golang.org/x/exp v0.0.0-20220921164117-439092de6870
)

require (
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/zeebo/wyhash v0.0.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
)

replace github.com/panmari/cuckoofilter => github.com/panmari/cuckoofilter v1.0.4-0.20220924130152-b9b2432b5494
replace github.com/panmari/cuckoofilter_local => /home/smoser/cuckoofilter