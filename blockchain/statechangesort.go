package blockchain

import (
	"sort"

	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
)

type contractStateChangeSlice []mediatedtransfer.ContractStateChange

func (c contractStateChangeSlice) Len() int {
	return len(c)
}
func (c contractStateChangeSlice) Less(i, j int) bool {
	return c[i].GetBlockNumber() < c[j].GetBlockNumber()
}
func (c contractStateChangeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

/*
must be stable
对于 ChannelOpenedAndDeposit 事件,会产生两个 stateChange,
严格要求有先后顺序
*/
/*
 *  This function must be stable.
 * 	As for ChannelOpenedAndDeposit event, two stateChange will be generated,
 *  and they must be in order.
 */
func sortContractStateChange(chs []mediatedtransfer.ContractStateChange) {
	sort.Stable(contractStateChangeSlice(chs))
}
