https://www.liwenzhou.com/posts/Go/go_gopsutil/

## 包安装

```bash
go get -u github.com/shirou/gopsutil
```



## cpu

```shell
"github.com/shirou/gopsutil/cpu"
```



### cpu.Info()

```shell
func Info() ([]InfoStat, error) {}
```



返回：

```shell
[]InfoStat, error

type InfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}
```



### cpu.Percent()

cpu使用百分比

```shell
func Percent(interval time.Duration, percpu bool) ([]float64, error) {}
```



### cpu.Counts()

```shell
func Counts(logical bool) (int, error) {}
```

true默认为逻辑核,反之为物理核数



### cpu.ProcInfo()

```go
func ProcInfo() ([]Win32_PerfFormattedData_PerfOS_System, error) {}


type Win32_PerfFormattedData_PerfOS_System struct {
	Processes            uint32   // 进程数
	ProcessorQueueLength uint32   // 进程队列长度
}
```



### load.avg()

```go
import "github.com/shirou/gopsutil/load"
```

```
func Avg() (*AvgStat, error) {}

type AvgStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
```



### cpu相关示例

```golang
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"time"
	//"github.com/shirou/gopsutil/mem"
)

// getCpuInfo cpuinfo,返回[]InfoStat类型的cpu数量的切片和err
func getCpuInfo()([]cpu.InfoStat,error){
	res,err := cpu.Info()
	if err != nil {
		return nil,err
	}
	return res,nil
}

// getCpuPercent // cpup平均使用率，返回float类型的切片,percpu为true时显示每个cpu,从执行开始一直等待到时间完毕，时时计算的
                 // percpu为false 显示第0个cpu
func getCpuPercent(interval time.Duration,percpu bool)([]float64, error) {
	percent,err := cpu.Percent(interval,percpu)
	if err != nil {
		return nil,err
	}
	return percent,nil
}

// Getcpucount true默认为逻辑核,反之为物理核数
func Getcpucount(){
	c,_ := cpu.Counts(false)
	fmt.Println(c)
}

// getCpuLoad cpu平均负载
func getCpuLoad(){
	cpuload,_ := load.Avg()
	fmt.Println(cpuload.Load1,cpuload.Load5,cpuload.Load15)
}

func main() {
	// 获取cpu信息
	cpuinfos,_ := getCpuInfo()
	for k,v := range cpuinfos{
		fmt.Println("=======>",k)
		fmt.Println(v)
	}

	// 获取cpu利用率 平均5秒
	t := time.Second*5
	cpupercent,_ := getCpuPercent(t,false)
	for k,v := range cpupercent{
		fmt.Println(k)
		fmt.Println(v)
	}

	// cpu load
	getCpuLoad()
}
```





## host

```go
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
)

/*
type InfoStat struct {
	Hostname             string `json:"hostname"`
	Uptime               uint64 `json:"uptime"`
	BootTime             uint64 `json:"bootTime"`
	Procs                uint64 `json:"procs"`           // number of processes
	OS                   string `json:"os"`              // ex: freebsd, linux
	Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"` // guest or host
	HostID               string `json:"hostid"`             // ex: uuid
}
*/

// getHostinfo 获取host信息
func getHostinfo(){
	// host.info 返回 *host.*InfoStat
	hostinfo,_ := host.Info()
	fmt.Println(hostinfo.Hostname,hostinfo.Uptime,hostinfo.OS,hostinfo.Procs,hostinfo.Platform,hostinfo.KernelArch)
}

func main() {
	getHostinfo()
}
```



## mem

```go
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
)

