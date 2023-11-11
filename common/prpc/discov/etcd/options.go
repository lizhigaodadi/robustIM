package etcd

import "time"

var (
	defaultDialTimeOut          = 10 * time.Second
	defaultDialKeepAliveTimeOut = 10 * time.Second
	defaultAutoSyncInterval     = 10 * time.Second
	defaultLeaseTTL             = int64(10)
)

type Options struct {
	DialTimeTimeOut      time.Duration
	DialKeepAliveTimeOut time.Duration
	AutoSyncInterval     time.Duration
	LeaseTTL             int64
}

type Option func(o *Options)

func NewOptions(os ...Option) *Options {
	opt := &Options{
		AutoSyncInterval:     defaultAutoSyncInterval,
		DialKeepAliveTimeOut: defaultDialKeepAliveTimeOut,
		DialTimeTimeOut:      defaultDialTimeOut,
		LeaseTTL:             defaultLeaseTTL,
	}

	for _, f := range os {
		f(opt)
	}
	return opt
}

func WithDialTimeOut(dialTimeTimeOut time.Duration) Option {
	return func(o *Options) {
		o.DialTimeTimeOut = dialTimeTimeOut
	}
}

func WithDialKeepAliveTimeOut(dialKeepAliveTimeOut time.Duration) Option {
	return func(o *Options) {
		o.DialKeepAliveTimeOut = dialKeepAliveTimeOut
	}
}

func WithAutoSyncInterval(autoSyncInterval time.Duration) Option {
	return func(o *Options) {
		o.AutoSyncInterval = autoSyncInterval
	}
}
