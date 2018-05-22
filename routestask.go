package smartraiden

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
)

/*
send ping to detect neighbours is online or not
*/

type routesTask struct {
	NewTask           chan *routesToDetect
	TaskResult        chan *routesToDetect
	PingSender        network.PingSender
	NodesStatusGetter network.NodesStatusGetter
	enabled           bool
}

type routesToDetect struct {
	RoutesState     *transfer.RoutesState
	StateManager    *transfer.StateManager
	InitStateChange transfer.StateChange
}

func newRoutesTask(sender network.PingSender, statusGetter network.NodesStatusGetter) *routesTask {
	return &routesTask{
		NewTask:           make(chan *routesToDetect, 10),
		TaskResult:        make(chan *routesToDetect, 10),
		PingSender:        sender,
		NodesStatusGetter: statusGetter,
		//enabled:           true, //每个节点都要进行探测,尤其是在 ice 模式下,很耗费资源,并且效果不佳,暂时禁用.
	}
}

func (rt *routesTask) start() {
	go rt.loop()
}
func (rt *routesTask) stop() {
	close(rt.NewTask)
	//let task result open,otherwise may crash
}
func (rt *routesTask) loop() {
	for {
		task, ok := <-rt.NewTask
		if !ok {
			break //user stop
		}
		rt.startTask(task)
	}
}

func (rt *routesTask) startTask(task *routesToDetect) {
	if rt.enabled {
		var availables []*transfer.RouteState
		var needWait = true
		pingcnt := 0
		const MaxPingOneTime = 10
		for i := 0; i < len(task.RoutesState.AvailableRoutes); i++ {
			status, lastAckTime := rt.NodesStatusGetter.GetNetworkStatusAndLastAckTime(task.RoutesState.AvailableRoutes[i].HopNode)
			if status == network.NodeNetworkReachable && lastAckTime.Add(time.Minute).After(time.Now()) {
				if i == 0 {
					needWait = false
				}
				continue //just detect seconds ago
			}
			err := rt.PingSender.SendPing(task.RoutesState.AvailableRoutes[i].HopNode)
			if err != nil {
				log.Error(fmt.Sprintf("sendping to %s err:%s", task.RoutesState.AvailableRoutes[i].HopNode.String(), err))
			}
			pingcnt++
			if pingcnt >= MaxPingOneTime {
				break
			}
		}
		//wait ack 5 seconds, long or short?
		needWait = false //for test only,shoulde be removed.
		if needWait {
			time.Sleep(5 * time.Second)
			for i := 0; i < len(task.RoutesState.AvailableRoutes); i++ {
				status, lastAckTime := rt.NodesStatusGetter.GetNetworkStatusAndLastAckTime(task.RoutesState.AvailableRoutes[i].HopNode)
				if status == network.NodeNetworkReachable && lastAckTime.Add(time.Minute).After(time.Now()) {
					availables = append(availables, task.RoutesState.AvailableRoutes[i])
				} else {
					task.RoutesState.IgnoredRoutes = append(task.RoutesState.IgnoredRoutes, task.RoutesState.AvailableRoutes[i])
				}
			}
		} else {
			availables = task.RoutesState.AvailableRoutes
		}

		task.RoutesState.AvailableRoutes = availables
		rt.TaskResult <- task
	} else {
		rt.TaskResult <- task
	}
}
