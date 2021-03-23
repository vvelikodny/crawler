package crawler

import "sync/atomic"

type Stat interface {
	AddPage()
	PagesCount() int32
	AddTotalDiscovered()
	TotalDiscovered() int32
	AddUniqDiscovered()
	UniqDiscovered() int32
	AddTotalFetched()
	TotalFetched() int32
}

type PublicStat interface {
	PagesCount() int32
	TotalDiscovered() int32
	UniqDiscovered() int32
	TotalFetched() int32
}

type inMemStat struct {
	pagesCount      int32
	totalDiscovered int32
	uniqDiscovered  int32
	totalFetched    int32
}

func NewStat() Stat {
	return &inMemStat{}
}

func (s *inMemStat) AddPage() {
	atomic.AddInt32(&s.pagesCount, 1)
}

func (s *inMemStat) PagesCount() int32 {
	return atomic.LoadInt32(&s.pagesCount)
}

func (s *inMemStat) AddTotalDiscovered() {
	atomic.AddInt32(&s.totalDiscovered, 1)
}

func (s *inMemStat) TotalDiscovered() int32 {
	return atomic.LoadInt32(&s.totalDiscovered)
}

func (s *inMemStat) AddUniqDiscovered() {
	atomic.AddInt32(&s.uniqDiscovered, 1)
}

func (s *inMemStat) UniqDiscovered() int32 {
	return atomic.LoadInt32(&s.uniqDiscovered)
}

func (s *inMemStat) AddTotalFetched() {
	atomic.AddInt32(&s.totalFetched, 1)
}

func (s *inMemStat) TotalFetched() int32 {
	return atomic.LoadInt32(&s.totalFetched)
}
