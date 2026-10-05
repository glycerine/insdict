insdict: deterministic insert-ordered iteration hash table for Golang
=======

A hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.


# benchmarks

~~~
go test -v -run=xxx -bench=.

goos: darwin
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-8           	 4409697	       270.9 ns/op	        27.09 put_ns/key	     576 B/op	       4 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-8          	  272703	      3906 ns/op	        39.06 put_ns/key	   13824 B/op	      12 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-8         	   39198	     29741 ns/op	        29.74 put_ns/key	  114176 B/op	      18 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-8        	    3501	    342937 ns/op	        34.29 put_ns/key	  851456 B/op	      24 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-8       	     241	   4935478 ns/op	        49.35 put_ns/key	13319685 B/op	      32 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-8      	      14	  77689064 ns/op	        77.69 put_ns/key	106307072 B/op	      38 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-8            	 3336882	       358.5 ns/op	        35.85 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-8           	  285188	      4003 ns/op	        40.03 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-8          	   21308	     55374 ns/op	        55.37 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-8         	    2336	    507079 ns/op	        50.71 put_ns/key	  591481 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-8        	     226	   5296537 ns/op	        52.97 put_ns/key	 4729349 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-8       	      12	  96560114 ns/op	        96.56 put_ns/key	75503097 B/op	    8196 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-8    	 5233291	       226.2 ns/op	        22.62 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-8   	  769030	      1554 ns/op	        15.54 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-8  	   64756	     17714 ns/op	        17.71 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-8 	    7084	    172923 ns/op	        17.29 put_ns/key	  295553 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-8         	     432	   2712740 ns/op	        27.13 put_ns/key	 2364737 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-8        	      18	  66453186 ns/op	        66.45 put_ns/key	37832741 B/op	    4097 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-8           	153849123	         7.718 ns/op	         7.718 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-8          	147025454	         8.178 ns/op	         8.178 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-8         	132315062	         9.027 ns/op	         9.027 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-8        	78457596	        15.49 ns/op	        15.49 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-8       	70082877	        17.36 ns/op	        17.36 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-8      	18502413	        65.29 ns/op	        65.29 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-8            	160886718	         7.342 ns/op	         7.342 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-8           	133408916	         7.905 ns/op	         7.905 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-8          	131697872	         9.406 ns/op	         9.406 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-8         	96409647	        13.26 ns/op	        13.26 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-8        	61332528	        19.27 ns/op	        19.27 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-8       	18802016	        63.35 ns/op	        63.35 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-8                    	152684310	         7.654 ns/op	         7.654 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-8                   	145802900	         7.997 ns/op	         7.997 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-8                  	139941637	         8.502 ns/op	         8.502 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-8                 	77532310	        14.83 ns/op	        14.83 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-8                	70864642	        16.81 ns/op	        16.81 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-8               	18016556	        69.93 ns/op	        69.93 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-8                     	174682490	         6.942 ns/op	         6.942 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-8                    	157363578	         7.050 ns/op	         7.050 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-8                   	147103096	         8.239 ns/op	         8.239 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-8                  	135085810	         8.939 ns/op	         8.939 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-8                 	74789551	        14.03 ns/op	        14.03 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-8                	27353666	        45.00 ns/op	        45.00 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-8                	195833452	         6.022 ns/op	         0.6022 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-8               	17983996	        64.51 ns/op	         0.6451 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-8              	 2159012	       557.8 ns/op	         0.5578 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-8             	  230773	      5412 ns/op	         0.5412 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-8            	   18652	     66159 ns/op	         0.6616 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-8           	     718	   1641733 ns/op	         1.642 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-8       	174001639	         6.761 ns/op	         0.6761 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-8      	13432494	        88.84 ns/op	         0.8884 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-8     	 1611127	       742.8 ns/op	         0.7428 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-8    	  153541	      7473 ns/op	         0.7473 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-8   	   13296	     92093 ns/op	         0.9209 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-8  	     571	   2184055 ns/op	         2.184 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-8                 	10823258	       106.5 ns/op	        10.65 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-8                	 1707716	       695.4 ns/op	         6.954 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-8               	  148258	      7764 ns/op	         7.764 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-8              	   14757	     80434 ns/op	         8.043 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-8             	    1567	    706002 ns/op	         7.060 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-8            	     100	  10560363 ns/op	        10.56 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	101.789s
~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
