package launcher

type Option func(l *Launcher)

func WithServer(server Server) Option {
	return func(l *Launcher) {
		l.servers = append(l.servers, server)
	}
}
