package client

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

// NewClient returns a new containerd client that communicates over
// the unix socket at the given path.
//
// Interesting note: there is no remote client for containerd.
// The debate is around security.
func NewClient(socketPath string) (*containerd.Client, error) {
	return containerd.New(socketPath)
}

// WithNamespace returns a new context with the given namespace.
// Context is used by containerd to pass around optional information.
//
// The namespace is used to isolate resources.
//
// You can use "well-known" namespaces like "docker" or "k8s.io" expose
// containers to other tools.
func WithNamespace(ctx context.Context, namespace string) context.Context {
	return namespaces.WithNamespace(ctx, namespace)
}
