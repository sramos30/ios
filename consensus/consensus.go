/*
Package consensus implements the Unanimous local replication algorithm.

This is INCOMPLETE as it currently:
	- assumes that all state is persistent
	- master does not recovery and assumes 3 is the last index allocated
	- master does all of its own coordination
	- master handles only 1 request at a time
	- log size is limited to 1000 entries
*/

package consensus

import (
	"github.com/golang/glog"
	"github.com/heidi-ann/hydra/msgs"
)

// Config describes the static configuration of the consensus algorithm
type Config struct {
	ID        int // id of node
	N         int // size of cluster (nodes numbered 0 to N-1)
	LogLength int // max log size
	BatchInterval  int // how often to batch process request in ms, 0 means no batching
	MaxBatch int // maximum requests in a batch, unused if BatchInterval=0
}

// Init runs the consensus algorithm.
// It will not return until the application is terminated.
func Init(io *msgs.Io, config Config) {

	// setup
	glog.Infof("Starting node %d of %d", config.ID, config.N)
	state := State{
		View:        0,
		Log:         make([]msgs.Entry, config.LogLength),
		CommitIndex: -1,
		MasterID:    0,
		LastIndex:   -1}

	// write intial term to persistent storage
	(*io).ViewPersist <- 0

	// if master, start master goroutine
	if config.ID == 0 {
		glog.Info("Starting leader module")
		go RunMaster(0, -1, true, io, config)
	}

	// operator as normal node
	glog.Info("Starting participant module, ID ", config.ID)
	RunParticipant(state, io, config)

}

func Recover(io *msgs.Io, config Config, view int, log []msgs.Entry) {
	// setup
	glog.Infof("Restarting node %d of %d", config.ID, config.N)

	new_log := make([]msgs.Entry, config.LogLength)
	copy(new_log, log)
	state := State{
		View:        view,
		Log:         new_log,
		CommitIndex: -1,
		MasterID:    mod(view, config.N),
		LastIndex:   len(log) - 1}

	// if master, start master goroutine
	if config.ID == state.MasterID {
		glog.Info("Starting leader module")
		go RunMaster(view, -1, true, io, config)
	}

	// apply recovered requests to state machine
	for i := 0; i <= state.LastIndex; i++ {
		if !state.Log[i].Committed {
			break
		}
		state.CommitIndex = i

		for _, request := range state.Log[i].Requests {
			(*io).OutgoingRequests <- request
		}
	}

	// operator as normal node
	glog.Info("Starting participant module, ID ", config.ID)
	RunParticipant(state, io, config)

}
