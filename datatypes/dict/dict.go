package dict

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/multierr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// revive:disable:confusing-naming  // 泛型不同类型，会误报

// Dict 是一个普通字典，支持并发
type Dict[T comparable, V any] struct {
	lock *sync.RWMutex
	data map[T]V
}

func NewDict[T comparable, V any](initCapacity int) *Dict[T, V] {
	return &Dict[T, V]{
		lock: &sync.RWMutex{},
		data: make(map[T]V, initCapacity),
	}
}

func (d Dict[T, V]) Get(key T) (info V, exist bool) {
	return d.getWithLock(key, true)
}

func (d Dict[T, V]) getWithLock(key T, needLock bool) (info V, exist bool) {
	if needLock {
		d.lock.RLock()

		defer d.lock.RUnlock()
	}

	info, exist = d.data[key]

	return info, exist
}

func (d *Dict[T, V]) Set(key T, value V) {
	d.setWithLock(key, value, true)
}

func (d *Dict[T, V]) setWithLock(key T, value V, needLock bool) {
	if needLock {
		d.lock.Lock()
		defer d.lock.Unlock()
	}

	d.data[key] = value
}

func (d *Dict[T, V]) SetWhen(key T, value V, match func(oldValue V) bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	old, exist := d.data[key]

	if !exist || match(old) {
		d.data[key] = value
	}
}
func (d *Dict[T, V]) SetWithoutLock(key T, value V) {
	d.data[key] = value
}

func (d *Dict[T, V]) SetDefault(keys []T, value V) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, key := range keys {
		d.data[key] = value
	}
}

func (d *Dict[T, V]) SetNX(key T, value V) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, exist := d.data[key]; exist {
		return false
	}

	d.data[key] = value

	return true
}

