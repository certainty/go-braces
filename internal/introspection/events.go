package introspection

type StartCompileModuleEvent struct{}

func (e StartCompileModuleEvent) EventInspect() string {
	return "StartCompileModule"
}

func EventStartCompileModule() StartCompileModuleEvent {
	return StartCompileModuleEvent{}
}

type EndCompileModuleEvent struct{}

func (e EndCompileModuleEvent) EventInspect() string {
	return "EndCompileModule"
}

func EventEndCompileModule() EndCompileModuleEvent {
	return EndCompileModuleEvent{}
}
