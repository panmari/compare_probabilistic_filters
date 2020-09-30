# Comparing probabilistic set membership datastructures

Most notable contenders in this category are

* Bloom filters
* Cuckoo filters

Both have various implementations readily available. Here we're mainly interested in the false positive rate acchieved. This is a function of the memory in use, but here I'm only testing the standard parameters provided by each library.

## Results

    steakknife/bloomfilter: fp=3, fp_rate=0.000651
    seiflotfy/cuckoofilter: fp=137, fp_rate=0.029724
