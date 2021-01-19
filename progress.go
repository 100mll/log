// Copyright (c) 2018-2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"strings"
	"sync"
	"time"
)

var ProgressInterval = 10 * time.Second

type ProgressLogger struct {
	name        string
	event       string
	calls       int64
	events      int64
	lastLogTime time.Time
	logger      Logger
	sync.Mutex
}

func NewProgressLogger(name, event string, logger Logger) *ProgressLogger {
	if logger == nil {
		logger = Log
	}
	return &ProgressLogger{
		name:        name,
		event:       event,
		lastLogTime: time.Now(),
		logger:      logger,
	}
}

func pluralize(str string, count int64) string {
	if str == "" {
		str = "call"
	}
	if count == 0 || count > 1 {
		str += "s"
	}
	return str
}

func (p *ProgressLogger) LogN(n int) {
	p.Log(n, time.Time{})
}

func (p *ProgressLogger) Log(n int, ts time.Time, extra ...string) {
	p.Lock()
	defer p.Unlock()
	p.calls++
	p.events += int64(n)
	now := time.Now()
	duration := now.Sub(p.lastLogTime)
	if duration < ProgressInterval || p.events == 0 {
		return
	}

	// Truncate the duration to 10s of milliseconds.
	tDuration := duration.Truncate(10 * time.Millisecond)

	// Log information about the event.
	eventStr := p.event
	if p.events == 1 {
		eventStr = eventStr[:len(eventStr)-1]
	}
	callType := ""
	if len(extra) > 0 {
		callType = extra[0]
		extra = extra[1:]
	}
	callStr := pluralize(callType, p.calls)
	extraString := ""
	if ex := strings.Join(extra, " "); len(ex) > 0 {
		extraString = ", " + ex
	}
	tm := ts.UTC().Format("2006-01-02 15:04:05 MST ")
	if ts.IsZero() {
		tm = ""
	}
	p.logger.Infof("%s: processed %d %sin %s (%s, %d %s%s)",
		p.name, p.events, eventStr, tDuration,
		tm, p.calls, callStr, extraString)

	p.calls = 0
	p.events = 0
	p.lastLogTime = now
}
