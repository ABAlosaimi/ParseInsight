package adapters

import "fmt"

// Registry holds all available parser adapters
type Registry struct {
	adapters map[string]func(messageType string) ParserAdapter
}

// NewRegistry creates a new adapter registry
func NewRegistry() *Registry {
	r := &Registry{
		adapters: make(map[string]func(messageType string) ParserAdapter),
	}

	// Register available adapters
	r.Register("net/http", func(mt string) ParserAdapter {
		return NewNetHTTPAdapter(mt)
	})
	r.Register("fasthttp", func(mt string) ParserAdapter {
		return NewFastHTTPAdapter(mt)
	})

	return r
}

// Register adds an adapter factory to the registry
func (r *Registry) Register(name string, factory func(messageType string) ParserAdapter) {
	r.adapters[name] = factory
}

// Get retrieves an adapter by name
func (r *Registry) Get(name, messageType string) (ParserAdapter, error) {
	factory, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("adapter '%s' not found", name)
	}
	return factory(messageType), nil
}

// Available returns list of available adapter names
func (r *Registry) Available() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}
