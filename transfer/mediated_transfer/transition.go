package mediated_transfer

import "github.com/SmartMeshFoundation/raiden-network/transfer"

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

/*

def update_route(next_state, route_state_change):

    if new_route.state != CHANNEL_STATE_OPENED:
        available_routes.pop(available_idx)

    elif new_route.state == CHANNEL_STATE_OPENED:
        if available_idx:
             overwrite it, balance might have changed
            available_routes[available_idx] = new_route

        else:
             TODO: re-add the new_route into the available_routes list if it can be used.
            ignored = any(
                route.node_address == new_route.node_address
                for route in next_state.routes.ignored_routes
            )

            canceled = any(
                route.node_address == new_route.node_address
                for route in next_state.routes.canceled_routes
            )

            if not canceled and not ignored:
                 new channel opened, add the route for use
                available_routes.append(new_route)

    next_state.routes.available_routes = available_routes

*/
