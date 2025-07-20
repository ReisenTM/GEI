package gei

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node //子节点
	isWild   bool    //是否精确匹配 part 含有 `:` 或 `*` 时为true
}

// 找第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	//如果没找到
	return nil
}

// 找所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//如果是最后一层节点，匹配结束
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//如果没有子节点，就插入
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	//递归插入
	child.insert(pattern, parts, height+1)
}

// search 匹配路由
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height {
		//查找结束
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		//逐层查找看是否存在绑定的路由
		result := child.search(parts, height+1)
		if result != nil {
			//如果找到就向上层返回
			return result
		}
	}
	return nil
}