/*
type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`

	// OS X / BSD specific numbers:
	// http://www.macyourself.com/2010/02/17/what-is-free-wired-active-and-inactive-system-memory-ram/
	Active   uint64 `json:"active"`
	Inactive uint64 `json:"inactive"`
	Wired    uint64 `json:"wired"`

	// FreeBSD specific numbers:
	// https://reviews.freebsd.org/D8467
	Laundry uint64 `json:"laundry"`

	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	Buffers        uint64 `json:"buffers"`
	Cached         uint64 `json:"cached"`
	Writeback      uint64 `json:"writeback"`
	Dirty          uint64 `json:"dirty"`
	WritebackTmp   uint64 `json:"writebacktmp"`
	Shared         uint64 `json:"shared"`
	Slab           uint64 `json:"slab"`
	SReclaimable   uint64 `json:"sreclaimable"`
	SUnreclaim     uint64 `json:"sunreclaim"`
	PageTables     uint64 `json:"pagetables"`
	SwapCached     uint64 `json:"swapcached"`
	CommitLimit    uint64 `json:"commitlimit"`
	CommittedAS    uint64 `json:"committedas"`
	HighTotal      uint64 `json:"hightotal"`
	HighFree       uint64 `json:"highfree"`
	LowTotal       uint64 `json:"lowtotal"`
	LowFree        uint64 `json:"lowfree"`
	SwapTotal      uint64 `json:"swaptotal"`
	SwapFree       uint64 `json:"swapfree"`
	Mapped         uint64 `json:"mapped"`
	VMallocTotal   uint64 `json:"vmalloctotal"`
	VMallocUsed    uint64 `json:"vmallocused"`
	VMallocChunk   uint64 `json:"vmallocchunk"`
	HugePagesTotal uint64 `json:"hugepagestotal"`
	HugePagesFree  uint64 `json:"hugepagesfree"`
	HugePageSize   uint64 `json:"hugepagesize"`
}
*/

// getMemoryInfo获取内存信息，返回指针类型结构体，total，available，used，free，usedPercent等
func getMemoryInfo()  {
	res,_ := mem.VirtualMemory()
	fmt.Println(res)
}

func main() {
	getMemoryInfo()
}
```



## Disk

### 源码

```go
type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type PartitionStat struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	Fstype     string `json:"fstype"`
	Opts       string `json:"opts"`
}

type IOCountersStat struct {
	ReadCount        uint64 `json:"readCount"`
	MergedReadCount  uint64 `json:"mergedReadCount"`
	WriteCount       uint64 `json:"writeCount"`
	MergedWriteCount uint64 `json:"mergedWriteCount"`
	ReadBytes        uint64 `json:"readBytes"`
	WriteBytes       uint64 `json:"writeBytes"`
	ReadTime         uint64 `json:"readTime"`
	WriteTime        uint64 `json:"writeTime"`
	IopsInProgress   uint64 `json:"iopsInProgress"`
	IoTime           uint64 `json:"ioTime"`
	WeightedIO       uint64 `json:"weightedIO"`
	Name             string `json:"name"`
	SerialNumber     string `json:"serialNumber"`
	Label            string `json:"label"`
}
........
// Usage returns a file system usage. path is a filesystem path such
// as "/", not device file path like "/dev/vda1".  If you want to use
// a return value of disk.Partitions, use "Mountpoint" not "Device".
func Usage(path string) (*UsageStat, error) {
	return UsageWithContext(context.Background(), path)
}

// Partitions returns disk partitions. If all is false, returns
// physical devices only (e.g. hard disks, cd-rom drives, USB keys)
// and ignore all others (e.g. memory partitions such as /dev/shm)
//
// 'all' argument is ignored for BSD, see: https://github.com/giampaolo/psutil/issues/906
func Partitions(all bool) ([]PartitionStat, error) {
	return PartitionsWithContext(context.Background(), all)
}

func IOCounters(names ...string) (map[string]IOCountersStat, error) {
	return IOCountersWithContext(context.Background(), names...)
}
```



### 示例

```go
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
)

// getDiskInfo获取逻辑盘信息，返回逻辑盘分区数长度切片
func getDiskInfo()  {
	res,_ := disk.Partitions(true)
	fmt.Println(res)

	usage,_ := disk.Usage("C:")
	fmt.Println(usage)

	ioStat, _ := disk.IOCounters()
	for k, v := range ioStat {
		fmt.Printf("%v:%v\n", k, v)
	}

}

func main() {
	getDiskInfo()
}
```



## net IO

```go
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
)

func getNetIo()  {
	info, _ := net.IOCounters(true)
	for index, v := range info {
		fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	}
}

func main() {
	getNetIo()
}
```




