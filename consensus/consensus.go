/*
Package consensus implements the Unanimous local replication algorithm.

This is INCOMPLETE as it currently:
	- assumes that all state is persistent
	- master does not recovery and assumes 3 is the last index allocated
	- master does all of its own coordination
	- master handles only 1 request at a time
	- log size is limited to 100 entries
*/

package consensus

import (
	"github.com/golang/glog"
	"github.com/heidi-ann/hydra/msgs"
)

// Config describes the static configuration of the consensus algorithm
type Config struct {
	ID int // id of node
	N  int // size of cluster (nodes numbered 0 to N-1)
}

// Init runs the consensus algorithm.
// It will not return until the application is terminated.
func Init(io *msgs.Io, config Config) {

	// setup
	glog.Infof("Starting node %d of %d", config.ID, config.N)
	state := State{
		View:        0,
		Log:         make([]msgs.Entry, 100), //TODO: Fix this
		CommitIndex: -1,
		MasterID:    0,
		LastIndex:   -1}

	// if master, start master goroutine
	if config.ID == 0 {
		glog.Info("Starting leader module")
		go RunMaster(0, 0, io, config)
	}

	// operator as normal node
	glog.Info("Starting participant module, ID ", config.ID)
	RunParticipant(state, io, config)

}
