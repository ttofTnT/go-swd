package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kirklin/go-swd/pkg/core"
	"github.com/kirklin/go-swd/pkg/types/category"
)

// mockDetector 是一个用于测试的mock检测器
type mockDetector struct {
	matchAllFunc   func(text string) []core.SensitiveWord
	matchAllInFunc func(text string, categories ...category.Category) []core.SensitiveWord
}

func (m *mockDetector) Detect(text string) bool {
	return len(m.MatchAll(text)) > 0
}

func (m *mockDetector) DetectIn(text string, categories ...category.Category) bool {
	return len(m.MatchAllIn(text, categories...)) > 0
}

func (m *mockDetector) Match(text string) *core.SensitiveWord {
	matches := m.MatchAll(text)
	if len(matches) > 0 {
		return &matches[0]
	}
	return nil
}

func (m *mockDetector) MatchIn(text string, categories ...category.Category) *core.SensitiveWord {
	matches := m.MatchAllIn(text, categories...)
	if len(matches) > 0 {
		return &matches[0]
	}
	return nil
}

func (m *mockDetector) MatchAll(text string) []core.SensitiveWord {
	if m.matchAllFunc != nil {
		return m.matchAllFunc(text)
	}
	return nil
}

func (m *mockDetector) MatchAllIn(text string, categories ...category.Category) []core.SensitiveWord {
	if m.matchAllInFunc != nil {
		return m.matchAllInFunc(text, categories...)
	}
	return nil
}

func TestFilter_Replace(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		replacement rune
		matches     []core.SensitiveWord
		want        string
	}{
		{
			name:        "empty text",
			text:        "",
			replacement: '*',
			matches:     nil,
			want:        "",
		},
		{
			name:        "no sensitive words",
			text:        "hello world",
			replacement: '*',
			matches:     nil,
			want:        "hello world",
		},
		{
			name:        "single sensitive word",
			text:        "hello bad world",
			replacement: '*',
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
			},
			want: "hello *** world",
		},
		{
			name:        "multiple sensitive words",
			text:        "bad hello bad world",
			replacement: '#',
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 0, EndPos: 3, Category: category.Violence},
				{Word: "bad", StartPos: 10, EndPos: 13, Category: category.Violence},
			},
			want: "### hello ### world",
		},
		{
			name:        "overlapping sensitive words",
			text:        "hello badword world",
			replacement: '*',
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
				{Word: "word", StartPos: 9, EndPos: 13, Category: category.Violence},
			},
			want: "hello ******* world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{
				matchAllFunc: func(text string) []core.SensitiveWord {
					return tt.matches
				},
			}
			f := NewFilter(detector)
			got := f.Replace(tt.text, tt.replacement)
			if got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_ReplaceIn(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		replacement rune
		categories  []category.Category
		matches     []core.SensitiveWord
		want        string
	}{
		{
			name:        "empty text",
			text:        "",
			replacement: '*',
			categories:  []category.Category{category.Violence},
			matches:     nil,
			want:        "",
		},
		{
			name:        "no categories",
			text:        "hello bad world",
			replacement: '*',
			categories:  nil,
			matches:     nil,
			want:        "hello bad world",
		},
		{
			name:        "single category match",
			text:        "hello bad world",
			replacement: '*',
			categories:  []category.Category{category.Violence},
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
			},
			want: "hello *** world",
		},
		{
			name:        "multiple categories",
			text:        "bad hello evil world",
			replacement: '#',
			categories:  []category.Category{category.Violence, category.Pornography},
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 0, EndPos: 3, Category: category.Violence},
				{Word: "evil", StartPos: 10, EndPos: 14, Category: category.Pornography},
			},
			want: "### hello #### world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{
				matchAllInFunc: func(text string, categories ...category.Category) []core.SensitiveWord {
					if !reflect.DeepEqual(categories, tt.categories) {
						t.Errorf("Categories mismatch, got %v, want %v", categories, tt.categories)
					}
					return tt.matches
				},
			}
			f := NewFilter(detector)
			got := f.ReplaceIn(tt.text, tt.replacement, tt.categories...)
			if got != tt.want {
				t.Errorf("ReplaceIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_ReplaceWithAsterisk(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		matches []core.SensitiveWord
		want    string
	}{
		{
			name:    "empty text",
			text:    "",
			matches: nil,
			want:    "",
		},
		{
			name:    "no sensitive words",
			text:    "hello world",
			matches: nil,
			want:    "hello world",
		},
		{
			name: "single sensitive word",
			text: "hello bad world",
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
			},
			want: "hello *** world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{
				matchAllFunc: func(text string) []core.SensitiveWord {
					return tt.matches
				},
			}
			f := NewFilter(detector)
			got := f.ReplaceWithAsterisk(tt.text)
			if got != tt.want {
				t.Errorf("ReplaceWithAsterisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_ReplaceWithStrategy(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		strategy func(word core.SensitiveWord) string
		matches  []core.SensitiveWord
		want     string
	}{
		{
			name: "custom replacement strategy",
			text: "hello bad world",
			strategy: func(word core.SensitiveWord) string {
				return "[REMOVED]"
			},
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
			},
			want: "hello [REMOVED] world",
		},
		{
			name: "length based replacement",
			text: "hello bad evil world",
			strategy: func(word core.SensitiveWord) string {
				return "?" + strings.Repeat("*", len([]rune(word.Word))-2) + "?"
			},
			matches: []core.SensitiveWord{
				{Word: "bad", StartPos: 6, EndPos: 9, Category: category.Violence},
				{Word: "evil", StartPos: 10, EndPos: 14, Category: category.Violence},
			},
			want: "hello ?*? ?**? world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{
				matchAllFunc: func(text string) []core.SensitiveWord {
					return tt.matches
				},
			}
			f := NewFilter(detector)
			got := f.ReplaceWithStrategy(tt.text, tt.strategy)
			if got != tt.want {
				t.Errorf("ReplaceWithStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}
