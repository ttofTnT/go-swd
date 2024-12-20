package detector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kirklin/go-swd/pkg/algorithm"
	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/detector/preprocessor"
	"github.com/kirklin/go-swd/pkg/dictionary"
	"github.com/kirklin/go-swd/pkg/types/category"
)

// detector 实现敏感词检测器接口
type detector struct {
	algo       core.Algorithm
	preprocess *preprocessor.Preprocessor
	mu         sync.RWMutex
	options    core.SWDOptions
}

// NewDetector 创建一个新的检测器实例
func NewDetector(options core.SWDOptions) (core.Detector, error) {
	ahoCorasick := algorithm.NewAhoCorasick()
	loader := dictionary.NewLoader()

	// 加载词典
	if err := loader.LoadDefaultWords(context.Background()); err != nil {
		return nil, fmt.Errorf("加载默认词典失败: %w", err)
	}

	// 获取词典内容
	words := loader.GetWords()
	if len(words) == 0 {
		return nil, fmt.Errorf("词典内容为空")
	}

	// 构建算法
	if err := ahoCorasick.Build(words); err != nil {
		return nil, fmt.Errorf("构建算法失败: %w", err)
	}

	d := &detector{
		algo:       ahoCorasick,
		preprocess: preprocessor.NewPreprocessor(options),
		options:    options,
	}

	// 注册为观察者
	loader.AddObserver(d)

	return d, nil
}

// OnWordsChanged 实现Observer接口,当词库变更时重建算法
func (d *detector) OnWordsChanged(words map[string]category.Category) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 重建算法
	if err := d.algo.Build(words); err != nil {
		// 这里只能记录错误,因为是回调方法
		log.Printf("重建算法失败: %v", err)
	}
}

// Detect 检查文本是否包含任何敏感词
func (d *detector) Detect(text string) bool {
	if text == "" {
		return false
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 使用读锁进行检测
	d.mu.RLock()
	match := d.algo.Match(processedText)
	d.mu.RUnlock()

	return match != nil
}

// DetectIn 检查文本是否包含指定分类的敏感词
func (d *detector) DetectIn(text string, categories ...category.Category) bool {
	if text == "" || len(categories) == 0 {
		return false
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 创建分类映射用于快速查找
	categoryMap := make(map[category.Category]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}

	// 使用读锁进行检测
	d.mu.RLock()
	match := d.algo.Match(processedText)
	d.mu.RUnlock()

	return match != nil && categoryMap[match.Category]
}

// Match 返回文本中找到的第一个敏感词
func (d *detector) Match(text string) *core.SensitiveWord {
	if text == "" {
		return nil
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 使用读锁进行检测
	d.mu.RLock()
	match := d.algo.Match(processedText)
	d.mu.RUnlock()

	return match
}

// MatchIn 返回文本中指定分类的第一个敏感词
func (d *detector) MatchIn(text string, categories ...category.Category) *core.SensitiveWord {
	if text == "" || len(categories) == 0 {
		return nil
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 创建分类映射
	categoryMap := make(map[category.Category]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}

	// 使用读锁进行检测
	d.mu.RLock()
	match := d.algo.Match(processedText)
	d.mu.RUnlock()

	if match != nil && categoryMap[match.Category] {
		return match
	}

	return nil
}

// MatchAll 返回文本中找到的所有敏感词
func (d *detector) MatchAll(text string) []core.SensitiveWord {
	if text == "" {
		return nil
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 使用读锁进行检测
	d.mu.RLock()
	matches := d.algo.MatchAll(processedText)
	d.mu.RUnlock()

	return matches
}

// MatchAllIn 返回文本中指定分类的所有敏感词
func (d *detector) MatchAllIn(text string, categories ...category.Category) []core.SensitiveWord {
	if text == "" || len(categories) == 0 {
		return nil
	}

	// 预处理文本
	processedText := d.preprocess.Process(text)

	// 创建分类映射
	categoryMap := make(map[category.Category]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}

	// 使用读锁进行检测
	d.mu.RLock()
	matches := d.algo.MatchAll(processedText)
	d.mu.RUnlock()

	// 按分类过滤匹配项
	var filteredMatches []core.SensitiveWord
	for _, match := range matches {
		if categoryMap[match.Category] {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}
