insdict: deterministic insert-ordered iteration hash table for Golang
=======

Dict provides a hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.

In terms of performance, the common full table scans over All() are 5-20x faster than the built-in
go map. Our point operations are tied or slightly faster than the built-in map. Benchmarks follow.

# Why our Dict range All() scans are so much faster than Go's built-in map.

Dict stores entries in one contiguous array, separate from its hash index.
All() walks that array in insertion order; Go inlines the iterator and its
callback. This gives predictable branches and sequential memory access that the
CPU can prefetch. The hash index is not touched during a scan.

Go's built-in map scans hash-table slots through a general runtime iterator.
Each returned entry involves a runtime call, iterator state updates, occupancy
checks, and checks for table growth and indirect key/value storage. It also
supports adding entries during iteration, which Dict.All() needs special
care to do correctly (see [the All docs](https://pkg.go.dev/github.com/glycerine/insdict#Dict.All); 
call Pack(true) first; or use [SlowWriteAll](https://pkg.go.dev/github.com/glycerine/insdict#Dict.SlowWriteAll)).

On Go 1.26.4/linux/amd64, a matched integer-value summation benchmark measured
about 112 instructions and 25 branches per map entry versus 12 instructions and
3 branches for Dict. With 1,000 entries fitting in cache, the map took
8.89 ns/key versus 0.49 ns/key for Dict. At one million entries it took
10.32 ns/key versus 2.03 ns/key. These measurements point mainly to iterator
overhead, with locality also contributing at larger sizes. The ratio depends
on table size, key/value types, and how much work the loop body does.

See [scan_profile.txt](scan_profile.txt) for Linux `perf` counters, CPU profiles,
and reproduction commands, and [scan_bench_test.go](scan_bench_test.go) for the
matched benchmarks. The
[Go runtime's iteration discussion](https://github.com/golang/go/blob/go1.26.4/src/internal/runtime/maps/map.go#L134-L175)
explains its more demanding mutation semantics.

# Dict memory overhead is about 1.34x to 1.7x

For 1,000,000 map[int]int entries vs Dict, heap after GC:

|              | insdict   | Go map    | insdict overhead |
| -------------|   -------:|    ------:| ----------------:|
| Grown        | 60.67 MiB | 36.04 MiB | 68%              |
| Presized     | 48.52 MiB | 36.08 MiB | 34%              |

See memory_test.go to evaluate for your data shape and size.

# benchmarks of insdict (Dict) versus the built-in Go map (Map).

~~~
$ go test -v -run=xxx -bench=.

goos: linux
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-48            	 2142226	       580.4 ns/op	        58.04 put_ns/key	     704 B/op	       4 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-48           	  132698	      9093 ns/op	        90.93 put_ns/key	   16352 B/op	      12 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-48          	   16454	     72757 ns/op	        72.76 put_ns/key	  135136 B/op	      18 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-48         	    1579	    729345 ns/op	        72.93 put_ns/key	 1019873 B/op	      24 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-48        	     126	   8959778 ns/op	        89.60 put_ns/key	15945696 B/op	      32 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-48       	      12	  86388814 ns/op	        86.39 put_ns/key	127283373 B/op	      39 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-48             	 1570562	       749.8 ns/op	        74.98 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-48            	  144189	      8655 ns/op	        86.55 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-48           	    9844	    121640 ns/op	       121.6 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-48          	    1082	   1110290 ns/op	       111.0 put_ns/key	  591483 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-48         	     100	  11964505 ns/op	       119.6 put_ns/key	 4729750 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-48        	       6	 187960179 ns/op	       188.0 put_ns/key	75475466 B/op	    8194 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-48     	 2527846	       484.7 ns/op	        48.47 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-48    	  370183	      3301 ns/op	        33.01 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-48   	   34626	     34640 ns/op	        34.64 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-48  	    3068	    384561 ns/op	        38.46 put_ns/key	  295552 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-48 	     205	   5819163 ns/op	        58.19 put_ns/key	 2364741 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-48         	       9	 113298323 ns/op	       113.3 put_ns/key	37832952 B/op	    4099 allocs/op
BenchmarkDictPutPresized
BenchmarkDictPutPresized/n=10
BenchmarkDictPutPresized/n=10-48             	 3127922	       378.8 ns/op	        37.88 put_ns/key	     464 B/op	       2 allocs/op
BenchmarkDictPutPresized/n=100
BenchmarkDictPutPresized/n=100-48            	  337417	      3550 ns/op	        35.50 put_ns/key	    5504 B/op	       2 allocs/op
BenchmarkDictPutPresized/n=1000
BenchmarkDictPutPresized/n=1000-48           	   36391	     33376 ns/op	        33.38 put_ns/key	   51200 B/op	       2 allocs/op
BenchmarkDictPutPresized/n=10000
BenchmarkDictPutPresized/n=10000-48          	    3205	    355963 ns/op	        35.60 put_ns/key	  475136 B/op	       2 allocs/op
BenchmarkDictPutPresized/n=100000
BenchmarkDictPutPresized/n=100000-48         	     334	   4359064 ns/op	        43.59 put_ns/key	 5562369 B/op	       2 allocs/op
BenchmarkDictPutPresized/n=1000000
BenchmarkDictPutPresized/n=1000000-48        	      19	  58475076 ns/op	        58.47 put_ns/key	50880700 B/op	       3 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-48            	184086732	         6.429 ns/op	         6.429 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-48           	165364309	         7.184 ns/op	         7.184 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-48          	150305198	         7.841 ns/op	         7.841 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-48         	110955352	        10.92 ns/op	        10.92 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-48        	90565683	        13.35 ns/op	        13.35 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-48       	19629327	        59.49 ns/op	        59.49 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-48             	95712310	        12.27 ns/op	        12.27 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-48            	97365656	        12.69 ns/op	        12.69 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-48           	91365162	        14.42 ns/op	        14.42 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-48          	71936695	        16.70 ns/op	        16.70 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-48         	50704756	        23.13 ns/op	        23.13 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-48        	11770156	       103.7 ns/op	       103.7 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-48                     	189820318	         6.122 ns/op	         6.122 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-48                    	190962265	         6.314 ns/op	         6.314 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-48                   	170471432	         6.816 ns/op	         6.816 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-48                  	117700486	        10.10 ns/op	        10.10 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-48                 	102885072	        11.70 ns/op	        11.70 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-48                	23650609	        50.63 ns/op	        50.63 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-48                      	149736351	         8.054 ns/op	         8.054 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-48                     	138917115	         8.072 ns/op	         8.072 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-48                    	130131655	         9.399 ns/op	         9.399 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-48                   	115108560	        10.07 ns/op	        10.07 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-48                  	84979269	        14.01 ns/op	        14.01 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-48                 	18817881	        66.53 ns/op	        66.53 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-48                 	208511940	         5.645 ns/op	         0.5645 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-48                	22164832	        58.82 ns/op	         0.5882 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-48               	 1993611	       507.7 ns/op	         0.5077 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-48              	  213732	      5803 ns/op	         0.5803 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-48             	   19012	     61477 ns/op	         0.6148 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-48            	     678	   1838104 ns/op	         1.838 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-48        	123750986	         8.969 ns/op	         0.8969 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-48       	13745991	        86.24 ns/op	         0.8624 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-48      	 1426533	       811.2 ns/op	         0.8112 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-48     	  147583	      8022 ns/op	         0.8022 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-48    	   14822	     80281 ns/op	         0.8028 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-48   	     514	   2355491 ns/op	         2.355 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-48                  	10896586	       115.0 ns/op	        11.50 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-48                 	 1647214	       742.6 ns/op	         7.426 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-48                	  130939	      9201 ns/op	         9.201 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-48               	   14010	     83011 ns/op	         8.301 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-48              	    1719	    714405 ns/op	         7.144 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-48             	     100	  10383562 ns/op	        10.38 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	118.009s
~~~

# (optional) Sorted iteration

Iterating keys in sorted order is lazy. If you do not use it, you do
not pay for extra memory that the sorted order index consumes. It
is merely a convenience for those times when you do want to 
traverse keys in sorted order (and possibly Del keys <= current key
while doing so -- this is a common use case for clearing out all old 
timestamps for me).

`Ascend()` scans all keys in ascending order; `Ascend(lo)` starts at `lo`,
and `Ascend(lo, hi)` scans `lo <= key < hi`.

`Descend()` scans all keys in descending order; `Descend(hi)` visits keys
`<= hi`, and `Descend(hi, lo)` scans `hi >= key > lo`. 

Both share the same sorted index, which is allocated on first
iteration and rebuilt lazily after inserts, deletes, or compaction.

Built-in ordered key types use typed comparators without reflection. NaNs sort last
during Ascend, first on Descend.

You can supply `EasyCompare[K]` for named ordered types or a custom comparator:

```go
d := insdict.NewDictFunc[int, string](insdict.EasyHashInt, insdict.EasyCompare[int])
d.Put(3, "three")
for key, value := range d.Ascend(0, 10) {
    fmt.Println(key, value)
}
```

`NewDictFuncSize(hash, hint, compare)` also accepts the optional comparator.

Ascend and Descend require exclusive access. During Ascend, `Del` of keys at or
below the current key is supported; during Descend, keys at or above the current
key may be deleted. Insertions, compaction, and clearing during either scan are
not supported. 

Benchmarks show we are pretty fast. 2x faster Put than an in-memory B-tree.
Sorted-order full table read (scan) time is expensive the first time, of course, but
is cheap if writes are infrequent compared to scans: the first sort basically costs
the same (303.54 ns/key) as inserting into a B-tree the first time (316.6 ns/key). 

| Operation (showing ns/key) |  insdict.Dict |BP-Tree | builtin Go map | tidwall/btree | red-black tree |
| -------------------------- |  -----------: | -----: | -------------: | ------------: | -------------: |
| Get                        |          25.5 |  115.2 |           19.6 |         133.1 |          232.6 |
| Put                        |         150.1 |  201.9 |          168.6 |         316.6 |          635.6 |
| Ordered scan (amortized)   |          7.46 |   3.57 |  not supported |          4.84 |          18.05 |
| First ordered scan         |        303.54 |  10.12 |  not supported |          7.75 |          32.03 |

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
