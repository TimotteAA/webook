package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

// 基于超时次数的轮询：如果单个服务超时次数超过n次，则切换服务
type failoverServiceV2 struct {
	svcs []sms.Service
	// 当前正在使用服务的超时次数
	count int64
	// 当前正在使用的服务
	idx int64
	// 超时上限
	threshold int64
}

// 判断理由：如果一个第三方服务不可用，最直接会出现的情况是连续n次超时
func NewFailoverServiceV2(svcs []sms.Service, threshold int64) sms.Service {
	return &failoverServiceV2{svcs: svcs, count: 0, idx: 0, threshold: threshold}
}

func (f *failoverServiceV2) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 由于并发操作，需要用到原子类
	// 当前正在使用的服务
	idx := atomic.LoadInt64(&f.idx)
	cnt := atomic.LoadInt64(&f.count)
	len := len(f.svcs)

	// 超时，切换服务
	if cnt > f.threshold {
		// 重新计算idx
		newIdx := (idx + int64(1)) % int64(len)
		// CAS操作失败，这说明别的地方已经切换了，这里就不用切换了，避免重复切换
		if atomic.CompareAndSwapInt64(&f.idx, idx, newIdx) {
			atomic.StoreInt64(&f.count, 0)
		}
	}
	// 未超时
	err := f.svcs[idx].Send(ctx, tplId, args, numbers...)
	switch err {
	// 未超时
	case nil:
		// 重置超时时间
		atomic.StoreInt64(&f.count, 0)
	// 超时了
	case context.DeadlineExceeded:
		atomic.AddInt64(&f.count, 1)
	}
	return err
}
