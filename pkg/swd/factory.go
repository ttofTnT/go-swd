package swd

import (
	"context"
	"fmt"

	"github.com/ttofTnT/go-swd/pkg/core"
	"github.com/ttofTnT/go-swd/pkg/detector"
	"github.com/ttofTnT/go-swd/pkg/dictionary"
	"github.com/ttofTnT/go-swd/pkg/filter"
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
	loader := dictionary.NewLoader()
	return loader
}

// CreateComponents 创建并关联所有组件
func (f *DefaultFactory) CreateComponents(options *core.SWDOptions) (core.Detector, core.Filter, core.Loader) {
	// 创建加载器
	loader := f.CreateLoader()

	// 加载默认词库
	if err := loader.LoadDefaultWords(context.Background()); err != nil {
		panic(fmt.Sprintf("加载默认词库失败: %v", err))
	}

	// 创建检测器
	d := f.CreateDetector(options)

	// 注册检测器为加载器的观察者
	if observer, ok := d.(core.Observer); ok {
		loader.(*dictionary.Loader).AddObserver(observer)
	}

	// 创建过滤器
	filter := f.CreateFilter(d)

	return d, filter, loader
}
