package introspection_events

import "encoding/gob"

func RegisterEventTypes() {
	gob.Register(EventBeginCompileString{})
	gob.Register(EventEndCompileString{})
	gob.Register(EventBeginCompileModule{})
	gob.Register(EventEndCompileModule{})
}
