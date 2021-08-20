内容摘自

https://zhuanlan.zhihu.com/p/27328476

# Golang的runtime

![1610165780370](Golang的runtime.assets\1610165780370.png)

## runtime概述

​        runtime包含Go运行时的系统交互的操作，例如控制go runtine的功能。还有debug，pprof进行排查问题和运行时性能分析，tracer来抓取异常事件信息，如 go routine的创建，加锁解锁状态，系统调用进入推出和锁定还有GC相关的事件，堆栈大小的改变以及进程的退出和开始事件等等；race进行竞态关系检查以及CGO的实现。总的来说运行时是调度器和GC 



## scheduler

​         首先说到调度，我们学习操作系统时知道，对于CPU时间片的调度，是系统的资源分配策略，如任务A在执行完后，选择哪个任务来执行，使得某个因素（如进程总执行时间，或者磁盘寻道时间等）最小，达到最优的服务。这就是调度关注的问题。那么Go的运行时的scheduler是什么呢？我们为什么需要它，因为我们知道OS内核已经有一个线程（进程）scheduler了嘛？ 

![1610166096932](Golang的runtime.assets\1610166096932.png)

​        为什么Go还要自己搞一套？想想我们是不是经常说Go牛逼啊，语言级别实现了并发，我们为什么会这样说呢？原因就在于此，Go有自己的scheduler。 

​        说了这么多，到底为什么？我们知道线程有自己的信号掩码，上下文环境以及各种控制信息等，但这些很多特征对于Go程序本身来说并不关心， 而且context上下文切换的耗时费时费力费资源，更重要的是GC的原因，也是本文下部分说的，就是Go的垃圾回收需要stop the world，所有的goroutine停止，才能使得内存保持在一个一致的状态。垃圾回收的时间会根据内存情况变化是不确定的，如果我们没有自己的scheduler我们交给了OS自己的scheduler，我们就失去了控制，并且会有大量的线程需要停止工作。所以Go就需要自己单独的开发一个自己使用的调度器，能够自己管理goruntines，并且知道在什么时候内存状态是一致的，也就是说，对于OS而言运行时只需要为当时正在CPU核上运行的那个线程等待即可，而不是等待所有的线程。
每一个Go程序都附带一个runtime，runtime负责与底层操作系统交互，也都会有scheduler对goruntines进行调度。在scheduler中有三个非常重要的概念：P，M，G。 

 查看源码/src/runtime/proc.go我们可以看到注释： 

```go
// Goroutine scheduler
// The scheduler's job is to distribute ready-to-run goroutines over worker threads.
//
// The main concepts are:
// G - goroutine.
// M - worker thread, or machine.
// P - processor, a resource that is required to execute Go code.
//     M must have an associated P to execute Go code, however it can be
//     blocked or in a syscall w/o an associated P.
//
// Design doc at https://golang.org/s/go11sched.
```

 我们也看下Go程序的启动流程： 

```go
// The bootstrap sequence is:
//
//	call osinit
//	call schedinit
//	make & queue new G
//	call runtime·mstart
//
// The new G calls runtime·main.
```

![1610166417054](Golang的runtime.assets\1610166417054.png)

