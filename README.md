# Go-SWD (Sensitive Words Detection)

一个高性能的敏感词检测和过滤库，基于 Go 语言开发，采用整洁架构设计。专注于中文文本的敏感词检测，支持多种检测策略和灵活的扩展机制。

## 主要特性

- 🚀 高性能：支持 Trie 和 AC 自动机等算法
- 🎯 精准检测：支持多种文本匹配策略
- 📚 内置词库：提供常用的内置敏感词库
- 🔄 灵活分类：支持多种敏感词分类（涉黄、涉政、暴力等），可独立开关
- 🛠 可扩展：支持自定义词库扩展
- 📦 轻量级：无外部依赖，即插即用
- 🔒 安全性：内置多种反规避机制

## 反规避特性

V1.0 版本支持：
- 基础文本匹配
- 特殊字符过滤

后续版本规划：
- 大小写混淆检测（如：FuCk -> fuck）
- 全半角混淆检测（如：ｆｕｃｋ -> fuck）
- 数字样式检测（如：9⓿二肆⁹₈ -> 902498）
- 特殊字符插入检测（如：f*u*c*k -> fuck）
- 中文拼音混合检测
- 同音字检测
- 形近字检测

## 快速开始

### 安装

```bash
go get github.com/kirklin/go-swd
```

### 基础使用

```go
package main

import (
	"fmt"
	"log"

	"github.com/kirklin/go-swd"
)

func main() {
	// 1. 创建实例
	detector, err := swd.New()
	if err != nil {
		log.Fatal(err)
	}

	// 2. 添加自定义敏感词（可选）
	customWords := map[string]swd.Category{
		"涉黄":    swd.Pornography,    // 涉黄分类
		"涉政":    swd.Political,      // 涉政分类
		"赌博词汇":  swd.Gambling,       // 赌博分类
		"毒品词汇":  swd.Drugs,          // 毒品分类
		"脏话词汇":  swd.Profanity,      // 脏话分类
		"歧视词汇":  swd.Discrimination, // 歧视分类
		"诈骗词汇":  swd.Scam,           // 诈骗分类
		"自定义词汇": swd.Custom,         // 自定义分类
	}
	if err := detector.AddWords(customWords); err != nil {
		log.Fatal(err)
	}

	// 3. 基本检测
	text := "这是一段包含敏感词涉黄和涉政的文本"
	fmt.Println("是否包含敏感词:", detector.Detect(text))

	// 4. 检测指定分类
	fmt.Println("是否包含涉黄内容:", detector.DetectIn(text, swd.Pornography))
	fmt.Println("是否包含涉政内容:", detector.DetectIn(text, swd.Political))
	fmt.Println("是否包含赌博内容:", detector.DetectIn(text, swd.Gambling))
	fmt.Println("是否包含毒品内容:", detector.DetectIn(text, swd.Drugs))

	// 5. 检测多个分类
	fmt.Println("是否包含涉黄或涉政内容:", detector.DetectIn(text, swd.Pornography, swd.Political))
	fmt.Println("是否包含任意预定义分类:", detector.DetectIn(text, swd.All))

	// 6. 获取匹配结果
	if word := detector.Match(text); word != nil {
		fmt.Printf("首个敏感词: %s (分类: %s)\n", word.Word, word.Category)
	}

	// 7. 获取所有匹配
	words := detector.MatchAll(text)
	for _, word := range words {
		fmt.Printf("敏感词: %s (分类: %s, 位置: %d-%d)\n",
			word.Word, word.Category, word.StartPos, word.EndPos)
	}

	// 8. 敏感词过滤
	filtered := detector.ReplaceWithAsterisk(text) // 使用 * 替换
	fmt.Println("过滤后的文本:", filtered)

	// 9. 自定义替换策略
	customFiltered := detector.ReplaceWithStrategy(text, func(word swd.SensitiveWord) string {
		return fmt.Sprintf("[%s]", word.Category) // 替换为分类名
	})
	fmt.Println("自定义替换后的文本:", customFiltered)

	// 10. 移除敏感词
	if err := detector.RemoveWord("自定义敏感词1"); err != nil {
		log.Printf("移除敏感词失败: %v", err)
	}

	// 11. 清空词库
	if err := detector.Clear(); err != nil {
		log.Printf("清空词库失败: %v", err)
	}
}
```

## 项目结构

```
pkg/
├── core/           # 核心接口定义
├── types/          # 基础类型定义
├── detector/       # 敏感词检测算法实现
├── filter/         # 敏感词过滤策略
├── dictionary/     # 词库管理
│   ├── default/    # 内置词库
│   │   ├── pornography.txt  # 涉黄词库
│   │   ├── political.txt    # 涉政词库
│   │   └── violence.txt     # 暴力词库
│   └── loader.go   # 词库加载器
├── normalize/      # 文本标准化处理
└── swd/         # 整合各个模块，提供统一的对外接口

```

## 版本规划

### V1.0
- 基础的 Trie 树实现
- 内置词库支持
- 基本的敏感词检测和过滤
- 支持自定义词库扩展

### V1.1
- 添加全半角转换
- 添加大小写转换
- 添加特殊字符过滤

### V1.2
- 添加 AC 自动机算法
- 优化性能

### V2.0
- 拼音检测
- 同音字检测
- 形近字检测

## 许可证

本项目采用 Apache 许可证。详情请见 [LICENSE](LICENSE) 文件。
