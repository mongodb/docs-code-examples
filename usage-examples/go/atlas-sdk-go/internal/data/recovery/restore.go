package recovery

import (
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// AddElectableNodesToRegion increases electable node count in the specified target region.
func AddElectableNodesToRegion(repl []admin.ReplicationSpec20240805, targetRegion string, addNodes int) (int, bool) {
	added := 0
	found := false
	for i := range repl {
		rcs := repl[i].GetRegionConfigs()
		for j := range rcs {
			regionName := ""
			if rcs[j].HasRegionName() {
				regionName = rcs[j].GetRegionName()
			}
			if regionName == targetRegion && rcs[j].HasElectableSpecs() {
				es := rcs[j].GetElectableSpecs()
				before := 0
				if es.HasNodeCount() {
					before = es.GetNodeCount()
				}
				es.SetNodeCount(before + addNodes)
				rcs[j].SetElectableSpecs(es)
				added += addNodes
				found = true
			}
		}
		repl[i].SetRegionConfigs(rcs)
	}
	return added, found
}

// ZeroElectableNodesInRegion sets electable node count to zero in the outage region, returning count of regions modified.
func ZeroElectableNodesInRegion(repl []admin.ReplicationSpec20240805, outageRegion string) int {
	zeroed := 0
	for i := range repl {
		rcs := repl[i].GetRegionConfigs()
		for j := range rcs {
			regionName := ""
			if rcs[j].HasRegionName() {
				regionName = rcs[j].GetRegionName()
			}
			if regionName == outageRegion && rcs[j].HasElectableSpecs() {
				es := rcs[j].GetElectableSpecs()
				if es.HasNodeCount() && es.GetNodeCount() > 0 {
					es.SetNodeCount(0)
					rcs[j].SetElectableSpecs(es)
					zeroed++
				}
			}
		}
		repl[i].SetRegionConfigs(rcs)
	}
	return zeroed
}
