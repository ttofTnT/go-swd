package swd

import (
	"fmt"

	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/detector"
	"github.com/kirklin/go-swd/pkg/dictionary"
	"github.com/kirklin/go-swd/pkg/filter"
)

// DefaultFactory 默认组件工厂实现
type DefaultFactory struct{}

// NewDefaultFactory 创建默认工厂实例
func NewDefaultFactory() ComponentFactory {
	return &DefaultFactory{}
}

// CreateDetector 创建检测器实例
func (f *DefaultFactory) CreateDetector(options *core.SWDOptions) core.Detector {
	detector, err := detector.NewDetector(*options)
	if err != nil {
		panic(fmt.Sprintf("创建检测器失败: %v", err))
	}
	return detector
}

// CreateFilter 创建过滤器实例
func (f *DefaultFactory) CreateFilter(detector core.Detector) core.Filter {
	return filter.NewFilter(detector)
}

// CreateLoader 创建加载器实例
func (f *DefaultFactory) CreateLoader() core.Loader {
	return dictionary.NewLoader()
}
