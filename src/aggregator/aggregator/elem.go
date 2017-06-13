// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aggregator

import (
	"errors"
	"sync"
	"time"

	"github.com/m3db/m3aggregator/aggregation"
	"github.com/m3db/m3metrics/metric/id"
	"github.com/m3db/m3metrics/metric/unaggregated"
	"github.com/m3db/m3metrics/policy"
)

const (
	defaultNumValues = 2
)

var (
	emptyTimedCounter timedCounter
	emptyTimedTimer   timedTimer
	emptyTimedGauge   timedGauge

	errCounterElemClosed = errors.New("counter element is closed")
	errTimerElemClosed   = errors.New("timer element is closed")
	errGaugeElemClosed   = errors.New("gauge element is closed")
)

type aggMetricFn func(
	idPrefix []byte,
	id id.RawID,
	idSuffix []byte,
	timeNanos int64,
	value float64,
	sp policy.StoragePolicy,
)

type newMetricElemFn func() metricElem

// metricElem is the common interface for metric elements
type metricElem interface {
	// ID returns the metric id.
	ID() id.RawID

	// ResetSetData resets the element and sets data
	ResetSetData(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes)

	// AddMetric adds a new metric value
	// TODO(xichen): a value union would suffice here
	AddMetric(timestamp time.Time, mu unaggregated.MetricUnion) error

	// Consume processes values before a given time and discards
	// them afterwards, returning whether the element can be collected
	// after discarding the values
	Consume(earlierThanNanos int64, fn aggMetricFn) bool

	// MarkAsTombstoned marks an element as tombstoned, which means this element
	// will be deleted once its aggregated values have been flushed
	MarkAsTombstoned()

	// Close closes the element
	Close()
}

type timedCounter struct {
	timeNanos int64
	counter   *aggregation.LockedCounter
}

type timedTimer struct {
	timeNanos int64
	timer     *aggregation.LockedTimer
}

type timedGauge struct {
	timeNanos int64
	gauge     *aggregation.LockedGauge
}

type elemBase struct {
	sync.RWMutex

	opts       Options
	id         id.RawID
	sp         policy.StoragePolicy
	aggTypes   policy.AggregationTypes
	aggOpts    aggregation.Options
	tombstoned bool
	closed     bool
}

// CounterElem is the counter element
type CounterElem struct {
	elemBase

	values    []timedCounter // aggregated counters sorted by time in ascending order
	toConsume []timedCounter
}

// TimerElem is the timer element
type TimerElem struct {
	elemBase

	isAggTypesPooled  bool
	isQuantilesPooled bool
	quantiles         []float64

	values    []timedTimer // aggregated timers sorted by time in ascending order
	toConsume []timedTimer
}

// GaugeElem is the gauge element
type GaugeElem struct {
	elemBase

	values    []timedGauge // aggregated gauges sorted by time in ascending order
	toConsume []timedGauge
}

func newElemBase(opts Options) elemBase {
	return elemBase{
		opts:    opts,
		aggOpts: aggregation.NewOptions(),
	}
}

// ResetSetData resets the counter and sets data.
func (e *elemBase) ResetSetData(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes) {
	e.id = id
	e.sp = sp
	e.aggTypes = aggTypes
	e.aggOpts.ResetSetData(aggTypes)
	e.tombstoned = false
	e.closed = false
}

// ID returns the metric id.
func (e *elemBase) ID() id.RawID {
	e.RLock()
	id := e.id
	e.RUnlock()
	return id
}

// MarkAsTombstoned marks an element as tombstoned, which means this element
// will be deleted once its aggregated values have been flushed.
func (e *elemBase) MarkAsTombstoned() {
	e.Lock()
	if e.closed {
		e.Unlock()
		return
	}
	e.tombstoned = true
	e.Unlock()
}

// NewCounterElem creates a new counter element
func NewCounterElem(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes, opts Options) *CounterElem {
	e := &CounterElem{
		elemBase: newElemBase(opts),
		values:   make([]timedCounter, 0, defaultNumValues), // in most cases values will have two entries
	}
	e.ResetSetData(id, sp, aggTypes)
	return e
}

// AddMetric adds a new counter value
func (e *CounterElem) AddMetric(timestamp time.Time, mu unaggregated.MetricUnion) error {
	alignedStart := timestamp.Truncate(e.sp.Resolution().Window).UnixNano()
	counter, err := e.findOrCreate(alignedStart)
	if err != nil {
		return err
	}
	counter.Lock()
	counter.Update(mu.CounterVal)
	counter.Unlock()
	return nil
}

