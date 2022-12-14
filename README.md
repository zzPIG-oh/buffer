# buffer
介于redis、服务之间的内存缓存系统
***
说明：传统的项目加速都会使用redis、memcached做缓存，牺牲网络IO减少磁盘IO的时间。buffer是想利用本地缓存的数据减少掉部分网络IO的时间。引入本地缓存会面临很多问题，内存数据的增大导致系统运行缓慢、多副本间本地缓存不平衡等等。buffer还支持制定数据源，用于服务启动时就加载进内存。  
buffer的定位是服务、redis之间，也就是操作redis也交给buffer来管理，同时buffer会缓存数据到本地。如何解决上述两个问题？  
* 多副本间的数据有redis来做数据的同步，一开始是想利用redis的pubsub来做多副本间的数据同步，但按照正常的系统设计，缓存不命中的时候会反查db，所以buffer使用了redis来做多副本间数据不均衡的问题。同时要注意，多副本使用的是同一个redis实例。
* 问题2解决方式buffer引入了两个标识来做内存的管控，可以设置。buffer会进行内存数据巡检，根据有无过期时间进行两种数据回收，有过期时间的key，buffer会判断是否过期，无过期时间的，buffer会根据热度来回收，当前实现还较为简单。
***
由于buffer是来自于日常开发中，当前buffer仅满足了hashmap的数据存取。
***
如何使用?  

```
func TestBufferSet(t *testing.T) {
	defaultFile = "../../../buffer.json"
	Init()
	bc := NewBufferClient()

	bl1 := bc.Hget("keyword", "zhangzhen").IsEmpty()
	bc.Hset("keyword", "zhangzhen", Empty)
	bl2 := bc.Hget("keyword", "zhangzhen").IsEmpty()
	fmt.Println(bl1, bl2)
}
```


