package worker

import "testing"

func Test_worker1(t *testing.T) {
	Initialize(0, 0, "127.0.0.1", []string{"127.0.0.1:2181"})
	StartWork()
}
