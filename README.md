insdict: deterministic insert-ordered iteration hash table for Golang
=======

Dict provides a hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.

We now use int64 indexes so that Dict size is not limited to 2^31.

In terms of performance, the common full table scans over All() are 10x faster than the built-in
go map. Our point operations are tied or slightly faster than the built-in map. Benchmarks follow.

# benchmarks of insdict (Dict) versus the built-in Go map (Map).

~~~
$ go test -v -run=xxx -bench=.
goos: linux
goarch: amd64
pkg: github.com/glycerine/insdict
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkDictPut
BenchmarkDictPut/n=10
BenchmarkDictPut/n=10-48            	 1840042	       649.3 ns/op	        64.93 put_ns/key	     696 B/op	       6 allocs/op
BenchmarkDictPut/n=100
BenchmarkDictPut/n=100-48           	  132122	      9668 ns/op	        96.68 put_ns/key	   16344 B/op	      18 allocs/op
BenchmarkDictPut/n=1000
BenchmarkDictPut/n=1000-48          	   14634	     80237 ns/op	        80.24 put_ns/key	  134616 B/op	      27 allocs/op
BenchmarkDictPut/n=10000
BenchmarkDictPut/n=10000-48         	    1394	    863299 ns/op	        86.33 put_ns/key	 1015256 B/op	      36 allocs/op
BenchmarkDictPut/n=100000
BenchmarkDictPut/n=100000-48        	     116	   9864317 ns/op	        98.64 put_ns/key	15941091 B/op	      48 allocs/op
BenchmarkDictPut/n=1000000
BenchmarkDictPut/n=1000000-48       	      10	 112526159 ns/op	       112.5 put_ns/key	127278851 B/op	      59 allocs/op
BenchmarkMapPut
BenchmarkMapPut/n=10
BenchmarkMapPut/n=10-48             	 1661924	       738.8 ns/op	        73.88 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPut/n=100
BenchmarkMapPut/n=100-48            	  131948	      8618 ns/op	        86.18 put_ns/key	    4456 B/op	       9 allocs/op
BenchmarkMapPut/n=1000
BenchmarkMapPut/n=1000-48           	    9469	    122765 ns/op	       122.8 put_ns/key	   74264 B/op	      20 allocs/op
BenchmarkMapPut/n=10000
BenchmarkMapPut/n=10000-48          	    1059	   1120873 ns/op	       112.1 put_ns/key	  591482 B/op	      79 allocs/op
BenchmarkMapPut/n=100000
BenchmarkMapPut/n=100000-48         	     100	  12183354 ns/op	       121.8 put_ns/key	 4729359 B/op	     530 allocs/op
BenchmarkMapPut/n=1000000
BenchmarkMapPut/n=1000000-48        	       6	 192212614 ns/op	       192.2 put_ns/key	75432496 B/op	    8190 allocs/op
BenchmarkMapPutPresized
BenchmarkMapPutPresized/n=10
BenchmarkMapPutPresized/n=10-48     	 2511164	       478.7 ns/op	        47.87 put_ns/key	     328 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=100
BenchmarkMapPutPresized/n=100-48    	  333108	      3317 ns/op	        33.17 put_ns/key	    2344 B/op	       3 allocs/op
BenchmarkMapPutPresized/n=1000
BenchmarkMapPutPresized/n=1000-48   	   34291	     35026 ns/op	        35.03 put_ns/key	   36944 B/op	       5 allocs/op
BenchmarkMapPutPresized/n=10000
BenchmarkMapPutPresized/n=10000-48  	    3007	    400571 ns/op	        40.06 put_ns/key	  295552 B/op	      33 allocs/op
BenchmarkMapPutPresized/n=100000
BenchmarkMapPutPresized/n=100000-48 	     199	   5901361 ns/op	        59.01 put_ns/key	 2364945 B/op	     257 allocs/op
BenchmarkMapPutPresized/n=1000000
BenchmarkMapPutPresized/n=1000000-48         	       9	 123909568 ns/op	       123.9 put_ns/key	37832990 B/op	    4099 allocs/op
BenchmarkDictPutPresized
BenchmarkDictPutPresized/n=10
BenchmarkDictPutPresized/n=10-48             	 2268607	       523.2 ns/op	        52.32 put_ns/key	     560 B/op	       4 allocs/op
BenchmarkDictPutPresized/n=100
BenchmarkDictPutPresized/n=100-48            	  239710	      5259 ns/op	        52.59 put_ns/key	    8544 B/op	       4 allocs/op
BenchmarkDictPutPresized/n=1000
BenchmarkDictPutPresized/n=1000-48           	   27454	     43714 ns/op	        43.71 put_ns/key	   67680 B/op	       4 allocs/op
BenchmarkDictPutPresized/n=10000
BenchmarkDictPutPresized/n=10000-48          	    2677	    455203 ns/op	        45.52 put_ns/key	  499808 B/op	       4 allocs/op
BenchmarkDictPutPresized/n=100000
BenchmarkDictPutPresized/n=100000-48         	     187	   6277602 ns/op	        62.78 put_ns/key	 7954535 B/op	       4 allocs/op
BenchmarkDictPutPresized/n=1000000
BenchmarkDictPutPresized/n=1000000-48        	      13	  83314926 ns/op	        83.31 put_ns/key	63619478 B/op	       6 allocs/op
BenchmarkDictPutOverwrite
BenchmarkDictPutOverwrite/n=10
BenchmarkDictPutOverwrite/n=10-48            	128248502	         9.069 ns/op	         9.069 put_ns/key
BenchmarkDictPutOverwrite/n=100
BenchmarkDictPutOverwrite/n=100-48           	118274371	        10.06 ns/op	        10.06 put_ns/key
BenchmarkDictPutOverwrite/n=1000
BenchmarkDictPutOverwrite/n=1000-48          	112684245	        10.31 ns/op	        10.31 put_ns/key
BenchmarkDictPutOverwrite/n=10000
BenchmarkDictPutOverwrite/n=10000-48         	95371587	        12.69 ns/op	        12.69 put_ns/key
BenchmarkDictPutOverwrite/n=100000
BenchmarkDictPutOverwrite/n=100000-48        	77702868	        15.12 ns/op	        15.12 put_ns/key
BenchmarkDictPutOverwrite/n=1000000
BenchmarkDictPutOverwrite/n=1000000-48       	18692292	        62.53 ns/op	        62.53 put_ns/key
BenchmarkMapPutOverwrite
BenchmarkMapPutOverwrite/n=10
BenchmarkMapPutOverwrite/n=10-48             	100000000	        12.03 ns/op	        12.03 put_ns/key
BenchmarkMapPutOverwrite/n=100
BenchmarkMapPutOverwrite/n=100-48            	92382820	        13.17 ns/op	        13.17 put_ns/key
BenchmarkMapPutOverwrite/n=1000
BenchmarkMapPutOverwrite/n=1000-48           	82489485	        14.38 ns/op	        14.38 put_ns/key
BenchmarkMapPutOverwrite/n=10000
BenchmarkMapPutOverwrite/n=10000-48          	72546896	        16.05 ns/op	        16.05 put_ns/key
BenchmarkMapPutOverwrite/n=100000
BenchmarkMapPutOverwrite/n=100000-48         	51026212	        22.90 ns/op	        22.90 put_ns/key
BenchmarkMapPutOverwrite/n=1000000
BenchmarkMapPutOverwrite/n=1000000-48        	11914714	       103.0 ns/op	       103.0 put_ns/key
BenchmarkDictGet
BenchmarkDictGet/n=10
BenchmarkDictGet/n=10-48                     	188871216	         6.327 ns/op	         6.327 get_ns/key
BenchmarkDictGet/n=100
BenchmarkDictGet/n=100-48                    	186473864	         6.386 ns/op	         6.386 get_ns/key
BenchmarkDictGet/n=1000
BenchmarkDictGet/n=1000-48                   	174650648	         6.845 ns/op	         6.845 get_ns/key
BenchmarkDictGet/n=10000
BenchmarkDictGet/n=10000-48                  	118308538	        10.16 ns/op	        10.16 get_ns/key
BenchmarkDictGet/n=100000
BenchmarkDictGet/n=100000-48                 	102340584	        11.60 ns/op	        11.60 get_ns/key
BenchmarkDictGet/n=1000000
BenchmarkDictGet/n=1000000-48                	23547938	        47.70 ns/op	        47.70 get_ns/key
BenchmarkMapGet
BenchmarkMapGet/n=10
BenchmarkMapGet/n=10-48                      	146293527	         9.809 ns/op	         9.809 get_ns/key
BenchmarkMapGet/n=100
BenchmarkMapGet/n=100-48                     	137665294	         9.375 ns/op	         9.375 get_ns/key
BenchmarkMapGet/n=1000
BenchmarkMapGet/n=1000-48                    	127637856	         9.185 ns/op	         9.185 get_ns/key
BenchmarkMapGet/n=10000
BenchmarkMapGet/n=10000-48                   	117256420	        10.10 ns/op	        10.10 get_ns/key
BenchmarkMapGet/n=100000
BenchmarkMapGet/n=100000-48                  	84523376	        14.05 ns/op	        14.05 get_ns/key
BenchmarkMapGet/n=1000000
BenchmarkMapGet/n=1000000-48                 	17750324	        67.21 ns/op	        67.21 get_ns/key
BenchmarkDictIterate
BenchmarkDictIterate/n=10
BenchmarkDictIterate/n=10-48                 	195536712	         5.981 ns/op	         0.5981 iter_ns/key
BenchmarkDictIterate/n=100
BenchmarkDictIterate/n=100-48                	22705027	        55.01 ns/op	         0.5501 iter_ns/key
BenchmarkDictIterate/n=1000
BenchmarkDictIterate/n=1000-48               	 2368802	       480.7 ns/op	         0.4807 iter_ns/key
BenchmarkDictIterate/n=10000
BenchmarkDictIterate/n=10000-48              	  201808	      5875 ns/op	         0.5875 iter_ns/key
BenchmarkDictIterate/n=100000
BenchmarkDictIterate/n=100000-48             	   19005	     60771 ns/op	         0.6077 iter_ns/key
BenchmarkDictIterate/n=1000000
BenchmarkDictIterate/n=1000000-48            	     637	   1821575 ns/op	         1.822 iter_ns/key
BenchmarkDictIterateWithHoles
BenchmarkDictIterateWithHoles/n=10
BenchmarkDictIterateWithHoles/n=10-48        	130624459	         8.994 ns/op	         0.8994 iter_ns/key
BenchmarkDictIterateWithHoles/n=100
BenchmarkDictIterateWithHoles/n=100-48       	12735891	        94.88 ns/op	         0.9488 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000
BenchmarkDictIterateWithHoles/n=1000-48      	 1335889	       890.9 ns/op	         0.8909 iter_ns/key
BenchmarkDictIterateWithHoles/n=10000
BenchmarkDictIterateWithHoles/n=10000-48     	  136861	      8804 ns/op	         0.8804 iter_ns/key
BenchmarkDictIterateWithHoles/n=100000
BenchmarkDictIterateWithHoles/n=100000-48    	   13524	     87543 ns/op	         0.8754 iter_ns/key
BenchmarkDictIterateWithHoles/n=1000000
BenchmarkDictIterateWithHoles/n=1000000-48   	     476	   2477983 ns/op	         2.478 iter_ns/key
BenchmarkMapIterate
BenchmarkMapIterate/n=10
BenchmarkMapIterate/n=10-48                  	 8641455	       116.7 ns/op	        11.67 iter_ns/key
BenchmarkMapIterate/n=100
BenchmarkMapIterate/n=100-48                 	 1634574	       732.2 ns/op	         7.322 iter_ns/key
BenchmarkMapIterate/n=1000
BenchmarkMapIterate/n=1000-48                	  127077	      9033 ns/op	         9.033 iter_ns/key
BenchmarkMapIterate/n=10000
BenchmarkMapIterate/n=10000-48               	   14263	     83944 ns/op	         8.394 iter_ns/key
BenchmarkMapIterate/n=100000
BenchmarkMapIterate/n=100000-48              	    1672	    701750 ns/op	         7.017 iter_ns/key
BenchmarkMapIterate/n=1000000
BenchmarkMapIterate/n=1000000-48             	     100	  10857856 ns/op	        10.86 iter_ns/key
PASS
ok  	github.com/glycerine/insdict	119.544s

~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
