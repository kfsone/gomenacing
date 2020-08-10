// Utilities for general help writing unit tests.
package main

// Returns the first waiting value from a channel or false.
func maybeReadChannel(source <-chan interface{}) interface{} {
	select {
	case value := <-source:
		return value
	default:
		return false
	}
}
