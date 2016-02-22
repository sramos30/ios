package simulator

import (
	"flag"
	"github.com/golang/glog"
	"github.com/heidi-ann/hydra/msgs"
	"testing"
	"time"
)

func checkRequest(t *testing.T, req msgs.ClientRequest, ios []*msgs.Io) {
	ios[0].IncomingRequests <- req

	for id := range ios {
		select {
		case reply := <-(*ios[id]).OutgoingRequests:
			if reply != req {
				t.Error(reply)
			}
		case <-time.After(time.Millisecond):
			t.Error("Participant not responding")
		}
	}
}

func TestSimulator(t *testing.T) {
	flag.Parse()
	defer glog.Flush()

	// create a system of 3 nodes
	ios := RunSimulator(3)

	// check that 3 nodes were created
	if len(ios) != 3 {
		t.Error("Correct number of nodes not created")
	}

	// check that master can replicate a request when no failures occur
	request1 := msgs.ClientRequest{
		ClientID:  2,
		RequestID: 0,
		Request:   "update A 3"}

	checkRequest(t, request1, ios)

	request2 := msgs.ClientRequest{
		ClientID:  2,
		RequestID: 1,
		Request:   "get A"}

	checkRequest(t, request2, ios)

	request3 := msgs.ClientRequest{
		ClientID:  4,
		RequestID: 0,
		Request:   "get C"}

	checkRequest(t, request3, ios)
}