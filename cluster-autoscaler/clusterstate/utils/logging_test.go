package utils

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestLogCollectorNoCompaction(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	logCollector := NewLogCollector()
	logCollector.Log("Event1", Debug)
	logCollector.Log("Event2", Info)
	logCollector.Log("Event3", Warning)
	log := logCollector.GetLogs()
	logCollector.Log("Event4", Error)
	assert.Equal(t, len(log), 3)
	assert.Equal(t, "Event1", log[0].Log)
	assert.EqualValues(t, Debug, log[0].Level)
	assert.Equal(t, "Event3", log[2].Log)
	assert.EqualValues(t, Warning, log[2].Level)
}
func TestLogCollectorSizeCompaction(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	logCollector := NewLogCollector()
	logCollector.maxItems = 2
	logCollector.Log("Event1", Info)
	logCollector.Log("Event2", Info)
	logCollector.Log("Event3", Info)
	log := logCollector.GetLogs()
	assert.Equal(t, 2, len(log))
	assert.Equal(t, "Event2", log[0].Log)
	assert.Equal(t, "Event3", log[1].Log)
}
func TestLogCollectorTimeCompaction(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	logCollector := NewLogCollector()
	start := time.Now()
	later := start.Add(10 * time.Minute)
	end := start.Add(20 * time.Minute)
	logCollector.Log("Event1", Info)
	logCollector.Log("Event2", Info)
	logCollector.Log("Event3", Info)
	logCollector.store[0].Timestamp = start
	logCollector.store[1].Timestamp = start
	logCollector.store[2].Timestamp = later
	logCollector.compact(end)
	log := logCollector.GetLogs()
	assert.Equal(t, 1, len(log))
	assert.Equal(t, LogItem{Log: "Event3", Timestamp: later, Level: Info}, log[0])
}
