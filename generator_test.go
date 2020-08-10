package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestNewGenerator(t *testing.T) {
	outputCh := make(chan interface{})
	defer close(outputCh)

	handler := func(sink GeneratorSink, generator *Generator) {
		defer generator.Close(nil)
		sink <- 1
		sink <- 2
	}

	generator := NewGenerator(outputCh, handler)
	require.NotNil(t, generator)
	var sink <-chan interface{} = outputCh
	assert.Equal(t, sink, generator.OutputCh)
	assert.False(t, generator.cancelled)
	assert.Nil(t, generator.err)

	var something = 0

	go func() {
		something = (<-outputCh).(int)
	}()
	assert.Eventually(t, func() bool { return something == 1 }, time.Millisecond*20, time.Microsecond)

	go func() {
		something = (<-outputCh).(int)
	}()
	assert.Eventually(t, func() bool { return something == 2 }, time.Millisecond*20, time.Microsecond)

	assert.False(t, generator.cancelled)
	assert.Nil(t, generator.err)
}

func TestGenerator_Cancel(t *testing.T) {
	// Test that calling cancel sends a message to the cancel channel, by
	// having the generator forward anything from the cancel channel to output.
	generator := Generator{}
	generator.wg.Add(1) // make sure cancel is not blocked by the wg
	assert.False(t, generator.cancelled)
	generator.Cancel(io.EOF)
	assert.True(t, generator.cancelled)
	assert.Equal(t, io.EOF, generator.err)
	generator.Cancel(nil)
	assert.True(t, generator.cancelled)
	assert.Equal(t, io.EOF, generator.err)
}

func TestGenerator_Cancelled(t *testing.T) {
	generator := Generator{}
	generator.wg.Add(1) // make sure cancelled is not blocked by the wg
	assert.False(t, generator.Cancelled())
	generator.cancelled = true
	assert.True(t, generator.Cancelled())
}

func TestGenerator_Close(t *testing.T) {
	generator := Generator{}
	generator.wg.Add(1)
	assert.Nil(t, generator.err)

	// Check that the wg was Done.
	var closed = false
	go func() {
		generator.wg.Wait()
		closed = true
	}()

	// Make sure we block at all
	assert.Never(t, func() bool { return closed }, time.Millisecond*10, time.Millisecond)

	// Now we should unblock
	generator.Close(io.EOF)
	assert.Eventually(t, func() bool { return closed }, time.Millisecond*20, time.Microsecond)
	assert.Equal(t, io.EOF, generator.err)

	// Close also shouldn't erase an existing error
	generator.wg.Add(1)
	generator.Close(nil)
	assert.Equal(t, io.EOF, generator.err)
}

func TestGenerator_Error(t *testing.T) {
	// Test that Error() returns what is in the generator's error but only after close.
	generator := Generator{}
	generator.wg.Add(1)
	assert.Nil(t, generator.err)

	var err error
	go func() {
		err = generator.Error()
	}()

	// Nothing should happen to err yet
	assert.Never(t, func() bool { return err != nil }, time.Millisecond*10, time.Millisecond)

	// Calling close should populate err
	generator.Close(io.EOF)
	assert.Eventually(t, func() bool { return err == io.EOF }, time.Millisecond*20, time.Microsecond)
}

func TestGenerator_setError(t *testing.T) {
	generator := Generator{}
	assert.Nil(t, generator.err)
	generator.setError(nil)
	assert.Nil(t, generator.err)
	generator.setError(io.EOF)
	assert.Equal(t, io.EOF, generator.err)
	// Verify existing error not overwritten
	generator.setError(io.ErrNoProgress)
	assert.Equal(t, io.EOF, generator.err)
	// Verify nil error does not clear existing error.
	generator.setError(nil)
	assert.Equal(t, io.EOF, generator.err)

}
