insdict: deterministic insert-ordered iteration hash table for Golang
=======

A hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.


# benchmarks versus built-in Go map.

~~~
$ go test -v -run=xxx -bench=.
goos: darwin
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-8           	 4236652	       269.1 ns/op	        26.91 put_ns/key	     576 B/op	       4 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-8          	  317474	      3744 ns/op	        37.44 put_ns/key	   13824 B/op	      12 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-8         	   37797	     29615 ns/op	        29.61 put_ns/key	  114176 B/op	      18 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-8        	    3422	    346873 ns/op	        34.69 put_ns/key	  851456 B/op	      24 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-8       	     244	   5298796 ns/op	        52.99 put_ns/key	13319685 B/op	      32 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-8      	      13	  91084816 ns/op	        91.08 put_ns/key	106307080 B/op	      38 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-8            	 3293583	       358.0 ns/op	        35.80 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-8           	  273104	      4116 ns/op	        41.16 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-8          	   20926	     56932 ns/op	        56.93 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-8         	    2217	    507611 ns/op	        50.76 put_ns/key	  591482 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-8        	     213	   5352151 ns/op	        53.52 put_ns/key	 4729548 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-8       	       9	 124655296 ns/op	       124.7 put_ns/key	75497990 B/op	    8196 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-8    	 4497073	       233.9 ns/op	        23.39 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-8   	  746996	      1611 ns/op	        16.11 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-8  	   65760	     16210 ns/op	        16.21 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-8 	    6532	    195156 ns/op	        19.52 put_ns/key	  295552 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-8         	     398	   3057381 ns/op	        30.57 put_ns/key	 2364750 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-8        	      15	  69897488 ns/op	        69.90 put_ns/key	37832711 B/op	    4097 allocs/op
BenchmarkDictPutPresized
BenchmarkDictPutPresized/n=10
BenchmarkDictPutPresized/n=10-8            	 5080443	       231.0 ns/op	        23.10 put_ns/key	     464 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100
BenchmarkDictPutPresized/n=100-8           	  536473	      2197 ns/op	        21.97 put_ns/key	    7248 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000
BenchmarkDictPutPresized/n=1000-8          	   68350	     19725 ns/op	        19.73 put_ns/key	   57424 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=10000
BenchmarkDictPutPresized/n=10000-8         	    5611	    224293 ns/op	        22.43 put_ns/key	  417872 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100000
BenchmarkDictPutPresized/n=100000-8        	     451	   2368295 ns/op	        23.68 put_ns/key	 6643792 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000000
BenchmarkDictPutPresized/n=1000000-8       	      22	  50979305 ns/op	        50.98 put_ns/key	53133392 B/op	       3 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-8           	149597185	         8.622 ns/op	         8.622 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-8          	126343340	         9.557 ns/op	         9.557 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-8         	89518482	        12.39 ns/op	        12.39 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-8        	69159433	        18.88 ns/op	        18.88 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-8       	65088783	        21.43 ns/op	        21.43 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-8      	16933659	        78.75 ns/op	        78.75 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-8            	145771962	         8.977 ns/op	         8.977 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-8           	138360951	         8.399 ns/op	         8.399 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-8          	123762282	         9.617 ns/op	         9.617 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-8         	82939708	        13.07 ns/op	        13.07 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-8        	62825701	        20.52 ns/op	        20.52 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-8       	18313032	        65.03 ns/op	        65.03 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-8                    	152804330	         7.861 ns/op	         7.861 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-8                   	144720765	         8.840 ns/op	         8.840 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-8                  	124300645	         9.551 ns/op	         9.551 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-8                 	72909276	        15.72 ns/op	        15.72 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-8                	68777526	        17.15 ns/op	        17.15 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-8               	17884054	        68.42 ns/op	        68.42 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-8                     	177290838	         7.119 ns/op	         7.119 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-8                    	154808145	         8.171 ns/op	         8.171 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-8                   	137821508	         8.612 ns/op	         8.612 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-8                  	100000000	        11.18 ns/op	        11.18 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-8                 	84324668	        13.89 ns/op	        13.89 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-8                	26665224	        46.08 ns/op	        46.08 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-8                	199423839	         6.567 ns/op	         0.6567 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-8               	16416741	        67.10 ns/op	         0.6710 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-8              	 1988001	       587.1 ns/op	         0.5871 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-8             	  209703	      5806 ns/op	         0.5806 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-8            	   18228	     67627 ns/op	         0.6763 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-8           	     698	   1726096 ns/op	         1.726 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-8       	186504960	         6.685 ns/op	         0.6685 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-8      	13047747	        88.79 ns/op	         0.8879 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-8     	 1598455	       755.9 ns/op	         0.7559 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-8    	  160752	      7506 ns/op	         0.7506 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-8   	   12213	     90922 ns/op	         0.9092 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-8  	     480	   2378436 ns/op	         2.378 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-8                 	11100625	       108.9 ns/op	        10.89 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-8                	 1640890	       697.6 ns/op	         6.976 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-8               	  141837	      8329 ns/op	         8.329 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-8              	   13678	     89862 ns/op	         8.986 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-8             	    1525	    770410 ns/op	         7.704 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-8            	     100	  10766240 ns/op	        10.77 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	109.440s
~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
