package common

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type DocumentInfo struct {
	DocID      string
	CategoryID string
	Url        string
	Title      string
}

func GetDocumentInfoMap(path string) map[string]DocumentInfo {
	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// ファイルから一行ずつ読み込み
	titles := make(map[string]DocumentInfo)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t") // タブで行を分割

		if len(parts) < 4 {
			continue
		}
		info := genDocumentInfoMap(parts)

		titles[info.DocID] = info
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return titles
}

func genDocumentInfoMap(parts []string) DocumentInfo {
	return DocumentInfo{
		DocID:      parts[0],
		CategoryID: parts[1],
		Url:        parts[2],
		Title:      parts[3],
	}
}
