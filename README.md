insdict: deterministic insert-ordered iteration hash table for Golang
=======

Dict provides a hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.

We now use int64 indexes so that Dict size is not limited to 2^31.

Dict with int32 indexes was faster than the built in Go map on all fronts.

With int64 indexes we are tied except on iteration (the most common
operation when using a Dict as a set). On iteration, also called a full 
table scan, we are still 10x faster than the built in Go map. 

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

# same benchmarks, but on linux

~~~
$ go test -v -run=xxx -bench=.
goos: linux
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-48            	 2058218	       587.6 ns/op	        58.76 put_ns/key	     672 B/op	       4 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-48           	  130516	      9219 ns/op	        92.19 put_ns/key	   15840 B/op	      12 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-48          	   15303	     76444 ns/op	        76.44 put_ns/key	  130528 B/op	      18 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-48         	    1497	    804537 ns/op	        80.45 put_ns/key	  982496 B/op	      24 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-48        	     112	  10024626 ns/op	       100.2 put_ns/key	15416804 B/op	      32 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-48       	       8	 138631261 ns/op	       138.6 put_ns/key	123084284 B/op	      38 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-48             	 1776561	       748.2 ns/op	        74.82 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-48            	  134764	      8329 ns/op	        83.29 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-48           	    8941	    123414 ns/op	       123.4 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-48          	    1093	   1082203 ns/op	       108.2 put_ns/key	  591483 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-48         	     100	  11510409 ns/op	       115.1 put_ns/key	 4729353 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-48        	       7	 177691655 ns/op	       177.7 put_ns/key	75476332 B/op	    8194 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-48     	 2489041	       478.9 ns/op	        47.89 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-48    	  361270	      3175 ns/op	        31.75 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-48   	   35690	     35706 ns/op	        35.71 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-48  	    3040	    388481 ns/op	        38.85 put_ns/key	  295552 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-48 	     207	   5548312 ns/op	        55.48 put_ns/key	 2364550 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-48         	       9	 116253361 ns/op	       116.3 put_ns/key	37832952 B/op	    4099 allocs/op
BenchmarkDictPutPresized
BenchmarkDictPutPresized/n=10
BenchmarkDictPutPresized/n=10-48             	 2393432	       497.5 ns/op	        49.75 put_ns/key	     528 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100
BenchmarkDictPutPresized/n=100-48            	  234302	      5116 ns/op	        51.16 put_ns/key	    8272 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000
BenchmarkDictPutPresized/n=1000-48           	   28227	     39621 ns/op	        39.62 put_ns/key	   65616 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=10000
BenchmarkDictPutPresized/n=10000-48          	    2728	    446897 ns/op	        44.69 put_ns/key	  483408 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=100000
BenchmarkDictPutPresized/n=100000-48         	     219	   5292380 ns/op	        52.92 put_ns/key	 7692370 B/op	       3 allocs/op
BenchmarkDictPutPresized/n=1000000
BenchmarkDictPutPresized/n=1000000-48        	      12	 100625582 ns/op	       100.6 put_ns/key	61522177 B/op	       4 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-48            	133164138	         8.948 ns/op	         8.948 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-48           	118515672	         9.928 ns/op	         9.928 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-48          	112551403	        10.55 ns/op	        10.55 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-48         	80782101	        15.08 ns/op	        15.08 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-48        	62836018	        18.78 ns/op	        18.78 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-48       	12226678	       102.7 ns/op	       102.7 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-48             	85168430	        12.35 ns/op	        12.35 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-48            	97880710	        13.44 ns/op	        13.44 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-48           	82813599	        14.76 ns/op	        14.76 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-48          	70141123	        16.53 ns/op	        16.53 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-48         	53616336	        24.24 ns/op	        24.24 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-48        	11967979	       105.1 ns/op	       105.1 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-48                     	137406115	         8.284 ns/op	         8.284 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-48                    	135717871	         8.670 ns/op	         8.670 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-48                   	115439620	        10.44 ns/op	        10.44 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-48                  	87206395	        13.55 ns/op	        13.55 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-48                 	70787023	        17.15 ns/op	        17.15 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-48                	13346227	        88.50 ns/op	        88.50 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-48                      	154370638	         7.836 ns/op	         7.836 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-48                     	139168990	         8.173 ns/op	         8.173 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-48                    	128744065	         9.233 ns/op	         9.233 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-48                   	119347287	         9.920 ns/op	         9.920 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-48                  	87308637	        14.10 ns/op	        14.10 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-48                 	22447066	        54.14 ns/op	        54.14 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-48                 	158656988	         7.474 ns/op	         0.7474 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-48                	15380718	        74.84 ns/op	         0.7484 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-48               	 1622821	       723.9 ns/op	         0.7239 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-48              	  165733	      6965 ns/op	         0.6965 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-48             	   16333	     72557 ns/op	         0.7256 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-48            	     724	   1669683 ns/op	         1.670 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-48        	143930596	         8.408 ns/op	         0.8408 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-48       	14979265	        81.08 ns/op	         0.8108 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-48      	 1614892	       732.2 ns/op	         0.7322 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-48     	  163437	      6542 ns/op	         0.6542 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-48    	   16822	     64835 ns/op	         0.6484 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-48   	     388	   3026659 ns/op	         3.027 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-48                  	10674592	       115.6 ns/op	        11.56 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-48                 	 1613738	       790.2 ns/op	         7.902 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-48                	  133922	      8919 ns/op	         8.919 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-48               	   14431	     81629 ns/op	         8.163 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-48              	    1686	    717282 ns/op	         7.173 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-48             	     100	  11376207 ns/op	        11.38 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	121.839s
~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
