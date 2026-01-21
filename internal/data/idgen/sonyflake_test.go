package idgen

import (
	"testing"
)

func TestNewIdGenerator(t *testing.T) {
	generator := NewIDGenerator()
	if generator == nil {
		t.Fatal("Expected generator to not be nil")
	}
	if generator.idgen == nil {
		t.Error("Expected sonyflake instance to not be nil")
	}
}

func TestNextID(t *testing.T) {
	generator := NewIDGenerator()

	// 测试基本功能
	prefix := "user_"
	id := generator.NextID(prefix)
	t.Log("Generated ID:", id)

	if len(id) <= len(prefix) {
		t.Errorf("Generated ID %q is not longer than prefix %q", id, prefix)
	}

	if !startsWith(id, prefix) {
		t.Errorf("Generated ID %q does not start with prefix %q", id, prefix)
	}
	nextID, _ := generator.idgen.NextID()
	decompose := generator.idgen.Decompose(nextID)
	t.Log("nextID", generator.NextOrderID())
	t.Log("parse", decompose)
	// 提取数字部分并验证是否为有效数字
	numStr := id[len(prefix):]
	t.Log("len(numStr)", len(numStr))
	if numStr == "" {
		t.Error("Generated ID has no numeric part after prefix")
	}
}

func TestNextIDUniqueness(t *testing.T) {
	generator := NewIDGenerator()
	count := 100
	ids := make(map[string]bool, count)

	for i := 0; i < count; i++ {
		id := generator.NextID("test_")
		if ids[id] {
			t.Fatalf("Duplicate ID generated: %s", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

func TestNextIDMonotonicity(t *testing.T) {
	generator := NewIDGenerator()

	// 生成多个ID并验证它们是递增的
	prevID := int64(-1)
	for i := 0; i < 100; i++ {
		fullID := generator.NextID("")
		currentID := parseID(fullID)

		if currentID <= prevID {
			t.Errorf("IDs are not monotonically increasing: %d <= %d", currentID, prevID)
		}
		prevID = currentID
	}
}

func BenchmarkNextID(b *testing.B) {
	generator := NewIDGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = generator.NextID("bench_")
	}
}

// 辅助函数
func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func parseID(s string) int64 {
	// 从带前缀的ID中提取数字部分
	underscoreIndex := -1
	for i, r := range s {
		if r == '_' {
			underscoreIndex = i
			break
		}
	}
	if underscoreIndex == -1 {
		// 如果没有下划线，则整个字符串都是数字
		underscoreIndex = -1
	} else {
		underscoreIndex++
	}

	numStr := s[underscoreIndex:]
	if numStr == "" {
		return 0
	}

	var result int64
	for _, r := range numStr {
		if r >= '0' && r <= '9' {
			result = result*10 + int64(r-'0')
		} else {
			break
		}
	}
	return result
}
