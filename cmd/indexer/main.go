package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/snowmerak/lux/v4/pkg/knowledge"
)

func main() {
	indexPath := "pkg/knowledge/data.bleve"

	// Remove existing index if it exists
	if _, err := os.Stat(indexPath); err == nil {
		os.RemoveAll(indexPath)
	}
	os.MkdirAll(indexPath, 0755)
	os.WriteFile(filepath.Join(indexPath, "placeholder.txt"), []byte("placehold"), 0644)

	// Create a new Bleve mapping
	mapping := bleve.NewIndexMapping()

	// 필드들이 검색 결과에 포함되도록 명시적 매핑 설정
	itemMapping := bleve.NewDocumentMapping()

	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Store = true
	itemMapping.AddFieldMappingsAt("title", titleMapping)

	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Store = true
	itemMapping.AddFieldMappingsAt("content", contentMapping)

	tagMapping := bleve.NewTextFieldMapping()
	tagMapping.Store = true
	itemMapping.AddFieldMappingsAt("tags", tagMapping)

	mapping.AddDocumentMapping("_default", itemMapping)

	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	// Walk through the knowledge directory
	err = filepath.Walk("knowledge", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		item, err := parseMarkdown(path)
		if err != nil {
			log.Printf("Failed to parse %s: %v", path, err)
			return nil
		}

		err = index.Index(item.ID, item)
		if err != nil {
			return err
		}

		fmt.Printf("Indexed: %s\n", item.Title)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Indexing completed successfully.")
}

func parseMarkdown(path string) (knowledge.Item, error) {
	file, err := os.Open(path)
	if err != nil {
		return knowledge.Item{}, err
	}
	defer file.Close()

	var item knowledge.Item
	item.ID = filepath.Base(path)

	scanner := bufio.NewScanner(file)
	var content strings.Builder
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if lineNum == 0 && strings.HasPrefix(line, "# ") {
			item.Title = strings.TrimPrefix(line, "# ")
		} else if lineNum == 1 && strings.HasPrefix(line, "Tags: ") {
			tagsStr := strings.TrimPrefix(line, "Tags: ")
			for _, t := range strings.Split(tagsStr, ",") {
				item.Tags = append(item.Tags, strings.TrimSpace(t))
			}
		} else {
			content.WriteString(line + "\n")
		}
		lineNum++
	}
	item.Content = content.String()
	return item, nil
}
