package wrr

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

var _ base.PickerBuilder = &PickerBuilderStck{}
var _ balancer.Picker = &PickerStck{}

// 定义一个常量，表示负载均衡器的名称
const Name = "custom_weighted_round_robin"

// 创建一个新的负载均衡器构建器
func newBuilder() balancer.Builder {
	// 使用base包中的NewBalancerBuilder函数创建一个新的负载均衡器构建器
	return base.NewBalancerBuilder(Name, &PickerBuilderStck{}, base.Config{HealthCheck: true})
}

// 在程序初始化时注册负载均衡器构建器
func init() {
	// 使用balancer包中的Register函数注册负载均衡器构建器
	balancer.Register(newBuilder())
}

// PickerBuilderStck 传统版本的基于权重的负载均衡算法
type PickerBuilderStck struct {
}

func (p *PickerBuilderStck) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {

		//test 2024年12月11日12:02:45
		md := sci.Address.Metadata.(map[string]any) //获取地址的元数据
		weightVal, _ := md["weight"]
		weight := weightVal.(float64)

		fmt.Println(md)

		conns = append(conns, &weightConn{
			weight:        int(weight),
			currentWeight: int(weight),
			SubConn:       sc,
		})
	}

	return &PickerStck{
		conns: conns,
	}
}

type PickerStck struct {
	// 执行负载均衡的地方
	conns []*weightConn
	lock  sync.Mutex
}

// Pick 实现基于权重的负载均衡算法
func (p *PickerStck) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// 如果没有可用的连接，则返回错误
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// 总权重
	var total int         // 最大权重
	var maxCC *weightConn // 最大权重节点
	for _, c := range p.conns {
		total += c.weight                            // 累加权重
		c.currentWeight = c.currentWeight + c.weight // 更新当前权重

		//  找到最大权重节点 (或者)
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}

	maxCC.currentWeight = maxCC.currentWeight - total

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		Done: func(info balancer.DoneInfo) {
			// 要在这里进一步调整weight/currentWeight
			// failover 要在这里做文章
			// 根据调用结果的具体错误信息进行容错
			// 1. 如果要是触发了限流了，
			// 1.1 你可以考虑直接挪走这个节点，后面再挪回来
			// 1.2 你可以考虑直接将 weight/currentWeight 调整到极低
			// 2. 触发了熔断呢？
			// 3. 降级呢？
		},
	}, nil

	//return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}

//type conn struct {
//	weight        int              // 权重
//	currentWeight int              // 当前权重
//	cc            balancer.SubConn //grpc中代表一个节点的表达
//}

type weightConn struct {
	balancer.SubConn     // grpc中代表一个节点的表达
	weight           int // 权重
	currentWeight    int // 当前权重

	// 可以用来标记不可用
	available bool
}
