package rearer

import "context"

var (
	// used to resolve circular dependency
	// if not a lock, maybe cause data race
	globalHooks []Hook
)

type HookInfo struct {
	Context    context.Context
	PanicError any
	Caller     string
	Stack      string
	Message    string
}

type Hook func(hookInfo *HookInfo)

func AddHook(hook Hook) {
	globalHooks = append(globalHooks, hook)
}
