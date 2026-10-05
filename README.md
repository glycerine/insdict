insdict: deterministic insert-ordered iteration hash table for Golang
=======

Dict provides a hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.

Dict is faster than the built in Go map, especially on iteration (the most common
operation when using a Dict as a set). For full table scan, we are 10x faster.

# benchmarks of this Dict versus the built-in Go map (Map).

~~~
$ go test -v -run=xxx -bench=.
goos: darwin
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-8           	 4265073	       268.8 ns/op	        26.88 put_ns/key	     672 B/op	       4 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-8          	  284005	      4074 ns/op	        40.74 put_ns/key	   15840 B/op	      12 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-8         	   30645	     33157 ns/op	        33.16 put_ns/key	  130528 B/op	      18 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-8        	    3154	    379238 ns/op	        37.92 put_ns/key	  982497 B/op	      24 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-8       	     214	   5537346 ns/op	        55.37 put_ns/key	15416805 B/op	      32 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-8      	      12	  98541083 ns/op	        98.54 put_ns/key	123084265 B/op	      38 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-8            	 3188865	       363.4 ns/op	        36.34 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-8           	  273422	      4206 ns/op	        42.06 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-8          	   21042	     56543 ns/op	        56.54 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-8         	    2299	    512161 ns/op	        51.22 put_ns/key	  591483 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-8        	     220	   5337156 ns/op	        53.37 put_ns/key	 4729534 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-8       	      12	  95244660 ns/op	        95.24 put_ns/key	75460033 B/op	    8192 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-8    	 5177088	       225.4 ns/op	        22.54 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-8   	  772227	      1587 ns/op	        15.87 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-8  	   73930	     16368 ns/op	        16.37 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-8 	    6886	    176001 ns/op	        17.60 put_ns/key	  295552 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-8         	     421	   2865708 ns/op	        28.66 put_ns/key	 2364647 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-8        	      16	  67281714 ns/op	        67.28 put_ns/key	37832746 B/op	    4097 allocs/op
BenchmarkDictPutPresized
BenchmarkDictPutPresized/n=10
BenchmarkDictPutPresized/n=10-8            	 5140898	       236.0 ns/op	        23.60 put_ns/key	     528 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100
BenchmarkDictPutPresized/n=100-8           	  525598	      2700 ns/op	        27.00 put_ns/key	    8272 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000
BenchmarkDictPutPresized/n=1000-8          	   54482	     21229 ns/op	        21.23 put_ns/key	   65616 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=10000
BenchmarkDictPutPresized/n=10000-8         	    5688	    209084 ns/op	        20.91 put_ns/key	  483408 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100000
BenchmarkDictPutPresized/n=100000-8        	     475	   2591179 ns/op	        25.91 put_ns/key	 7692373 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000000
BenchmarkDictPutPresized/n=1000000-8       	      16	  65298355 ns/op	        65.30 put_ns/key	61522000 B/op	       3 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-8           	151176457	         7.913 ns/op	         7.913 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-8          	143957914	         8.234 ns/op	         8.234 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-8         	128587779	         9.155 ns/op	         9.155 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-8        	71275582	        15.96 ns/op	        15.96 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-8       	63347269	        18.26 ns/op	        18.26 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-8      	14016556	        86.04 ns/op	        86.04 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-8            	154188910	         7.707 ns/op	         7.707 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-8           	142141297	         7.886 ns/op	         7.886 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-8          	121396922	         9.277 ns/op	         9.277 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-8         	88077816	        13.11 ns/op	        13.11 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-8        	60847532	        19.30 ns/op	        19.30 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-8       	16962398	        67.58 ns/op	        67.58 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-8                    	156466657	         7.671 ns/op	         7.671 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-8                   	148983884	         8.174 ns/op	         8.174 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-8                  	136610372	         8.628 ns/op	         8.628 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-8                 	77886370	        15.17 ns/op	        15.17 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-8                	69101317	        17.78 ns/op	        17.78 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-8               	14427088	        83.04 ns/op	        83.04 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-8                     	177402972	         6.690 ns/op	         6.690 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-8                    	170530435	         7.301 ns/op	         7.301 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-8                   	149602292	         8.015 ns/op	         8.015 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-8                  	136648825	         8.733 ns/op	         8.733 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-8                 	87789901	        13.28 ns/op	        13.28 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-8                	22245355	        46.06 ns/op	        46.06 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-8                	207572511	         5.777 ns/op	         0.5777 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-8               	19052665	        62.38 ns/op	         0.6238 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-8              	 2276166	       566.5 ns/op	         0.5665 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-8             	  210451	      5413 ns/op	         0.5413 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-8            	   16990	     70401 ns/op	         0.7040 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-8           	     654	   1767863 ns/op	         1.768 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-8       	186027747	         6.481 ns/op	         0.6481 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-8      	13609515	        82.46 ns/op	         0.8246 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-8     	 1742098	       702.4 ns/op	         0.7024 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-8    	  162056	      7081 ns/op	         0.7081 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-8   	   13774	     88166 ns/op	         0.8817 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-8  	     523	   2322603 ns/op	         2.323 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-8                 	11323458	       111.8 ns/op	        11.18 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-8                	 1607286	       729.0 ns/op	         7.290 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-8               	  146666	      8010 ns/op	         8.010 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-8              	   14500	     83564 ns/op	         8.356 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-8             	    1586	    737025 ns/op	         7.370 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-8            	     100	  11136606 ns/op	        11.14 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	107.395s
~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
