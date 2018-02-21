package raiden_network

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/labstack/gommon/log"
)

/*
send ping to detect target is online or not
*/

type RoutesTask struct {
	NewTask           chan *RoutesToDetect
	TaskResult        chan *RoutesToDetect
	PingSender        network.PingSender
	NodesStatusGetter network.NodesStatusGetter
}

type RoutesToDetect struct {
	RoutesState     *transfer.RoutesState
	StateManager    *transfer.StateManager
	InitStateChange transfer.StateChange
}

func NewRoutesTask(sender network.PingSender, statusGetter network.NodesStatusGetter) *RoutesTask {
	return &RoutesTask{
		NewTask:           make(chan *RoutesToDetect, 10),
		TaskResult:        make(chan *RoutesToDetect, 10),
		PingSender:        sender,
		NodesStatusGetter: statusGetter,
	}
}

func (this *RoutesTask) Start() {
	go this.loop()
}
func (this *RoutesTask) Stop() {
	close(this.NewTask)
	//let task result open,otherwise may crash
}
func (this *RoutesTask) loop() {
	for {
		task, ok := <-this.NewTask
		if !ok {
			break //user stop
		}
		this.startTask(task)
	}
}

func (this *RoutesTask) startTask(task *RoutesToDetect) {
	var availables []*transfer.RouteState
	for i := 0; i < len(task.RoutesState.AvailableRoutes); i++ {
		status, lastVisitTime := this.NodesStatusGetter.GetNetworkStatusAndLastVisitTime(task.RoutesState.AvailableRoutes[i].HopNode)
		if status == network.NODE_NETWORK_REACHABLE && lastVisitTime.Add(5*time.Second).After(time.Now()) {
			continue //just detect seconds ago
		}
		err := this.PingSender.SendPing(task.RoutesState.AvailableRoutes[i].HopNode)
		if err != nil {
			log.Error(fmt.Sprintf("sendping to %s err:%s", task.RoutesState.AvailableRoutes[i].HopNode.String(), err))
		}
	}
	//wait ack 3 seconds, long or short?
	time.Sleep(3 * time.Second)
	for i := 0; i < len(task.RoutesState.AvailableRoutes); i++ {
		status, lastVisitTime := this.NodesStatusGetter.GetNetworkStatusAndLastVisitTime(task.RoutesState.AvailableRoutes[i].HopNode)
		if status == network.NODE_NETWORK_REACHABLE && lastVisitTime.Add(10*time.Second).After(time.Now()) {
			availables = append(availables, task.RoutesState.AvailableRoutes[i])
		} else {
			task.RoutesState.IgnoredRoutes = append(task.RoutesState.IgnoredRoutes, task.RoutesState.AvailableRoutes[i])
		}
	}
	task.RoutesState.AvailableRoutes = availables
	this.TaskResult <- task
}