// Consume processes values before a given time and discards
// them afterwards, returning whether the element can be collected
// after discarding the values.
func (e *CounterElem) Consume(earlierThanNanos int64, fn aggMetricFn) bool {
	e.Lock()
	if e.closed {
		e.Unlock()
		return false
	}
	idx := 0
	for range e.values {
		// Bail as soon as the timestamp is no later than the target time
		if e.values[idx].timeNanos >= earlierThanNanos {
			break
		}
		idx++
	}
	e.toConsume = e.toConsume[:0]
	if idx > 0 {
		// Shift remaining values to the left and shrink the values slice
		e.toConsume = append(e.toConsume, e.values[:idx]...)
		n := copy(e.values[0:], e.values[idx:])
		e.values = e.values[:n]
	}
	canCollect := len(e.values) == 0 && e.tombstoned
	e.Unlock()

	// Process the counters.
	for i := range e.toConsume {
		endAtNanos := e.toConsume[i].timeNanos + int64(e.sp.Resolution().Window)
		e.toConsume[i].counter.Lock()
		e.processValue(endAtNanos, e.toConsume[i].counter.Counter, fn)
		e.toConsume[i].counter.Unlock()
		e.toConsume[i] = emptyTimedCounter
	}

	return canCollect
}

// Close closes the counter element
func (e *CounterElem) Close() {
	e.Lock()
	if e.closed {
		e.Unlock()
		return
	}
	e.closed = true
	e.id = nil
	e.values = e.values[:0]
	e.toConsume = e.toConsume[:0]
	aggTypesPool := e.opts.AggregationTypesPool()
	pool := e.opts.CounterElemPool()
	e.Unlock()

	aggTypesPool.Put(e.aggTypes)
	pool.Put(e)
}

// findOrCreate finds the counter for a given time, or creates one
// if it doesn't exist.
func (e *CounterElem) findOrCreate(alignedStart int64) (*aggregation.LockedCounter, error) {
	e.RLock()
	if e.closed {
		e.RUnlock()
		return nil, errCounterElemClosed
	}
	idx, found := e.indexOfWithLock(alignedStart)
	if found {
		counter := e.values[idx].counter
		e.RUnlock()
		return counter, nil
	}
	e.RUnlock()

	e.Lock()
	if e.closed {
		e.Unlock()
		return nil, errCounterElemClosed
	}
	idx, found = e.indexOfWithLock(alignedStart)
	if found {
		counter := e.values[idx].counter
		e.Unlock()
		return counter, nil
	}

	// If not found, create a new counter.
	numValues := len(e.values)
	e.values = append(e.values, emptyTimedCounter)
	copy(e.values[idx+1:numValues+1], e.values[idx:numValues])
	e.values[idx] = timedCounter{
		timeNanos: alignedStart,
		counter:   aggregation.NewLockedCounter(aggregation.NewCounter(e.aggOpts)),
	}
	counter := e.values[idx].counter
	e.Unlock()
	return counter, nil
}

