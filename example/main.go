package main

import (
	"fmt"
	"github.com/kirklin/go-swd/pkg/core"
	"log"

	"github.com/kirklin/go-swd/pkg/swd"
	"github.com/kirklin/go-swd/pkg/types/category"
)

func main() {
	// 1. 创建实例
	factory := swd.NewDefaultFactory()
	detector, err := swd.New(factory)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 添加自定义敏感词（可选）
	customWords := map[string]category.Category{
		"自定义敏感词1": category.Pornography,
		"自定义敏感词2": category.Political,
		"seqing":  category.Pornography,
	}
	if err := detector.AddWords(customWords); err != nil {
		log.Fatal(err)
	}

	// 3. 基本检测
	text := "这是一段包含色情词的文本,自定义敏感词1"
	fmt.Println("是否包含敏感词:", detector.Detect(text))

	// 4. 检测指定分类
	fmt.Println("是否包含涉黄内容:", detector.DetectIn(text, category.Pornography))
	fmt.Println("是否包含涉政内容:", detector.DetectIn(text, category.Political))

	// 5. 获取匹配结果
	if word := detector.Match(text); word != nil {
		fmt.Printf("首个敏感词: %s (分类: %s)\n", word.Word, word.Category)
	}

	// 6. 获取所有匹配
	words := detector.MatchAll(text)
	for i, word := range words {
		fmt.Printf("敏感词%d: %s (分类: %s, 位置: %d-%d)\n",
			i+1, word.Word, word.Category, word.StartPos, word.EndPos)
	}

	// 7. 敏感词过滤
	fmt.Println("过滤前的文本:", text)
	filtered := detector.ReplaceWithAsterisk(text) // 使用 * 替换
	fmt.Println("过滤后的文本:", filtered)

	// 8. 自定义替换策略
	customFiltered := detector.ReplaceWithStrategy(text, func(word core.SensitiveWord) string {
		return fmt.Sprintf("[%s]", word.Category) // 替换为分类名
	})
	fmt.Println("自定义替换后的文本:", customFiltered)

	// 9. 移除敏感词
	if err := detector.RemoveWord("自定义敏感词1"); err != nil {
		log.Printf("移除敏感词失败: %v", err)
	}

	// 10. 清空词库
	if err := detector.Clear(); err != nil {
		log.Printf("清空词库失败: %v", err)
	}
}
