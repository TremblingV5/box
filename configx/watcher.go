package configx

func Watch(storeKey string, options ...WatchOption) WatchFunc {
	return func(watchMap WatchMap) {
		opt := &WatchSetting{
			storeKey: storeKey,
		}

		for _, option := range options {
			option(opt)

		}

		watchMap[storeKey] = opt
	}
}

type (
	WatchMap map[string]*WatchSetting

	WatchSetting struct {
		storeKey   string
		sourceType string
		location   string
		model      any
		bindSetting
	}

	WatchFunc   func(watchMap WatchMap)
	WatchOption func(*WatchSetting)

	bindSetting struct {
		callback func(value any)
	}
	BindOption func(*bindSetting)
)

func WithStoreKey(storeKey string) WatchOption {
	return func(w *WatchSetting) {
		w.storeKey = storeKey
	}
}

func WithConfigSource(sourceType, location string) WatchOption {
	return func(w *WatchSetting) {
		w.sourceType = sourceType
		w.location = location
	}
}

func WithModel(model any, options ...BindOption) WatchOption {
	return func(w *WatchSetting) {
		w.model = model

		for _, option := range options {
			option(&w.bindSetting)
		}
	}
}

func WithBindCallback(fn func(value any)) BindOption {
	return func(b *bindSetting) {
		b.callback = fn
	}
}
