package streamserver

import "log"

type ConnLimiter struct {
	concurrentConn int
	bucket chan int
}

func newConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket: make(chan int, cc),
	}
}

func (cl *ConnLimiter) getConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate connection")
		return false
	}

	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) releaseConn() {
	c := <- cl.bucket
	log.Printf("New connection comming: %s", c)
}