package simulator

import (
	"time"
)

const (
	maxUsageRecorded = 50
)

type UsageRecord struct {
	usingTooMany	bool
	using			map[string]time.Time
	usedByTooMany	bool
	usedBy			map[string]time.Time
}
type UsageTracker struct{ usage map[string]*UsageRecord }

func NewUsageTracker() *UsageTracker {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &UsageTracker{usage: make(map[string]*UsageRecord)}
}
func (tracker *UsageTracker) Get(node string) (data *UsageRecord, found bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, found = tracker.usage[node]
	return data, found
}
func (tracker *UsageTracker) RegisterUsage(nodeA string, nodeB string, timestamp time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if record, found := tracker.usage[nodeA]; found {
		if len(record.using) >= maxUsageRecorded {
			record.usingTooMany = true
		} else {
			record.using[nodeB] = timestamp
		}
	} else {
		record := UsageRecord{using: make(map[string]time.Time), usedBy: make(map[string]time.Time)}
		record.using[nodeB] = timestamp
		tracker.usage[nodeA] = &record
	}
	if record, found := tracker.usage[nodeB]; found {
		if len(record.usedBy) >= maxUsageRecorded {
			record.usedByTooMany = true
		} else {
			record.usedBy[nodeA] = timestamp
		}
	} else {
		record := UsageRecord{using: make(map[string]time.Time), usedBy: make(map[string]time.Time)}
		record.usedBy[nodeA] = timestamp
		tracker.usage[nodeB] = &record
	}
}
func (tracker *UsageTracker) Unregister(node string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if record, found := tracker.usage[node]; found {
		for using := range record.using {
			if record2, found := tracker.usage[using]; found {
				delete(record2.usedBy, node)
			}
		}
		for usedBy := range record.usedBy {
			if record2, found := tracker.usage[usedBy]; found {
				delete(record2.using, node)
			}
		}
		delete(tracker.usage, node)
	}
}
func filterOutOld(timestampMap map[string]time.Time, cutoff time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	toRemove := make([]string, 0)
	for key, timestamp := range timestampMap {
		if timestamp.Before(cutoff) {
			toRemove = append(toRemove, key)
		}
	}
	for _, key := range toRemove {
		delete(timestampMap, key)
	}
}
func (tracker *UsageTracker) CleanUp(cutoff time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	toDelete := make([]string, 0)
	for key, usageRecord := range tracker.usage {
		if !usageRecord.usingTooMany {
			filterOutOld(usageRecord.using, cutoff)
		}
		if !usageRecord.usedByTooMany {
			filterOutOld(usageRecord.usedBy, cutoff)
		}
		if !usageRecord.usingTooMany && !usageRecord.usedByTooMany && len(usageRecord.using) == 0 && len(usageRecord.usedBy) == 0 {
			toDelete = append(toDelete, key)
		}
	}
	for _, key := range toDelete {
		delete(tracker.usage, key)
	}
}
func RemoveNodeFromTracker(tracker *UsageTracker, node string, utilization map[string]time.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	keysToRemove := make([]string, 0)
	if mainRecord, found := tracker.Get(node); found {
		if mainRecord.usingTooMany {
			keysToRemove = getAllKeys(utilization)
		} else {
		usingloop:
			for usedNode := range mainRecord.using {
				if usedNodeRecord, found := tracker.Get(usedNode); found {
					if usedNodeRecord.usedByTooMany {
						keysToRemove = getAllKeys(utilization)
						break usingloop
					} else {
						for anotherNode := range usedNodeRecord.usedBy {
							keysToRemove = append(keysToRemove, anotherNode)
						}
					}
				}
			}
		}
	}
	tracker.Unregister(node)
	delete(utilization, node)
	for _, key := range keysToRemove {
		delete(utilization, key)
	}
}
func getAllKeys(m map[string]time.Time) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}