// indexOfWithLock finds the smallest element index whose timestamp
// is no smaller than the start time passed in, and true if it's an
// exact match, false otherwise.
func (e *CounterElem) indexOfWithLock(alignedStart int64) (int, bool) {
	numValues := len(e.values)
	// Optimize for the common case
	if numValues > 0 && e.values[numValues-1].timeNanos == alignedStart {
		return numValues - 1, true
	}
	// Binary search for the unusual case. We intentionally do not
	// use the sort.Search() function because it requires passing
	// in a closure
	left, right := 0, numValues
	for left < right {
		mid := left + (right-left)/2 // avoid overflow
		if e.values[mid].timeNanos < alignedStart {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// If the current timestamp is equal to or larger than the target time,
	// return the index as is
	if left < numValues && e.values[left].timeNanos == alignedStart {
		return left, true
	}
	return left, false
}

func (e *CounterElem) processValue(timeNanos int64, agg aggregation.Counter, fn aggMetricFn) {
	var fullCounterPrefix = e.opts.FullCounterPrefix()
	if e.aggOpts.UseDefaultAggregation {
		// No suffix for default aggregations.
		// TODO(cw) Make aggregation types for Counter configurable.
		fn(fullCounterPrefix, e.id, e.suffix(policy.Sum), timeNanos, float64(agg.Sum()), e.sp)
		return
	}

	for _, aggType := range e.aggTypes {
		fn(fullCounterPrefix, e.id, e.suffix(aggType), timeNanos, agg.ValueOf(aggType), e.sp)
	}
}

// NB(cw) suffix overrides default suffixes to not add
// suffix for default Counter metrics for backward compatibility reasons.
func (e *CounterElem) suffix(aggType policy.AggregationType) []byte {
	if aggType == policy.Sum {
		return nil
	}
	return e.opts.Suffix(aggType)
}

// NewTimerElem creates a new timer element.
func NewTimerElem(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes, opts Options) *TimerElem {
	e := &TimerElem{
		elemBase: newElemBase(opts),
		values:   make([]timedTimer, 0, defaultNumValues), // in most cases values will have two entries
	}
	e.ResetSetData(id, sp, aggTypes)
	return e
}

// ResetSetData overrides elemBase.ResetData, to include quantile update for timer elements.
func (e *TimerElem) ResetSetData(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes) {
	if aggTypes.IsDefault() {
		aggTypes = e.opts.DefaultTimerAggregationTypes()
		e.isAggTypesPooled = false
		e.quantiles, e.isQuantilesPooled = e.opts.TimerQuantiles(), false
	} else {
		e.isAggTypesPooled = true
		e.quantiles, e.isQuantilesPooled = aggTypes.PooledQuantiles(e.opts.QuantilesPool())
	}

	e.elemBase.ResetSetData(id, sp, aggTypes)
}

// AddMetric adds a new batch of timer values.
func (e *TimerElem) AddMetric(timestamp time.Time, mu unaggregated.MetricUnion) error {
	alignedStart := timestamp.Truncate(e.sp.Resolution().Window).UnixNano()
	timer, err := e.findOrCreate(alignedStart)
	if err != nil {
		return err
	}
	timer.Lock()
	timer.AddBatch(mu.BatchTimerVal)
	timer.Unlock()
	return nil
}

// Consume processes values before a given time and discards
// them afterwards, returning whether the element can be collected
// after discarding the values.
func (e *TimerElem) Consume(earlierThanNanos int64, fn aggMetricFn) bool {
	e.Lock()
	if e.closed {
		e.Unlock()
		return false
	}
	idx := 0
	for range e.values {
		// Bail as soon as the timestamp is no later than the target time.
		if e.values[idx].timeNanos >= earlierThanNanos {
			break
		}
		idx++
	}
	e.toConsume = e.toConsume[:0]
	if idx > 0 {
		// Shift remaining values to the left and shrink the values slice.
		e.toConsume = append(e.toConsume, e.values[:idx]...)
		n := copy(e.values[0:], e.values[idx:])
		// Clear out the invalid timer items to avoid holding onto the underlying streams.
		for i := n; i < len(e.values); i++ {
			e.values[i] = emptyTimedTimer
		}
		e.values = e.values[:n]
	}
	canCollect := len(e.values) == 0 && e.tombstoned
	e.Unlock()

	// Process the timers.
	for i := range e.toConsume {
		endAtNanos := e.toConsume[i].timeNanos + int64(e.sp.Resolution().Window)
		e.toConsume[i].timer.Lock()
		e.processValue(endAtNanos, e.toConsume[i].timer.Timer, fn)
		e.toConsume[i].timer.Unlock()
		// Close the timer after it's processed.
		e.toConsume[i].timer.Close()
		e.toConsume[i] = emptyTimedTimer
	}

	return canCollect
}

// Close closes the timer element.
func (e *TimerElem) Close() {
	e.Lock()
	if e.closed {
		e.Unlock()
		return
	}
	e.closed = true
	e.id = nil
	for idx := range e.values {
		// Returning the underlying stream to pool
		e.values[idx].timer.Close()
		e.values[idx] = emptyTimedTimer
	}
	e.values = e.values[:0]
	e.toConsume = e.toConsume[:0]
	quantileFloatsPool := e.opts.QuantilesPool()
	aggTypesPool := e.opts.AggregationTypesPool()
	pool := e.opts.TimerElemPool()
	e.Unlock()

	if e.isQuantilesPooled {
		quantileFloatsPool.Put(e.quantiles)
	}
	if e.isAggTypesPooled {
		aggTypesPool.Put(e.aggTypes)
	}
	pool.Put(e)
}

// findOrCreate finds the timer for a given time, or creates one
// if it doesn't exist.
func (e *TimerElem) findOrCreate(alignedStart int64) (*aggregation.LockedTimer, error) {
	e.RLock()
	if e.closed {
		e.RUnlock()
		return nil, errTimerElemClosed
	}
	idx, found := e.indexOfWithLock(alignedStart)
	if found {
		timer := e.values[idx].timer
		e.RUnlock()
		return timer, nil
	}
	e.RUnlock()

	e.Lock()
	if e.closed {
		e.Unlock()
		return nil, errTimerElemClosed
	}
	idx, found = e.indexOfWithLock(alignedStart)
	if found {
		timer := e.values[idx].timer
		e.Unlock()
		return timer, nil
	}

	// If not found, create a new timer.
	numValues := len(e.values)
	e.values = append(e.values, emptyTimedTimer)
	copy(e.values[idx+1:numValues+1], e.values[idx:numValues])
	newTimer := aggregation.NewTimer(e.opts.TimerQuantiles(), e.opts.StreamOptions(), e.aggOpts)
	e.values[idx] = timedTimer{
		timeNanos: alignedStart,
		timer:     aggregation.NewLockedTimer(newTimer),
	}
	timer := e.values[idx].timer
	e.Unlock()
	return timer, nil
}

// indexOfWithLock finds the smallest element index whose timestamp
// is no smaller than the start time passed in, and true if it's an
// exact match, false otherwise.
func (e *TimerElem) indexOfWithLock(alignedStart int64) (int, bool) {
	numValues := len(e.values)
	// Optimize for the common case
	if numValues > 0 && e.values[numValues-1].timeNanos == alignedStart {
		return numValues - 1, true
	}
	// Binary search for the unusual case. We intentionally do not
	// use the sort.Search() function because it requires passing
	// in a closure
	left, right := 0, numValues
	for left < right {
		mid := left + (right-left)/2 // avoid overflow
		if e.values[mid].timeNanos < alignedStart {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// If the current timestamp is equal to or larger than the target time,
	// return the index as is
	if left < numValues && e.values[left].timeNanos == alignedStart {
		return left, true
	}
	return left, false
}

func (e *TimerElem) processValue(timeNanos int64, agg aggregation.Timer, fn aggMetricFn) {
	fullTimerPrefix := e.opts.FullTimerPrefix()
	if e.aggOpts.UseDefaultAggregation {
		// NB(cw) Use default suffix slice for faster look up.
		suffixes := e.opts.DefaultTimerAggregationSuffixes()
		aggTypes := e.opts.DefaultTimerAggregationTypes()
		for i, aggType := range aggTypes {
			fn(fullTimerPrefix, e.id, suffixes[i], timeNanos, agg.ValueOf(aggType), e.sp)
		}
		return
	}

	for _, aggType := range e.aggTypes {
		fn(fullTimerPrefix, e.id, e.opts.Suffix(aggType), timeNanos, agg.ValueOf(aggType), e.sp)
	}
}

// NewGaugeElem creates a new gauge element
func NewGaugeElem(id id.RawID, sp policy.StoragePolicy, aggTypes policy.AggregationTypes, opts Options) *GaugeElem {
	e := &GaugeElem{
		elemBase: newElemBase(opts),
		values:   make([]timedGauge, 0, defaultNumValues), // in most cases values will have two entries
	}
	e.ResetSetData(id, sp, aggTypes)
	return e
}

// AddMetric adds a new gauge value
func (e *GaugeElem) AddMetric(timestamp time.Time, mu unaggregated.MetricUnion) error {
	alignedStart := timestamp.Truncate(e.sp.Resolution().Window).UnixNano()
	gauge, err := e.findOrCreate(alignedStart)
	if err != nil {
		return err
	}
	gauge.Lock()
	gauge.Update(mu.GaugeVal)
	gauge.Unlock()
	return nil
}

// Consume processes values before a given time and discards
// them afterwards, returning whether the element can be collected
// after discarding the values
func (e *GaugeElem) Consume(earlierThanNanos int64, fn aggMetricFn) bool {
	e.Lock()
	if e.closed {
		e.Unlock()
		return false
	}
	idx := 0
	for range e.values {
		// Bail as soon as the timestamp is no later than the target time
		if e.values[idx].timeNanos >= earlierThanNanos {
			break
		}
		idx++
	}
	e.toConsume = e.toConsume[:0]
	if idx > 0 {
		// Shift remaining values to the left and shrink the values slice
		e.toConsume = append(e.toConsume, e.values[:idx]...)
		n := copy(e.values[0:], e.values[idx:])
		e.values = e.values[:n]
	}
	canCollect := len(e.values) == 0 && e.tombstoned
	e.Unlock()

	for i := range e.toConsume {
		endAtNanos := e.toConsume[i].timeNanos + int64(e.sp.Resolution().Window)
		e.toConsume[i].gauge.Lock()
		e.processValue(endAtNanos, e.toConsume[i].gauge.Gauge, fn)
		e.toConsume[i].gauge.Unlock()
		e.toConsume[i] = emptyTimedGauge
	}

	return canCollect
}

// Close closes the gauge element
func (e *GaugeElem) Close() {
	e.Lock()
	if e.closed {
		e.Unlock()
		return
	}
	e.closed = true
	e.id = nil
	e.values = e.values[:0]
	e.toConsume = e.toConsume[:0]
	aggTypesPool := e.opts.AggregationTypesPool()
	pool := e.opts.GaugeElemPool()
	e.Unlock()

	aggTypesPool.Put(e.aggTypes)
	pool.Put(e)
}

// findOrCreate finds the gauge for a given time, or creates one
// if it doesn't exist.
func (e *GaugeElem) findOrCreate(alignedStart int64) (*aggregation.LockedGauge, error) {
	e.RLock()
	if e.closed {
		e.RUnlock()
		return nil, errGaugeElemClosed
	}
	idx, found := e.indexOfWithLock(alignedStart)
	if found {
		gauge := e.values[idx].gauge
		e.RUnlock()
		return gauge, nil
	}
	e.RUnlock()

	e.Lock()
	if e.closed {
		e.Unlock()
		return nil, errGaugeElemClosed
	}
	idx, found = e.indexOfWithLock(alignedStart)
	if found {
		gauge := e.values[idx].gauge
		e.Unlock()
		return gauge, nil
	}

	// If not found, create a new gauge.
	numValues := len(e.values)
	e.values = append(e.values, emptyTimedGauge)
	copy(e.values[idx+1:numValues+1], e.values[idx:numValues])
	e.values[idx] = timedGauge{
		timeNanos: alignedStart,
		gauge:     aggregation.NewLockedGauge(aggregation.NewGauge(e.aggOpts)),
	}
	gauge := e.values[idx].gauge
	e.Unlock()
	return gauge, nil
}

// indexOfWithLock finds the smallest element index whose timestamp
// is no smaller than the start time passed in, and true if it's an
// exact match, false otherwise.
func (e *GaugeElem) indexOfWithLock(alignedStart int64) (int, bool) {
	numValues := len(e.values)
	// Optimize for the common case
	if numValues > 0 && e.values[numValues-1].timeNanos == alignedStart {
		return numValues - 1, true
	}
	// Binary search for the unusual case. We intentionally do not
	// use the sort.Search() function because it requires passing
	// in a closure
	left, right := 0, numValues
	for left < right {
		mid := left + (right-left)/2 // avoid overflow
		if e.values[mid].timeNanos < alignedStart {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// If the current timestamp is equal to or larger than the target time,
	// return the index as is
	if left < numValues && e.values[left].timeNanos == alignedStart {
		return left, true
	}
	return left, false
}

func (e *GaugeElem) processValue(timeNanos int64, agg aggregation.Gauge, fn aggMetricFn) {
	var fullGaugePrefix = e.opts.FullGaugePrefix()
	if e.aggOpts.UseDefaultAggregation {
		// No suffix for default aggregations.
		// TODO(cw) Make aggregation types for Gauge configurable.
		fn(fullGaugePrefix, e.id, e.suffix(policy.Last), timeNanos, agg.Last(), e.sp)
		return
	}

	for _, aggType := range e.aggTypes {
		fn(fullGaugePrefix, e.id, e.suffix(aggType), timeNanos, agg.ValueOf(aggType), e.sp)
	}
}

// NB(cw) suffix overrides default suffixes to not add
// suffix for default Gauge metrics for backward compatibility reasons.
func (e *GaugeElem) suffix(aggType policy.AggregationType) []byte {
	if aggType == policy.Last {
		return nil
	}
	return e.opts.Suffix(aggType)
}