func (d *Dict[T, V]) ForRange(process func(key T, value V)) {
	if d == nil {
		return
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	for key, value := range d.data {
		process(key, value)
	}
}

func (d *Dict[T, V]) Slice() []T {
	result := make([]T, 0, len(d.data))

	d.lock.RLock()
	defer d.lock.RUnlock()

	for key := range d.data {
		result = append(result, key)
	}

	return result
}
func (d *Dict[T, V]) ForRangeConcurrent(process func(key T, value V)) {
	if d == nil {
		return
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	wg := &sync.WaitGroup{}

	for key, value := range d.data {
		wg.Add(1)

		go func(loopKey T, loopValue V) {
			defer wg.Done()

			process(loopKey, loopValue)
		}(key, value)
	}

	wg.Wait()
}

func (d *Dict[T, V]) ForRangeConcurrentOnlyValue(process func(value V), wgPool *Pool[*sync.WaitGroup]) {
	if d == nil {
		return
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	wg := wgPool.Get()

	for _, value := range d.data {
		wg.Add(1)

		go func(loopValue V) {
			defer wg.Done()

			process(loopValue)
		}(value)
	}

	wg.Wait()

	wgPool.Put(wg)
}

func (d *Dict[T, V]) ForRangeWithSet(process func(key T, value V)) {
	if d == nil {
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	for key, value := range d.data {
		process(key, value)
		d.data[key] = value
	}
}

func (d *Dict[T, V]) Do(key T, do func(key T, value V) V) {
	if d == nil {
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	value, exist := d.data[key]

	if !exist {
		return
	}

	newValue := do(key, value)

	d.data[key] = newValue
}

func (d *Dict[T, V]) DeleteAll() {
	d.deleteAllWithLock(true)
}

func (d *Dict[T, V]) deleteAllWithLock(needLock bool) {
	if needLock {
		d.lock.Lock()
		defer d.lock.Unlock()
	}

	for key := range d.data {
		delete(d.data, key)
	}
}

func (d *Dict[T, V]) Len() int {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return len(d.data)
}

func (d *Dict[T, V]) ForRangeWithError(process func(key T, value V) error) error {
	if d == nil {
		return nil
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	var (
		singleErr, err error
	)

	for key, value := range d.data {
		if singleErr = process(key, value); singleErr != nil {
			err = multierr.Append(err, singleErr)
		}
	}

	return err
}

func (d *Dict[T, V]) ForRangeWithErrorConcurrent(process func(key T, value V) error) error {
	if d == nil {
		return nil
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	var (
		singleErr, err error
		errCh          = make(chan error, len(d.data))
		wg             = &sync.WaitGroup{}
	)

	for key, value := range d.data {
		wg.Add(1)

		go func(loopKey T, loopValue V) {
			defer wg.Done()

			singleErr = process(loopKey, loopValue)

			errCh <- singleErr
		}(key, value)
	}

	wg.Wait()

	close(errCh)

	for singleErr = range errCh {
		err = multierr.Append(err, singleErr)
	}

	return err
}

func (d *Dict[T, V]) Delete(keys ...T) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.data) == 0 {
		return
	}

	for _, key := range keys {
		delete(d.data, key)
	}
}

func (d *Dict[T, V]) DeleteNonExist(keys ...T) {
	d.lock.Lock()
	defer d.lock.Unlock()

	deleteKeys := make([]T, 0, len(d.data))

	keysMap := make(map[T]struct{}, len(keys))

	for i := range keys {
		keysMap[keys[i]] = struct{}{}
	}

	for key := range d.data {
		if _, exist := keysMap[key]; !exist {
			deleteKeys = append(deleteKeys, key)
		}
	}

	for _, key := range deleteKeys {
		delete(d.data, key)
	}
}
func (d *Dict[T, V]) DeleteWithCheck(key T) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.data) == 0 {
		return false
	}

	if _, exist := d.data[key]; exist {
		delete(d.data, key)
		return true
	}

	return false
}
func (d *Dict[T, V]) DeleteAndLen(key T) (length int, deleted bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.data) == 0 {
		return 0, false
	}

	if _, exist := d.data[key]; exist {
		delete(d.data, key)
		return len(d.data), true
	}

	return len(d.data), false
}

func (d Dict[T, V]) Copy(fields ...T) map[T]V {
	var (
		result     = make(map[T]V, d.Len())
		fieldsDict = NewDictFromSlice(fields)
	)

	d.ForRange(func(key T, value V) {
		if _, exist := fieldsDict.Get(key); exist || len(fields) == 0 {
			result[key] = value
		}
	})

	return result
}

func (d Dict[T, V]) GormDataType() string {
	return jsonType
}

func (d Dict[T, V]) Value() (driver.Value, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return json.Marshal(d.data)
}

func (d *Dict[T, V]) Scan(data any) error {
	if d.lock == nil {
		d.lock = &sync.RWMutex{}
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.data == nil {
		d.data = map[T]V{}
	}

	return ScanFromGorm(data, &d.data)
}

func (d Dict[T, V]) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	d.lock.RLock()
	defer d.lock.RUnlock()

	data, _ := json.Marshal(d.data) //nolint:errchkjson //没有空间返回

	if db.Dialector.Name() == `mysql` {
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

func (d *Dict[T, V]) Reset(data map[T]V) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.deleteAllWithLock(false)

	for a := range data {
		d.setWithLock(a, data[a], false)
	}
}

func ScanFromGorm(data any, value any) error {
	switch d := data.(type) {
	case []byte:
		return json.Unmarshal(d, value)
	default:
		return fmt.Errorf(`不支持的类型[%s]`, reflect.TypeOf(data).String())
	}
}

func NewDictFromSlice[T comparable](slice []T) *Dict[T, struct{}] {
	result := NewDict[T, struct{}](len(slice))

	for i := range slice {
		result.data[slice[i]] = struct{}{}
	}

	return result
}

func NewDictFromContainerSlice[T comparable, V any](slice []V, get func(V) T) *Dict[T, struct{}] {
	result := NewDict[T, struct{}](len(slice))

	for i := range slice {
		result.data[get(slice[i])] = struct{}{}
	}

	return result
}

type Counter[T comparable] Dict[T, int]

func NewCounter[T comparable](_ T, initCapacity int) *Counter[T] {
	data := Counter[T](*NewDict[T, int](initCapacity))

	return &data
}

func (c *Counter[T]) Add(t T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[t] += 1
}

func (c *Counter[T]) Remove(t T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[t] -= 1

	if c.data[t] <= 0 {
		delete(c.data, t)
	}
}

func (c *Counter[T]) Slice() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()

	result := make([]T, 0, len(c.data))

	for key := range c.data {
		result = append(result, key)
	}

	return result
}

const (
	jsonType = `JSON`
)

type Pool[T any] struct {
	data *sync.Pool
}

func NewPool[T any](initCapacity int, newFunc func() T) *Pool[T] {
	result := &Pool[T]{
		data: &sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
	}

	for i := 0; i < initCapacity; i++ {
		result.data.Put(newFunc())
	}

	return result
}

func (p *Pool[T]) Get() T {
	return p.data.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	p.data.Put(t)
}
