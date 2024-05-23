package ff

import "strings"

type node struct {
	pattern  string  // 待匹配的路由 e.g /p/:lang
	part     string  // 路由中的一部分 e.g :lang
	children []*node // 子节点
	isWild   bool    // 是否精确匹配，part含有 : 或 *时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		newNode := &node{
			part:   part,
			isWild: strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"),
		}

		n.children = append(n.children, newNode)
		newNode.insert(pattern, parts, height+1)
	} else {
		child.insert(pattern, parts, height+1)
	}
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	for _, child := range n.matchChildren(part) {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
