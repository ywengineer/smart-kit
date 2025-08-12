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

func WithServiceName(serviceName string) Option {
	return func(p *Provider) {
		p.serviceName = serviceName
	}
}

func WithGroupName(groupName string) Option {
	return func(p *Provider) {
		p.groupName = groupName
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

func WithEphemeral() Option {
	return func(p *Provider) {
		p.ephemeral = true
	}
}
