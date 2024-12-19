package swd

import (
	"context"

	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/types"
)

// ComponentFactory 定义了创建各种组件的工厂接口
type ComponentFactory interface {
	CreateDetector() core.Detector
	CreateFilter(detector core.Detector) core.Filter
	CreateLoader() core.Loader
}

// SWD 敏感词检测与过滤引擎的实现
type SWD struct {
	detector core.Detector
	filter   core.Filter
	loader   core.Loader
	options  core.SWDOptions
}

// Config SWD 的配置选项
type Config struct {
	Factory ComponentFactory
	Options core.SWDOptions
}

// New 创建新的敏感词检测引擎
func New(config Config) (*SWD, error) {
	if config.Factory == nil {
		return nil, ErrNoFactory
	}

	detector := config.Factory.CreateDetector()
	if detector == nil {
		return nil, ErrNoDetector
	}

	filter := config.Factory.CreateFilter(detector)
	if filter == nil {
		return nil, ErrNoFilter
	}

	loader := config.Factory.CreateLoader()
	if loader == nil {
		return nil, ErrNoLoader
	}

	options := config.Options
	if options == (core.SWDOptions{}) {
		options = core.DefaultSWDOptions()
	}

	return &SWD{
		detector: detector,
		filter:   filter,
		loader:   loader,
		options:  options,
	}, nil
}

// SetOptions 设置引擎选项
func (swd *SWD) SetOptions(options core.SWDOptions) error {
	swd.options = options
	return nil
}

// LoadDefaultWords 加载默认词库
func (swd *SWD) LoadDefaultWords(ctx context.Context) error {
	return swd.loader.LoadDefaultWords(ctx)
}

// LoadCustomWords 加载自定义词库
func (swd *SWD) LoadCustomWords(ctx context.Context, words map[string]types.Category) error {
	return swd.loader.LoadCustomWords(ctx, words)
}

// AddWord 添加单个敏感词
func (swd *SWD) AddWord(word string, category types.Category) error {
	return swd.loader.AddWord(word, category)
}

// AddWords 批量添加敏感词
func (swd *SWD) AddWords(words map[string]types.Category) error {
	return swd.loader.AddWords(words)
}

// RemoveWord 移除敏感词
func (swd *SWD) RemoveWord(word string) error {
	return swd.loader.RemoveWord(word)
}

// RemoveWords 批量移除敏感词
func (swd *SWD) RemoveWords(words []string) error {
	return swd.loader.RemoveWords(words)
}

// Clear 清空所有敏感词
func (swd *SWD) Clear() error {
	return swd.loader.Clear()
}

// Detect 检查文本是否包含敏感词
func (swd *SWD) Detect(text string) bool {
	return swd.detector.Detect(text)
}

// DetectIn 检查文本是否包含指定分类的敏感词
func (swd *SWD) DetectIn(text string, categories ...types.Category) bool {
	return swd.detector.DetectIn(text, categories...)
}

// Match 返回文本中第一个敏感词
func (swd *SWD) Match(text string) *core.SensitiveWord {
	return swd.detector.Match(text)
}

// MatchIn 返回文本中第一个指定分类的敏感词
func (swd *SWD) MatchIn(text string, categories ...types.Category) *core.SensitiveWord {
	return swd.detector.MatchIn(text, categories...)
}

// MatchAll 返回文本中所有敏感词
func (swd *SWD) MatchAll(text string) []core.SensitiveWord {
	return swd.detector.MatchAll(text)
}

// MatchAllIn 返回文本中所有指定分类的敏感词
func (swd *SWD) MatchAllIn(text string, categories ...types.Category) []core.SensitiveWord {
	return swd.detector.MatchAllIn(text, categories...)
}

// Replace 使用指定的替换字符替换敏感词
func (swd *SWD) Replace(text string, replacement rune) string {
	return swd.filter.Replace(text, replacement)
}

// ReplaceIn 使用指定的替换字符替换指定分类的敏感词
func (swd *SWD) ReplaceIn(text string, replacement rune, categories ...types.Category) string {
	return swd.filter.ReplaceIn(text, replacement, categories...)
}

// ReplaceWithAsterisk 使用星号替换敏感词
func (swd *SWD) ReplaceWithAsterisk(text string) string {
	return swd.filter.ReplaceWithAsterisk(text)
}

// ReplaceWithAsteriskIn 使用星号替换指定分类的敏感词
func (swd *SWD) ReplaceWithAsteriskIn(text string, categories ...types.Category) string {
	return swd.filter.ReplaceWithAsteriskIn(text, categories...)
}

// ReplaceWithStrategy 使用自定义策略替换敏感词
func (swd *SWD) ReplaceWithStrategy(text string, strategy func(word core.SensitiveWord) string) string {
	return swd.filter.ReplaceWithStrategy(text, strategy)
}

// ReplaceWithStrategyIn 使用自定义策略替换指定分类的敏感词
func (swd *SWD) ReplaceWithStrategyIn(text string, strategy func(word core.SensitiveWord) string, categories ...types.Category) string {
	return swd.filter.ReplaceWithStrategyIn(text, strategy, categories...)
}
