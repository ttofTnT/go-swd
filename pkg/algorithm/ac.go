package algorithm

import (
	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/types/category"
)

// AhoCorasickNode Aho-Corasick算法节点
type AhoCorasickNode struct {
	children map[rune]*AhoCorasickNode // 子节点映射
	failLink *AhoCorasickNode          // 失败指针
	isEnd    bool                      // 是否是单词结尾
	word     string                    // 如果是结尾节点，存储完整词
	category category.Category         // 敏感词分类
	parent   *AhoCorasickNode          // 父节点 (用于重建单词)
	depth    int                       // 在字典树中的深度
}

// newAhoCorasickNode 创建新的Aho-Corasick算法节点
func newAhoCorasickNode() *AhoCorasickNode {
	return &AhoCorasickNode{
		children: make(map[rune]*AhoCorasickNode),
	}
}

// AhoCorasick Aho-Corasick算法实现
type AhoCorasick struct {
	root  *AhoCorasickNode
	built bool // 是否已构建失败指针
}

// NewAhoCorasick 创建新的Aho-Corasick算法实例
func NewAhoCorasick() *AhoCorasick {
	return &AhoCorasick{
		root: newAhoCorasickNode(),
	}
}

// Type 返回算法类型
func (ac *AhoCorasick) Type() core.AlgorithmType {
	return core.AlgorithmAhoCorasick
}

// Build 构建Aho-Corasick算法词库
func (ac *AhoCorasick) Build(words map[string]category.Category) error {
	ac.root = newAhoCorasickNode()
	for word, category := range words {
		ac.insert(word, category)
	}

	ac.buildFailureLinks()
	return nil
}

// insert 向自动机中添加一个词
func (ac *AhoCorasick) insert(word string, category category.Category) {
	if word == "" {
		return
	}

	current := ac.root
	for i, char := range word {
		if _, exists := current.children[char]; !exists {
			current.children[char] = newAhoCorasickNode()
			current.children[char].parent = current
			current.children[char].depth = i + 1
		}
		current = current.children[char]
	}

	current.isEnd = true
	current.word = word
	current.category = category
	ac.built = false // 需要重新构建失败指针
}

// buildFailureLinks 构建失败指针
func (ac *AhoCorasick) buildFailureLinks() {
	if ac.built {
		return
	}

	// 使用BFS构建失败指针
	queue := make([]*AhoCorasickNode, 0)

	// 先处理根节点的子节点
	for _, child := range ac.root.children {
		child.failLink = ac.root
		queue = append(queue, child)
	}

	// 处理剩余节点
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			queue = append(queue, child)

			// 寻找失败指针
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

// Match 查找文本中的第一个匹配
func (ac *AhoCorasick) Match(text string) *core.SensitiveWord {
	if !ac.built {
		ac.buildFailureLinks()
	}

	current := ac.root
	runes := []rune(text)

	for pos, char := range runes {
		// 查找下一个状态
		for current != ac.root && current.children[char] == nil {
			current = current.failLink
		}

		if next, exists := current.children[char]; exists {
			current = next
		} else {
			continue
		}

		// 检查当前节点的匹配
		for node := current; node != ac.root; node = node.failLink {
			if node.isEnd {
				wordRunes := []rune(node.word)
				startPos := pos - len(wordRunes) + 1
				return &core.SensitiveWord{
					Word:     node.word,
					StartPos: startPos,
					EndPos:   pos + 1,
					Category: node.category,
				}
			}
		}
	}

	return nil
}

// MatchAll 返回文本中所有敏感词
func (ac *AhoCorasick) MatchAll(text string) []core.SensitiveWord {
	if !ac.built {
		ac.buildFailureLinks()
	}

	var matches []core.SensitiveWord
	current := ac.root
	runes := []rune(text)

	for pos, char := range runes {
		// 查找下一个状态
		for current != ac.root && current.children[char] == nil {
			current = current.failLink
		}

		if next, exists := current.children[char]; exists {
			current = next
		} else {
			continue
		}

		// 检查当前节点的所有匹配
		for node := current; node != ac.root; node = node.failLink {
			if node.isEnd {
				wordRunes := []rune(node.word)
				startPos := pos - len(wordRunes) + 1
				match := core.SensitiveWord{
					Word:     node.word,
					StartPos: startPos,
					EndPos:   pos + 1,
					Category: node.category,
				}
				matches = append(matches, match)
			}
		}
	}

	return matches
}

// Replace 替换敏感词
func (ac *AhoCorasick) Replace(text string, replacement rune) string {
	matches := ac.MatchAll(text)
	if len(matches) == 0 {
		return text
	}

	runes := []rune(text)
	for _, match := range matches {
		for i := match.StartPos; i < match.EndPos; i++ {
			runes[i] = replacement
		}
	}
	return string(runes)
}

// Detect 检查文本是否包含敏感词
func (ac *AhoCorasick) Detect(text string) bool {
	return ac.Match(text) != nil
}
