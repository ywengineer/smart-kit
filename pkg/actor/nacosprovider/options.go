package nacosprovider

import "time"

type Option func(p *Provider)

func WithWeight(w int) Option {
	return func(p *Provider) {
		p.weight = w
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(p *Provider) {
		p.metadata = metadata
	}
}

func WithClusterName(clusterName string) Option {
	return func(p *Provider) {
		p.clusterName = clusterName
	}
}

func WithGroupName(groupName string) Option {
	return func(p *Provider) {
		p.groupName = groupName
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(p *Provider) {
		p.ttl = ttl
	}
}

func WithRefreshTTL(refreshTTL time.Duration) Option {
	return func(p *Provider) {
		p.refreshTTL = refreshTTL
	}
}

func WithNamespace(namespace string) Option {
	return func(p *Provider) {
		p.namespace = namespace
	}
}
