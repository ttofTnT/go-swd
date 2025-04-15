package category

import (
	"fmt"
	"sync"
)

// Category 敏感词分类
type Category int

var (
	mu              sync.RWMutex
	categoryMap     = make(map[string]Category) // name => Category
	categoryString  = make(map[Category]string) // Category => name
	nextDynamicVal  Category                    // 动态分类起始值
	predefinedFlags Category                    // 所有静态分类合集
)

// 静态分类常量
const (
	None           Category = 0
	Pornography    Category = 1 << iota // 涉黄
	Political                           // 涉政
	Violence                            // 暴力
	Gambling                            // 赌博
	Drugs                               // 毒品
	Profanity                           // 脏话
	Discrimination                      // 歧视
	Scam                                // 诈骗
	Custom                              // 自定义
)

func init() {
	// 注册静态分类
	registerStatic("未分类", None)
	registerStatic("涉黄", Pornography)
	registerStatic("涉政", Political)
	registerStatic("暴力", Violence)
	registerStatic("赌博", Gambling)
	registerStatic("毒品", Drugs)
	registerStatic("脏话", Profanity)
	registerStatic("歧视", Discrimination)
	registerStatic("诈骗", Scam)
	registerStatic("自定义", Custom)

	// 动态分类起始值 = 静态最大值 << 1
	nextDynamicVal = Custom << 1
}

// 注册静态分类
func registerStatic(name string, val Category) {
	categoryMap[name] = val
	categoryString[val] = name
	predefinedFlags |= val
}

// RegisterCategory 注册一个新分类（如果已存在则返回旧值）
func RegisterCategory(name string) Category {
	mu.Lock()
	defer mu.Unlock()

	if val, ok := categoryMap[name]; ok {
		return val
	}

	val := nextDynamicVal
	nextDynamicVal <<= 1
	categoryMap[name] = val
	categoryString[val] = name
	return val
}

// ParseCategory 获取分类值（从名称）
func ParseCategory(name string) (Category, bool) {
	mu.RLock()
	defer mu.RUnlock()
	val, ok := categoryMap[name]
	return val, ok
}

// String 返回分类名称
func (c Category) String() string {
	mu.RLock()
	defer mu.RUnlock()
	if name, ok := categoryString[c]; ok {
		return name
	}
	return fmt.Sprintf("未知分类(%d)", c)
}

// Contains 判断是否包含某分类（支持组合）
func (c Category) Contains(other Category) bool {
	if other == None {
		return c == None
	}
	if !other.IsValid() {
		return false
	}
	return c&other != 0
}

// IsValid 判断是否合法（静态或动态）
func (c Category) IsValid() bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := categoryString[c]
	return ok
}

// All 返回当前所有已注册的分类（含静态和动态）
func All() Category {
	mu.RLock()
	defer mu.RUnlock()

	var all Category
	for cat := range categoryString {
		if cat != None {
			all |= cat
		}
	}
	return all
}
