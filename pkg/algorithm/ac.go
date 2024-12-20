package algorithm

import (
	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/types"
)

// ACNode AC自动机节点
// AC Automaton Node
type ACNode struct {
	children map[rune]*ACNode // 子节点映射 / Child nodes mapping
	failLink *ACNode          // 失败指针 / Failure link
	isEnd    bool             // 是否是单词结尾 / Whether it's the end of a word
	word     string           // 如果是结尾节点，存储完整词 / Store complete word if it's an end node
	category types.Category   // 敏感词分类 / Sensitive word category
	parent   *ACNode          // 父节点 (用于重建单词) / Parent node (for word reconstruction)
	depth    int              // 在字典树中的深度 / Depth in the trie
}

// newACNode 创建新的AC自动机节点
// Create a new AC automaton node
func newACNode() *ACNode {
	return &ACNode{
		children: make(map[rune]*ACNode),
	}
}

// ACAutomaton AC自动机
// AC Automaton
type ACAutomaton struct {
	root  *ACNode
	built bool // 是否已构建失败指针 / Whether failure links have been built
}

// NewACAutomaton 创建新的AC自动机
// Create a new AC automaton
func NewACAutomaton() *ACAutomaton {
	return &ACAutomaton{
		root: newACNode(),
	}
}

// Insert 向自动机中添加一个词
// Add a word to the automaton
func (ac *ACAutomaton) Insert(word string, category types.Category) {
	if word == "" {
		return
	}

	current := ac.root
	for i, char := range word {
		if _, exists := current.children[char]; !exists {
			current.children[char] = newACNode()
			current.children[char].parent = current
			current.children[char].depth = i + 1
		}
		current = current.children[char]
	}

	current.isEnd = true
	current.word = word
	current.category = category
	ac.built = false // 需要重新构建失败指针 / Need to rebuild failure links
}

// buildFailureLinks 构建失败指针
// Build failure links
func (ac *ACAutomaton) buildFailureLinks() {
	if ac.built {
		return
	}

	// 使用BFS构建失败指针 / Use BFS to build failure links
	queue := make([]*ACNode, 0)

	// 先处理根节点的子节点 / Handle root's children first
	for _, child := range ac.root.children {
		child.failLink = ac.root
		queue = append(queue, child)
	}

	// 处理剩余节点 / Process remaining nodes
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			queue = append(queue, child)

			// 寻找失败指针 / Find failure link
			failNode := current.failLink
			for failNode != nil {
				if next, exists := failNode.children[char]; exists {
					child.failLink = next
					break
				}
				failNode = failNode.failLink
			}
			if failNode == nil {
				child.failLink = ac.root
			}
		}
	}

	ac.built = true
}

// MatchFirst 查找文本中的第一个匹配
// Find the first match in the text
func (ac *ACAutomaton) MatchFirst(text string) *core.SensitiveWord {
	if !ac.built {
		ac.buildFailureLinks()
	}

	current := ac.root

	for pos, char := range text {
		// 查找下一个状态 / Find next state
		for current != ac.root && current.children[char] == nil {
			current = current.failLink
		}

		if next, exists := current.children[char]; exists {
			current = next
		} else {
			continue
		}

		// 检查当前节点的匹配 / Check matches at current node
		for node := current; node != ac.root; node = node.failLink {
			if node.isEnd {
				return &core.SensitiveWord{
					Word:     node.word,
					StartPos: pos - node.depth + 1,
					EndPos:   pos + 1,
					Category: node.category,
				}
			}
		}
	}

	return nil
}

// MatchAll 查找文本中的所有匹配
// Find all matches in the text
func (ac *ACAutomaton) MatchAll(text string) []core.SensitiveWord {
	if !ac.built {
		ac.buildFailureLinks()
	}

	var matches []core.SensitiveWord
	current := ac.root

	for pos, char := range text {
		// 查找下一个状态 / Find next state
		for current != ac.root && current.children[char] == nil {
			current = current.failLink
		}

		if next, exists := current.children[char]; exists {
			current = next
		} else {
			continue
		}

		// 检查当前节点的所有匹配 / Check all matches at current node
		for node := current; node != ac.root; node = node.failLink {
			if node.isEnd {
				match := core.SensitiveWord{
					Word:     node.word,
					StartPos: pos - node.depth + 1,
					EndPos:   pos + 1,
					Category: node.category,
				}
				matches = append(matches, match)
			}
		}
	}

	return matches
}
