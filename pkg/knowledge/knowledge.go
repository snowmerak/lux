package knowledge

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

// Item은 지식의 단위입니다.
type Item struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

//go:embed data.bleve
var embeddedIndex embed.FS

var (
	index     bleve.Index
	indexOnce sync.Once
)

// initIndex는 임베딩된 인덱스를 임시 디렉토리에 풀어서 Bleve로 엽니다.
func initIndex() error {
	var err error
	indexOnce.Do(func() {
		tmpDir, tErr := os.MkdirTemp("", "lux-index-*")
		if tErr != nil {
			err = tErr
			return
		}

		// 임베딩된 파일들을 임시 디렉토리로 복사
		err = copyEmbedToDisk(embeddedIndex, "data.bleve", tmpDir)
		if err != nil {
			return
		}

		index, err = bleve.Open(filepath.Join(tmpDir, "data.bleve"))
	})
	return err
}

func copyEmbedToDisk(fs embed.FS, srcDir, dstDir string) error {
	entries, err := fs.ReadDir(srcDir)
	if err != nil {
		return err
	}

	outDir := filepath.Join(dstDir, srcDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := path.Join(srcDir, entry.Name())
		dstPath := filepath.Join(outDir, entry.Name())

		if entry.IsDir() {
			if err := copyEmbedToDisk(fs, srcPath, dstDir); err != nil {
				return err
			}
			continue
		}

		srcFile, err := fs.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}

// SearchKnowledge는 Bleve를 사용하여 태그 및 내용을 검색합니다.
func SearchKnowledge(ctx context.Context, tags []string) ([]SearchResult, error) {
	if err := initIndex(); err != nil {
		return nil, err
	}

	// 태그들을 OR 쿼리로 결합 (tags와 title 필드 모두 검색)
	var queries []query.Query
	for _, t := range tags {
		tq := bleve.NewMatchQuery(t)
		tq.SetField("tags")
		queries = append(queries, tq)

		titleQ := bleve.NewMatchQuery(t)
		titleQ.SetField("title")
		queries = append(queries, titleQ)

		contentQ := bleve.NewMatchQuery(t)
		contentQ.SetField("content")
		queries = append(queries, contentQ)
	}

	if len(queries) == 0 {
		return make([]SearchResult, 0), nil
	}

	q := bleve.NewDisjunctionQuery(queries...)
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 10
	searchRequest.Fields = []string{"title"} // 소문자로 통일

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		title := getStringField(hit.Fields, "title")
		results = append(results, SearchResult{
			Item: Item{
				ID:    hit.ID,
				Title: title,
			},
			MatchCount: int(hit.Score),
		})
	}

	return results, nil
}

// SearchResult는 검색 결과를 담습니다.
type SearchResult struct {
	Item       Item
	MatchCount int
}

// GetContentByID는 ID를 통해 지식의 상세 내용을 가져옵니다.
func GetContentByID(ctx context.Context, id string) (Item, bool) {
	if err := initIndex(); err != nil {
		return Item{}, false
	}

	// 모든 필드를 가져오기 위해 검색 요청
	q := bleve.NewDocIDQuery([]string{id})
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Fields = []string{"Title", "Content", "Tags"}

	results, err := index.Search(searchRequest)
	if err != nil || results.Total == 0 {
		return Item{}, false
	}

	hit := results.Hits[0]
	item := Item{
		ID:      hit.ID,
		Title:   getStringField(hit.Fields, "title"),
		Content: getStringField(hit.Fields, "content"),
		Tags:    getStringSliceField(hit.Fields, "tags"),
	}

	return item, true
}

func getStringField(fields map[string]any, name string) string {
	if val, ok := fields[name]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func getStringSliceField(fields map[string]any, name string) []string {
	val, ok := fields[name]
	if !ok {
		return nil
	}

	switch v := val.(type) {
	case []any:
		var res []string
		for _, item := range v {
			res = append(res, fmt.Sprintf("%v", item))
		}
		return res
	case []string:
		return v
	case string:
		return []string{v}
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}
