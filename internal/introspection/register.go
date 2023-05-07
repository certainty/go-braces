package introspection

import "encoding/gob"

// Required for gob to work properly
// We should find a single place for them somehow
func RegisterTypes() {
	gob.Register(EventSingleStepBarrierReached{})
	gob.Register(HeloRequest{})
	gob.Register(HeloResponse{})
	gob.Register(StartSingleStepRequest{})
	gob.Register(StartSingleStepResponse{})
	gob.Register(NextSingleStepRequest{})
	gob.Register(NextSingleStepResponse{})
	gob.Register(ContinueSingleStepRequest{})
	gob.Register(ContinueSingleStepResponse{})
}
