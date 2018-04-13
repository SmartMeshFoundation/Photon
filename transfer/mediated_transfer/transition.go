package mediated_transfer

import "github.com/SmartMeshFoundation/SmartRaiden/transfer"

func UpdateRoute(Routes *transfer.RoutesState, stateChange *transfer.ActionRouteChangeStateChange) {
	newRoute := stateChange.Route
	idx := -1
	var oldRoute *transfer.RouteState
	availableRoutes := make([]*transfer.RouteState, len(Routes.AvailableRoutes))
	copy(availableRoutes, Routes.AvailableRoutes)
	for idx, oldRoute = range Routes.AvailableRoutes {
		if newRoute.HopNode == oldRoute.HopNode {
			break
		}
	}
	//this ActionRouteChangeStateChange has no relation with Routes.
	if len(Routes.AvailableRoutes) == 0 || idx == len(Routes.AvailableRoutes)-1 {
		return
	}
	//如果为空一定能找到?
	// TODO: what if the route that changed is the current route?
	if newRoute.State != transfer.CHANNEL_STATE_OPENED {
		//如果没有找到就替换最后一个?为什么呢
		availableRoutes = append(availableRoutes[0:idx], availableRoutes[idx+1:]...)
	} else {
		if idx >= 0 { //only when AvailableRoutes is empty
			//overwrite it, balance might have changed
			availableRoutes[idx] = newRoute
		} else {
			ignored := false
			canceled := false
			//  TODO: re-add the new_route into the available_routes list if it can be used.
			for _, r := range Routes.IgnoredRoutes {
				if r.HopNode == newRoute.HopNode {
					ignored = true
				}
			}
			for _, r := range Routes.CanceledRoutes {
				if r.HopNode == newRoute.HopNode {
					canceled = true
				}
			}
			if !canceled && !ignored {
				availableRoutes = append(availableRoutes, newRoute)
			}
		}
	}
	Routes.AvailableRoutes = availableRoutes
}
