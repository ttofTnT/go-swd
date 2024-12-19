package types

// Category 敏感词分类
type Category uint32

const (
	// None 未分类
	None Category = 0
	// Pornography 涉黄
	Pornography Category = 1 << iota
	// Political 涉政
	Political
	// Violence 暴力
	Violence
	// Gambling 赌博
	Gambling
	// Drugs 毒品
	Drugs
	// Profanity 谩骂
	Profanity
	// Discrimination 歧视
	Discrimination
	// Scam 诈骗
	Scam
	// Custom 自定义
	Custom
)

// CategoryNames 分类名称映射
var CategoryNames = map[Category]string{
	None:           "未分类",
	Pornography:    "涉黄",
	Political:      "涉政",
	Violence:       "暴力",
	Gambling:       "赌博",
	Drugs:          "毒品",
	Profanity:      "谩骂",
	Discrimination: "歧视",
	Scam:           "诈骗",
	Custom:         "自定义",
}

// String 获取分类名称
func (c Category) String() string {
	if name, ok := CategoryNames[c]; ok {
		return name
	}
	return "未知分类"
}

// Contains 判断是否包含指定分类
func (c Category) Contains(other Category) bool {
	return c&other != 0
}

// All 所有预定义分类
var All = Political | Pornography | Gambling | Violence | Discrimination | Drugs | Scam
