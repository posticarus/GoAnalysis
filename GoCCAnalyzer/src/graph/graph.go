// The MIT License (MIT)

// Copyright (c) 2015-2016 Christian Bergum Bergersen

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package graph

import (
	"github.com/chrisbbe/GoAnalysis/Stack/src"
)

type Graph struct {
	Root     *Node
	Nodes    map[interface{}]*Node
	scc      []*StronglyConnectedComponent
	preCount int         // Used by SCC.
	stack    stack.Stack // Used by SCC.
}

type Node struct {
	Value    interface{}
	visited  bool
	inEdges  []*Node
	outEdges []*Node
	low      int // Used by SCC.
}

type StronglyConnectedComponent struct {
	Nodes []*Node
}

func NewGraph() *Graph {
	return &Graph{Nodes: map[interface{}]*Node{}}
}

func (n *Node) insertOutNode(node *Node) {
	n.outEdges = append(n.outEdges, node)
}

func (n *Node) insertInNode(node *Node) {
	n.inEdges = append(n.inEdges, node)
}

func (graph *Graph) InsertEdge(leftNode *Node, rightNode *Node) {
	if len(graph.Nodes) == 0 {
		graph.Root = leftNode
		graph.Root.insertOutNode(rightNode)
		rightNode.insertInNode(graph.Root)
		graph.Nodes[graph.Root.Value] = graph.Root
		graph.Nodes[rightNode.Value] = rightNode
	} else {
		//Get left and right node if they already exist.
		if graph.Nodes[leftNode.Value] != nil {
			leftNode = graph.Nodes[leftNode.Value]
		}
		if graph.Nodes[rightNode.Value] != nil {
			rightNode = graph.Nodes[rightNode.Value]
		}

		if graph.Nodes[leftNode.Value] == nil {
			leftNode.insertOutNode(rightNode)
			rightNode.insertInNode(leftNode)
			graph.Nodes[leftNode.Value] = leftNode
			graph.Nodes[rightNode.Value] = rightNode
		} else {
			leftNode.insertOutNode(rightNode)
			rightNode.insertInNode(leftNode)
			graph.Nodes[rightNode.Value] = rightNode
		}
	}
}

func (node *Node) getDFS() (nodes []*Node) {
	if !node.visited {
		nodes = append(nodes, node)
	}
	node.visited = true
	for _, n := range node.outEdges {
		if !n.visited {
			nodes = append(nodes, n.getDFS()...)
		}
	}
	return nodes
}

func (graph *Graph) GetDFS() (nodes []*Node) {
	nodes = graph.Root.getDFS()
	//Clean up nodes by setting visited = false
	for _, node := range graph.Nodes {
		node.visited = false
	}
	return nodes
}

func (graph *Graph) dfs(v *Node) {
	v.low = graph.preCount
	graph.preCount++
	v.visited = true
	graph.stack.Push(v)

	min := v.low
	for _, w := range v.outEdges {
		if !w.visited {
			graph.dfs(w)
		}
		if w.low < min {
			min = w.low
		}
	}

	if min < v.low {
		v.low = min
		return
	}

	component := []*Node{}
	var w interface{}
	for ok := true; ok; ok = w.(*Node) != v {
		w, _ = graph.stack.Pop()
		component = append(component, w.(*Node))
		w.(*Node).low = len(graph.Nodes) - 1
	}

	graph.scc = append(graph.scc, &StronglyConnectedComponent{Nodes:component})
}

func (graph *Graph) GetSCComponents() ([]*StronglyConnectedComponent) {
	graph.scc = []*StronglyConnectedComponent{} //Init list.
	for _, node := range graph.Nodes {
		if !node.visited {
			graph.dfs(node)
		}
	}
	//Clean up nodes by setting visited = false
	for _, node := range graph.Nodes {
		node.visited = false
		node.low = 0
	}
	//Clean up graph.
	graph.stack = stack.Stack{}
	graph.preCount = 0
	return graph.scc
}

func (graph *Graph) GetNumberOfNodes() (int) {
	return len(graph.Nodes)
}

func (graph *Graph) GetNumberOfEdges() (numberOfEdges int) {
	for _, node := range graph.Nodes {
		numberOfEdges += len(node.outEdges)
	}
	return numberOfEdges
}

func (graph *Graph) GetNumberOfSCComponents() (int) {
	//We don't want to run Tarjan's once again if algorithm is already executed.
	if graph.scc == nil {
		return len(graph.GetSCComponents())
	}
	return len(graph.scc)
}
