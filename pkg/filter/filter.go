package filter

import (
	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/types/category"
)

// filter 实现了敏感词过滤器接口
type filter struct {
	detector core.Detector
}

// NewFilter 创建一个新的过滤器实例
func NewFilter(detector core.Detector) core.Filter {
	return &filter{
		detector: detector,
	}
}

// Replace 使用指定的替换字符替换敏感词
func (f *filter) Replace(text string, replacement rune) string {
	if text == "" {
		return text
	}
	return f.ReplaceWithStrategy(text, func(word core.SensitiveWord) string {
		chars := make([]rune, len([]rune(word.Word)))
		for i := range chars {
			chars[i] = replacement
		}
		return string(chars)
	})
}

// ReplaceIn 使用指定的替换字符替换指定分类的敏感词
func (f *filter) ReplaceIn(text string, replacement rune, categories ...category.Category) string {
	if text == "" || len(categories) == 0 {
		return text
	}
	matches := f.detector.MatchAllIn(text, categories...)
	return f.replaceWords(text, matches, func(word core.SensitiveWord) string {
		chars := make([]rune, len([]rune(word.Word)))
		for i := range chars {
			chars[i] = replacement
		}
		return string(chars)
	})
}

// ReplaceWithAsterisk 使用 * 号替换敏感词
func (f *filter) ReplaceWithAsterisk(text string) string {
	return f.Replace(text, '*')
}

// ReplaceWithAsteriskIn 使用 * 号替换指定分类的敏感词
func (f *filter) ReplaceWithAsteriskIn(text string, categories ...category.Category) string {
	return f.ReplaceIn(text, '*', categories...)
}

// ReplaceWithStrategy 使用自定义替换策略替换敏感词
func (f *filter) ReplaceWithStrategy(text string, strategy func(word core.SensitiveWord) string) string {
	if text == "" || strategy == nil {
		return text
	}
	matches := f.detector.MatchAll(text)
	return f.replaceWords(text, matches, strategy)
}

// ReplaceWithStrategyIn 使用自定义替换策略替换指定分类的敏感词
func (f *filter) ReplaceWithStrategyIn(text string, strategy func(word core.SensitiveWord) string, categories ...category.Category) string {
	if text == "" || strategy == nil || len(categories) == 0 {
		return text
	}
	matches := f.detector.MatchAllIn(text, categories...)
	return f.replaceWords(text, matches, strategy)
}

// replaceWords 替换文本中的敏感词
func (f *filter) replaceWords(text string, matches []core.SensitiveWord, strategy func(word core.SensitiveWord) string) string {
	if len(matches) == 0 {
		return text
	}

	runes := []rune(text)
	result := make([]rune, 0, len(runes))
	lastPos := 0

	for _, match := range matches {
		// 添加敏感词前的文本
		result = append(result, runes[lastPos:match.StartPos]...)
		// 添加替换后的文本
		replacement := []rune(strategy(match))
		result = append(result, replacement...)
		lastPos = match.EndPos
	}

	// 添加最后一个敏感词后的文本
	if lastPos < len(runes) {
		result = append(result, runes[lastPos:]...)
	}

	return string(result)
}