想要明白详细的流程可见：[golang internals - Genius0101 - 博客园](https://link.zhihu.com/?target=http%3A//www.cnblogs.com/genius0101/archive/2012/04/16/2447147.html)那么scheduler究竟解决了什么问题并如何管理goruntines呢？
*全部拷贝*
既然要调度那么肯定要有自己的调度策略了，go使用抢占式调度，goroutine的执行是可以被抢占的。如果一个goroutine一直占用CPU，长时间没有被调度过， 就会被runtime抢占掉，把CPU时间交给其他goroutine。详见： runtime在程序启动时，会自动创建一个系统线程，运行sysmon()函数， sysmon()函数在整个程序生命周期中一直执行，负责监视各个Goroutine的状态、判断是否要进行垃圾回收等，sysmon()会调用retake()函数，retake()函数会遍历所有的P，如果一个P处于执行状态， 且已经连续执行了较长时间，就会被抢占。
然后retake()调用preemptone()将P的stackguard0设为stackPreempt，这将导致该P中正在执行的G进行下一次函数调用时， 导致栈空间检查失败，进而触发morestack()，在goschedImpl()函数中，会通过调用dropg()将G与M解除绑定；再调用globrunqput()将G加入全局runnable队列中；最后调用schedule() 来用为当前P设置新的可执行的G。

![img](https://pic3.zhimg.com/80/v2-ebfdb863e8b40e3f95a5608ff9e446c6_720w.png)


如上图：go function 即可启动一个goroutine，所以每go出去一个语句被执行，runqueue队列就在其末尾加入一个goroutine，并在下一个调度点，就从runqueue中取出，一个goroutine执行。同时每个P可以转而投奔另一个OS线程，保证有足够的线程来运行所以的context P，也就是说goruntine可以在合适时机在多个OS线程间切换，也可以一直在一个线程，这由调度器决定。



GC1.7

![img](https://pic4.zhimg.com/80/v2-b144895d4724b1a80ca5caf9f8efc09b_720w.png)

可能从图上我们不好看出变化：在 Go 1.4 版本的时候它的 GC 在 300 毫秒的时候，但是在 1.5 版本 GC 已经优化得非常好了，压缩到了40 毫秒。从 1.6 版本的 15 到 20 毫秒升级到 1.63 版本的 5 毫秒。又从 1.6.3 升级到 1. 7 版本的 3 毫秒以内，同样在刚发布的1.8版本中，GC在低延迟方面的优化又给了我们大的惊喜，由于消除了GC的“stop-the-world stack re-scanning”，使得GC STW(stop-the-world)的时间通常低于100微秒，甚至经常低于10微秒，现在 GC 已经不是他们的问题了。GC 降下来了，CPU 使用率就上去了，1.7.3 和 1.8 版本中，CPU 会多利用一些，CPU 的使用率相对上升了一点，但是 GC 有很大的提升，当然这或多或少是以牺牲“吞吐”作为代价的，因此在Go 1.9中，GC的改进将持续进行，会在吞吐和低延迟上做一个很好的平衡。应该说，在 1.8 版本发布之后，1.9 版本现在引入了一个理念——goroutine 级别的GC，所以 1.9 版本可能还有更大的提升。
GC优化之路：

1. 1.3 以前，使用的是比较蠢的传统 Mark-Sweep 算法。
2. 1.3 版本进行了一下改进，把 Sweep 改为了并行操作。
3. 1.5 版本进行了较大改进，使用了改进三色标记算法，叫做“非分代的、非移动的、并发的、三色的标记清除垃圾收集器”，go 除了标准的三色收集以外，还有一个辅助回收功能，防止垃圾产生过快。分为两个主要阶段－markl阶段:GC对对象和不再使用的内存进行标记；sweep阶段，准备进行回收。这中间还分为两个子阶段，第一阶段，暂停应用，结束上一次sweep，接着进入并发mark阶段：找到正在使用的内存；第二阶段，mark结束阶段，这期间应用再一次暂停。最后，未使用的内存会被逐步回收，这个阶段是异步的，不会STW。
4. 1.6中，finalizer的扫描被移到了并发阶段中，对于大量连接的应用来说，GC的性能得到了显著提升。
5. 1.7号称史上改进最多的版本，在GC上的改进也很显著：并发的进行栈收缩，这样我们既实现了低延迟，又避免了对runtime进行调优，只要使用标准的runtime就可以。
6. 1.8 消除了GC的“stop-the-world stack re-scanning”

Go的GC目前来说已经做的非常好了，未来在1.9将更多在GC优化下对于吞吐和效率的平衡，我们一起期待！想要明白详细的流程可见：[golang internals - Genius0101 - 博客园](https://link.zhihu.com/?target=http%3A//www.cnblogs.com/genius0101/archive/2012/04/16/2447147.html)那么scheduler究竟解决了什么问题并如何管理goruntines呢？
*全部拷贝*
既然要调度那么肯定要有自己的调度策略了，go使用抢占式调度，goroutine的执行是可以被抢占的。如果一个goroutine一直占用CPU，长时间没有被调度过， 就会被runtime抢占掉，把CPU时间交给其他goroutine。详见： runtime在程序启动时，会自动创建一个系统线程，运行sysmon()函数， sysmon()函数在整个程序生命周期中一直执行，负责监视各个Goroutine的状态、判断是否要进行垃圾回收等，sysmon()会调用retake()函数，retake()函数会遍历所有的P，如果一个P处于执行状态， 且已经连续执行了较长时间，就会被抢占。
然后retake()调用preemptone()将P的stackguard0设为stackPreempt，这将导致该P中正在执行的G进行下一次函数调用时， 导致栈空间检查失败，进而触发morestack()，在goschedImpl()函数中，会通过调用dropg()将G与M解除绑定；再调用globrunqput()将G加入全局runnable队列中；最后调用schedule() 来用为当前P设置新的可执行的G。

![img](https://pic3.zhimg.com/80/v2-ebfdb863e8b40e3f95a5608ff9e446c6_720w.png)


如上图：go function 即可启动一个goroutine，所以每go出去一个语句被执行，runqueue队列就在其末尾加入一个goroutine，并在下一个调度点，就从runqueue中取出，一个goroutine执行。同时每个P可以转而投奔另一个OS线程，保证有足够的线程来运行所以的context P，也就是说goruntine可以在合适时机在多个OS线程间切换，也可以一直在一个线程，这由调度器决定。



GC1.7

![img](https://pic4.zhimg.com/80/v2-b144895d4724b1a80ca5caf9f8efc09b_720w.png)

可能从图上我们不好看出变化：在 Go 1.4 版本的时候它的 GC 在 300 毫秒的时候，但是在 1.5 版本 GC 已经优化得非常好了，压缩到了40 毫秒。从 1.6 版本的 15 到 20 毫秒升级到 1.63 版本的 5 毫秒。又从 1.6.3 升级到 1. 7 版本的 3 毫秒以内，同样在刚发布的1.8版本中，GC在低延迟方面的优化又给了我们大的惊喜，由于消除了GC的“stop-the-world stack re-scanning”，使得GC STW(stop-the-world)的时间通常低于100微秒，甚至经常低于10微秒，现在 GC 已经不是他们的问题了。GC 降下来了，CPU 使用率就上去了，1.7.3 和 1.8 版本中，CPU 会多利用一些，CPU 的使用率相对上升了一点，但是 GC 有很大的提升，当然这或多或少是以牺牲“吞吐”作为代价的，因此在Go 1.9中，GC的改进将持续进行，会在吞吐和低延迟上做一个很好的平衡。应该说，在 1.8 版本发布之后，1.9 版本现在引入了一个理念——goroutine 级别的GC，所以 1.9 版本可能还有更大的提升。
GC优化之路：

1. 1.3 以前，使用的是比较蠢的传统 Mark-Sweep 算法。
2. 1.3 版本进行了一下改进，把 Sweep 改为了并行操作。
3. 1.5 版本进行了较大改进，使用了改进三色标记算法，叫做“非分代的、非移动的、并发的、三色的标记清除垃圾收集器”，go 除了标准的三色收集以外，还有一个辅助回收功能，防止垃圾产生过快。分为两个主要阶段－markl阶段:GC对对象和不再使用的内存进行标记；sweep阶段，准备进行回收。这中间还分为两个子阶段，第一阶段，暂停应用，结束上一次sweep，接着进入并发mark阶段：找到正在使用的内存；第二阶段，mark结束阶段，这期间应用再一次暂停。最后，未使用的内存会被逐步回收，这个阶段是异步的，不会STW。
4. 1.6中，finalizer的扫描被移到了并发阶段中，对于大量连接的应用来说，GC的性能得到了显著提升。
5. 1.7号称史上改进最多的版本，在GC上的改进也很显著：并发的进行栈收缩，这样我们既实现了低延迟，又避免了对runtime进行调优，只要使用标准的runtime就可以。
6. 1.8 消除了GC的“stop-the-world stack re-scanning”

Go的GC目前来说已经做的非常好了，未来在1.9将更多在GC优化下对于吞吐和效率的平衡，我们一起期待！



## GO runtime包

    runtime
    尽管 Go 编译器产生的是本地可执行代码，这些代码仍旧运行在 Go 的 runtime（这部分的代码可以在 runtime 包中找到）当中。这个 runtime 类似 Java 和 .NET 语言所用到的虚拟机，它负责管理包括内存分配、垃圾回收、栈处理、goroutine、channel、切片（slice）、map 和反射（reflection）等等。
        
    runtime 调度器是个非常有用的东西，关于 runtime 包几个方法:
    
        Gosched：     让当前线程让出 cpu 以让其它线程运行,它不会挂起当前线程，因此当前线程未来会继续执行
        NumCPU：      返回当前系统的 CPU 核数量
        GOMAXPROCS：  设置最大的可同时使用的 CPU 核数
        Goexit：      退出当前 goroutine(但是defer语句会照常执行)
        NumGoroutine：返回正在执行和排队的任务总数
        GOOS：        目标操作系统
        Caller:   
```go
package main

import (
	"fmt"
	"runtime"
)

func f1() {
	pc, file, line, ok := runtime.Caller(1)
	// pc当前调用的函数名，runtime.FuncForPc(pc).name()
	// file 当n为0时为当前执行的文件全路径
	// line 执行该方法的第几行
	// 执行是否成功
	// n要得到执行调用的层级
	funcname := runtime.FuncForPC(pc).Name()
	fmt.Println(funcname, file, line, ok)
}
func f2() {
	f1()
}
func main() {
	f2() // 多调用一层 n为1
}s
```

