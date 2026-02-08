package knowledge

import (
	"context"
	"sort"
)

// Item은 바이너리에 내장될 지식 조각입니다.
type Item struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// Data는 사용자가 직접 채워넣을 내장 지식 창고입니다.
var Data = []Item{
	// 예시 데이터: 사용자가 여기에 지식을 채워넣게 됩니다.
	{ID: "1", Title: "Lux Framework Overview", Content: "Lux is an AI-native framework...", Tags: []string{"lux", "overview", "go"}},
}

// SearchResult는 검색 결과와 매칭된 태그 개수를 담습니다.
type SearchResult struct {
	Item       Item
	MatchCount int
}

// SearchKnowledge는 태그 매칭 개수가 많은 순서대로 10개의 지식 타이틀을 반환합니다.
func SearchKnowledge(ctx context.Context, tags []string) []SearchResult {
	var results []SearchResult

	for _, item := range Data {
		matchCount := 0
		for _, t := range tags {
			for _, it := range item.Tags {
				if t == it {
					matchCount++
					break
				}
			}
		}

		if matchCount > 0 {
			results = append(results, SearchResult{
				Item:       item,
				MatchCount: matchCount,
			})
		}
	}

	// 매칭 개수 내림차순 정렬
	sort.Slice(results, func(i, j int) bool {
		return results[i].MatchCount > results[j].MatchCount
	})

	// 최대 10개 제한
	if len(results) > 10 {
		results = results[:10]
	}

	return results
}

// GetContentByID는 ID를 통해 지식의 상세 내용을 가져옵니다.
func GetContentByID(ctx context.Context, id string) (Item, bool) {
	for _, item := range Data {
		if item.ID == id {
			return item, true
		}
	}
	return Item{}, false
}
