package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"JobQueueGo/utils/types"
)

func TestJobQueue(t *testing.T) {
	queue := NewJobQueue()

	job1 := types.Job{Id: "job1"}
	job2 := types.Job{Id: "job2"}
	job3 := types.Job{Id: "job3"}

	// Lägg till jobb i kön
	queue.Enqueue(job1)
	queue.Enqueue(job2)
	queue.Enqueue(job3)

	// Kontrollera längd
	assert.Equal(t, 3, queue.Length())

	// Plocka ut första jobbet
	j, ok := queue.Dequeue()
	require.True(t, ok)
	assert.Equal(t, "job1", j.Id)
	assert.Equal(t, 2, queue.Length())

	// Plocka ut andra
	j, ok = queue.Dequeue()
	require.True(t, ok)
	assert.Equal(t, "job2", j.Id)
	assert.Equal(t, 1, queue.Length())

	// Plocka ut tredje
	j, ok = queue.Dequeue()
	require.True(t, ok)
	assert.Equal(t, "job3", j.Id)
	assert.Equal(t, 0, queue.Length())

	// Testa att det nu är tomt
	_, ok = queue.Dequeue()
	require.False(t, ok)
}

