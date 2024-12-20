package dictionary

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/kirklin/go-swd/pkg/types/category"
)

// Loader 实现core.Loader接口
type Loader struct {
	words map[string]category.Category
	mu    sync.RWMutex
}

// NewLoader 创建新的加载器实例
func NewLoader() *Loader {
	return &Loader{
		words: make(map[string]category.Category),
	}
}

//go:embed default/political.txt
var politicalWords string

//go:embed default/pornography.txt
var pornographyWords string

//go:embed default/violence.txt
var violenceWords string

//go:embed default/all.txt
var allWords string

// LoadDefaultWords 加载默认词库
func (l *Loader) LoadDefaultWords(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 加载所有分类词典
	categories := map[string]struct {
		content string
		cat     category.Category
	}{
		"political.txt":   {content: politicalWords, cat: category.Political},
		"pornography.txt": {content: pornographyWords, cat: category.Pornography},
		"violence.txt":    {content: violenceWords, cat: category.Violence},
	}

	for filename, data := range categories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := l.loadFromString(ctx, data.content, data.cat); err != nil {
				return fmt.Errorf("failed to load %s: %w", filename, err)
			}
		}
	}

	// 加载通用词典
	if err := l.loadFromString(ctx, allWords, category.None); err != nil {
		return fmt.Errorf("failed to load all.txt: %w", err)
	}

	return nil
}

// LoadCustomWords 加载自定义词库
func (l *Loader) LoadCustomWords(ctx context.Context, words map[string]category.Category) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	const batchSize = 1000
	count := 0

	for word, cat := range words {
		if count%batchSize == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
		if err := l.addWordLocked(word, cat); err != nil {
			return err
		}
		count++
	}
	return nil
}

// AddWord 添加单个敏感词
func (l *Loader) AddWord(word string, cat category.Category) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.addWordLocked(word, cat)
}

// addWordLocked 在已获得锁的情况下添加单个敏感词
func (l *Loader) addWordLocked(word string, cat category.Category) error {
	if word = strings.TrimSpace(word); word == "" {
		return fmt.Errorf("word cannot be empty")
	}

	// 如果词已存在且有效分类，且当前要设置的是 None 分类，则保留原有分类
	if existingCat, exists := l.words[word]; exists && existingCat != category.None && cat == category.None {
		return nil
	}

	l.words[word] = cat
	return nil
}

// AddWords 批量添加敏感词
func (l *Loader) AddWords(words map[string]category.Category) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for word, cat := range words {
		if err := l.addWordLocked(word, cat); err != nil {
			return err
		}
	}
	return nil
}

// RemoveWord 移除单个敏感词
func (l *Loader) RemoveWord(word string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.words, word)
	return nil
}

// RemoveWords 批量移除敏感词
func (l *Loader) RemoveWords(words []string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, word := range words {
		delete(l.words, word)
	}
	return nil
}

// Clear 清空所有敏感词
func (l *Loader) Clear() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.words = make(map[string]category.Category)
	return nil
}

// loadFromString 从字符串加载敏感词
func (l *Loader) loadFromString(ctx context.Context, content string, cat category.Category) error {
	reader := strings.NewReader(content)
	return l.loadFromReader(ctx, reader, cat)
}

// loadFromReader 从Reader加载敏感词
func (l *Loader) loadFromReader(ctx context.Context, reader io.Reader, cat category.Category) error {
	scanner := bufio.NewScanner(reader)
	const batchSize = 1000
	count := 0

	for scanner.Scan() {
		if count%batchSize == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		word := strings.TrimSpace(scanner.Text())
		if word == "" || strings.HasPrefix(word, "#") {
			continue
		}
		if err := l.addWordLocked(word, cat); err != nil {
			return err
		}
		count++
	}
	return scanner.Err()
}

// GetWords 获取所有已加载的敏感词
func (l *Loader) GetWords() map[string]category.Category {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 创建一个副本以避免并发访问问题
	words := make(map[string]category.Category, len(l.words))
	for k, v := range l.words {
		words[k] = v
	}
	return words
}
