package decoders

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewQueue(t *testing.T) {
	regularpacket := "Af2|3|0||-1"
	emptypacket := "Af"

	queue, err := NewQueue(regularpacket)
	assert.NilError(t, err)

	queue, err = NewQueue(emptypacket)
	assert.Assert(t, queue == nil)
	assert.Assert(t, err != nil)

	queue, err = NewQueue("")
	assert.Assert(t, queue == nil)
	assert.Assert(t, err != nil)
}

func TestUpdateQueuePosition(t *testing.T) {
	regularpacket := "Af2|3|0||-1"
	emptypacket := "Af"

	queue, err := NewQueue(regularpacket)
	assert.NilError(t, err)
	err = queue.UpdateQueuePosition(regularpacket)
	assert.Assert(t, queue.currentPos == 2)
	assert.Assert(t, queue.totalSub == 3)
	assert.Assert(t, !queue.IsSub)

	queue, err = NewQueue(emptypacket)
	assert.Assert(t, err != nil)
	err = queue.UpdateQueuePosition(emptypacket)
	assert.Assert(t, queue == nil)
}
