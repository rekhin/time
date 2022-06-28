package nats

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var (
			expected = time.Millisecond * 5
			actual   time.Duration
			wg       sync.WaitGroup
		)
		wg.Add(1)
		start := time.Now()
		_ = AfterFunc(time.Millisecond*5, func() {
			actual = time.Since(start)
			wg.Done()
		})
		wg.Wait()
		assert.Equal(t, expected.Milliseconds(), actual.Milliseconds(), "they should be equal")
	})
	t.Run("with reset", func(t *testing.T) {
		var (
			expected = time.Millisecond * 15
			actual   time.Duration
			wg       sync.WaitGroup
		)
		wg.Add(1)
		start := time.Now()
		timer := AfterFunc(time.Millisecond*10, func() {
			actual = time.Since(start)
			wg.Done()
		})
		time.Sleep(time.Millisecond * 5)
		timer.Reset()
		wg.Wait()
		assert.Equal(t, expected.Milliseconds(), actual.Milliseconds(), "they should be equal")
	})
	t.Run("with double reset", func(t *testing.T) {
		var (
			expected = time.Millisecond * 20
			actual   time.Duration
			wg       sync.WaitGroup
		)
		wg.Add(1)
		start := time.Now()
		timer := AfterFunc(time.Millisecond*10, func() {
			actual = time.Since(start)
			assert.Equal(t, expected.Milliseconds(), actual.Milliseconds(), "they should be equal")
			wg.Done()
		})
		time.Sleep(time.Millisecond * 5)
		timer.Reset()
		time.Sleep(time.Millisecond * 5)
		timer.Reset()
		wg.Wait()
	})
	t.Run("with reset after trigger", func(t *testing.T) {
		var (
			expectedList = []time.Duration{time.Millisecond * 5, time.Millisecond * 15}
			actualList   []time.Duration
			wg           sync.WaitGroup
		)
		wg.Add(len(expectedList))
		start := time.Now()
		timer := AfterFunc(time.Millisecond*5, func() {
			actualList = append(actualList, time.Since(start))
			wg.Done()
		})
		time.Sleep(time.Millisecond * 10)
		timer.Reset()
		wg.Wait()
		for i, expected := range expectedList {
			assert.Equal(t, expected.Milliseconds(), actualList[i].Milliseconds(), "they should be equal")
		}
	})
}
