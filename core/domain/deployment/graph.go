package deployment

import (
	"fmt"
	"github.com/vite-cloud/vite/core/domain/config"
	"sort"
)

// Node contains information about a service and its position in the dependency graph.
type Node struct {
	// Parent is a node that depends on this node.
	Parent *Node
	// Service is the service that this node represents.
	Service *config.Service
	// Edges is a list of nodes that this node depends on.
	Edges []*Node
	// Depth is the number of edges from the root node to this node.
	// The root node has a depth of 0.
	Depth int
}

func (n *Node) AddEdge(e *Node) {
	n.Edges = append(n.Edges, e)
}

func (n *Node) Walk(f func(n *Node)) {
	// We visit the current node only if it isn't the root node.
	// The root node is a stub that has no service.
	if n.Service != nil {
		f(n)
	}

	for _, e := range n.Edges {
		e.Walk(f)
	}
}

type ServiceMap map[string]*config.Service

func (s ServiceMap) Layered() ([][]*config.Service, error) {
	root := &Node{}
	unresolved := map[string]bool{}

	for _, service := range s {
		edge, err := s.graph(root, service, unresolved)
		if err != nil {
			return nil, err
		}

		root.AddEdge(edge)
	}

	nodeDepth := map[*config.Service]int{}
	depthNode := map[int][]*config.Service{}

	root.Walk(func(n *Node) {
		if nodeDepth[n.Service] < n.Depth {
			nodeDepth[n.Service] = n.Depth
		}
	})

	for node, depth := range nodeDepth {
		depthNode[depth] = append(depthNode[depth], node)
	}

	reversed := make([][]*config.Service, len(depthNode))

	for key, nodes := range depthNode {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})

		reversed[len(depthNode)-key] = nodes
	}

	return reversed, nil
}

func (s ServiceMap) graph(parent *Node, service *config.Service, unresolved map[string]bool) (*Node, error) {
	node := &Node{
		Parent:  parent,
		Service: service,
		Depth:   parent.Depth + 1,
	}

	unresolved[service.Name] = true

	for _, require := range service.Requires {
		if unresolved[require] {
			return nil, fmt.Errorf("circular dependency detected: %s -> %s", service.Name, require)
		}

		edge, err := s.graph(node, s[require], unresolved)
		if err != nil {
			return nil, err
		}

		node.AddEdge(edge)
	}

	delete(unresolved, service.Name)

	return node, nil
}
