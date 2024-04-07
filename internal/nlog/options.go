package nlog

type optionList struct {
	rootKVAdder func() []KV
	kvConverter func(KV) []KV
}

func defaultOptionList() *optionList {
	return &optionList{
		kvConverter: func(kv KV) []KV {
			switch kv.Kind {
			case KindContext:
				return nil
			case KindError:
				if kv.Value == nil {
					return []KV{}
				}

				return []KV{KVVAny(kv.Key, kv.Value)}
			default:
				return []KV{KVVAny(kv.Key, kv.Value)}
			}
		},
	}
}

type Option func(*optionList)

func WithKVConverter(kvConverter func(KV) []KV) Option {
	return func(o *optionList) {
		o.kvConverter = kvConverter
	}
}

func WithROotKVAdder(rootKVAdder func() []KV) Option {
	return func(o *optionList) {
		o.rootKVAdder = rootKVAdder
	}
}
