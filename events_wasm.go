// events_wasm.go provides a mechanism to get orientation events from JavaScript in a WASM environment
package main

import (
	"syscall/js"
)

// pushOrientationEvent receives orientation data from JavaScript
func pushOrientationEvent(this js.Value, p []js.Value) interface{} {
	if len(p) != 3 {
		return nil
	}

	event := OrientationEvent{
		Alpha: p[0].Float(),
		Beta:  p[1].Float(),
		Gamma: p[2].Float(),
	}

	// Non-blocking send to channel
	select {
	case orientationChannel <- event:
		// Successfully sent
	default:
		// Channel is full, skip this event
	}

	return nil
}

// GetOrientationEvent returns the latest orientation event if available
func GetOrientationEvent() (OrientationEvent, bool) {
	select {
	case event := <-orientationChannel:
		return event, true
	default:
		return OrientationEvent{}, false
	}
}

func init() {
	js.Global().Set("pushOrientationEvent", js.FuncOf(pushOrientationEvent))
}
