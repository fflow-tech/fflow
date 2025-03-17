# 本地缓存功能开发设计文档（Feature Dev Spec）

[TOC]

## 产品需求文档（PRD）

Workflow 需要一个非常快速的本地缓存服务，并满足以下几点需求：

- 作为进程全局的缓存服务，支持并发访问；
- 能够存储结构化的数据，例如 struct；
- 一定时间后支持剔除；
- `SET`添加缓存后，`GET`能立即获取到最新结果；
- 内存限制（限制最大的可使用空间）；

## 技术选型

目前 Go 语言本地缓存的实现方案有：BigCacheSet、FreeCache、原生 map + 读写锁。在并发访问的条件下，原生 map + 读写锁的实现方式往往会出现写请求阻塞读请求的情况，造成大量协程阻塞。而 BigCacheSet 根据 key 的哈希将数据分为 shards。每个分片都包含一个映射和一个 ring buffer。每当设置新元素时，它都会将该元素追加到相应分片的 ring buffer 中，并且缓冲区中的偏移量将存储在 map 中。如果缓冲区太小，则将其扩展直到达到最大容量。。根据 [bigcache-bench](https://github.com/allegro/bigcache-bench) 的测试数据，在并行访问的条件下，BigCacheSet 的 Set/Get 所需的时间约为另外两种实现的一半。

```shell
BenchmarkBigCacheSetParallel-8        	34233472	       148 ns/op	     317 B/op	       3 allocs/op
BenchmarkFreeCacheSetParallel-8       	34222654	       268 ns/op	     350 B/op	       3 allocs/op
BenchmarkConcurrentMapSetParallel-8   	19635688	       240 ns/op	     200 B/op	       6 allocs/op
BenchmarkBigCacheGetParallel-8        	60547064	        86.1 ns/op	     152 B/op	       4 allocs/op
BenchmarkFreeCacheGetParallel-8       	50701280	       147 ns/op	     136 B/op	       3 allocs/op
BenchmarkConcurrentMapGetParallel-8   	27353288	       175 ns/op	      24 B/op	       2 allocs/op
PASS
ok  	github.com/allegro/bigcache/v3/caches_bench	256.257s

go run caches_gc_overhead_comparison.go
Number of entries:  20000000
GC pause for bigcache:  22.382827ms
GC pause for freecache:  41.264651ms
GC pause for map:  72.236853ms
```

除了读写性能上的优势，bigcache 的作者还对 GC 扫描做了一定的优化，这和 Go1.5 中一个修复有关([#9477](https://github.com/golang/go/issues/9477))，Go的开发者优化了垃圾回收时对于 map 的处理，如果 map 对象中的 key 和 value 不包含指针，虽然该 map 也是分配在堆上，但是 GC 可以无视它们。所以 bigcache 使用哈希值作为`map[int]int`的key。 把缓存对象序列化后放到一个预先分配的大的字节数组中，然后将它在数组中的 offset 作为`map[int]int`的 value。

## 干系人（Stakeholders）

beihaizheng; haidenzhang

## 类与接口设计（Classes & APIs）

- Set(key string, value []byte)：向 LocalCache 插入新数据；
- Get(key string)：查询 key 的值，如果 Key 不存在会返回 ErrEntryNotFound；
- Append(key string, value []byte)：在原有 value 后面追加新数据，如果 key 不存在则新建 key。

## 相关链接（Links）

- [Writing a very fast cache service with millions of entries in Go](https://blog.allegro.tech/2016/03/writing-fast-cache-service-in-go.html)
- [github.com/allegro/bigcache](https://github.com/allegro/bigcache)
- [bigcache优化技巧](https://colobu.com/2019/11/18/how-is-the-bigcache-is-fast/)
