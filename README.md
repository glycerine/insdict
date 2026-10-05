insdict: deterministic insert-ordered iteration hash table for Golang
=======

Dict provides a hash table for Go that iterates in insert-order for deterministic 
(reproducible) range over All(). Just like the Python 3.7+ dict dictionary.

We now use int64 indexes so that Dict size is not limited to 2^31.

In terms of performance, the common full table scans over All() are 10x faster than the built-in
go map. Our point operations are tied or slightly faster than the built-in map. Benchmarks follow.


# benchmarks of insdict (Dict) versus the built-in Go map (Map).

~~~
~~~

------------
Copyright(C) 2026 Jason E. Aten, Ph.D.

MIT License
