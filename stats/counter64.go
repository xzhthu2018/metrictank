package stats

import "sync/atomic"

type Counter64 struct {
	val uint64
}

func NewCounter64(name string) *Counter64 {
	return registry.add(name, func() GraphiteMetric {
		return &Counter64{}
	}).(*Counter64)
}

func (c *Counter64) SetUint64(val uint64) {
	atomic.StoreUint64(&c.val, val)
}

func (c *Counter64) Inc() {
	atomic.AddUint64(&c.val, 1)
}

func (c *Counter64) AddUint64(val uint64) {
	atomic.AddUint64(&c.val, val)
}

func (c *Counter64) ReportGraphite(prefix, buf []byte, now int64) []byte {
	val := atomic.LoadUint64(&c.val)
	buf = WriteUint64(buf, prefix, nil, val, now)
	return buf
}
