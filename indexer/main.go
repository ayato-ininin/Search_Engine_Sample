package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"search_engine/variable_byte_code"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	var entriesDefaultPath string
	var textsDefaultDir string

	flag.StringVar(&entriesDefaultPath, "entriesDefaultPath", "./indexer/entries_sample.txt", "エントリーデフォルトパス")
	flag.StringVar(&textsDefaultDir, "textsDefaultDir", "./indexer/texts", "textが格納されているディレクトリ")
	flag.Parse()

	// titleとdocIdのマップを作成
	titles := getTitleDocIdMap(entriesDefaultPath)
	// 転置インデックスを作成
	invertedIndex := createInvertedIndex(textsDefaultDir, titles)

	// 転置インデックスをシリアライズしてバイナリファイルに書き込む
	indexFile, err := os.Create("./index.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()

	// エンコーダの作成
	enc := gob.NewEncoder(indexFile)

	// マップのエンコードと書き込み
	if err := enc.Encode(invertedIndex); err != nil {
		log.Fatal("Encode error: ", err)
	}

	fmt.Println("Create index.bin successfully!")
}

// titleとdocIdのマップを作成
func getTitleDocIdMap(path string) map[string]string {
		// ファイルを開く
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// ファイルから一行ずつ読み込み
		titles := make(map[string]string)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, "\t")  // タブで行を分割

			if len(parts) < 4 {
				continue
			}

			docID := parts[0]
			title := parts[3]
			titles[docID] = title
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return titles
}

func createInvertedIndex(path string, titles map[string]string) map[string][]string {
	// 本文データがあるディレクトリを開く
	// ディレクトリ内のファイルを一つずつ処理
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// 転置インデックスを作成
	invertedIndex := make(map[string][]string)
	for _, file := range files {
		title := titles[file.Name()]  // ファイル名からタイトルを取得
		// fileの本文を読み込む
		filePath := fmt.Sprintf("%s/%s", path, file.Name())
		document, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}
		// titleとdocumentを結合
		text := fmt.Sprintf("%s\n%s", title, document)
		seg := getWakati(text)

		// 転置インデックスを作成
		for _, word := range seg {
			// 改行や半角全角空白、タブは無視(関数で文字出ないときは無視)
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}
			if _, ok := invertedIndex[word]; !ok {
				invertedIndex[word] = []string{file.Name()}
			} else {
				invertedIndex[word] = append(invertedIndex[word], file.Name())
			}
		}
	}
	return sortAndEncodeIndex(invertedIndex)
}

// textを形態素解析して分かち書きを取得
func getWakati(text string) []string {
	// textを形態素解析
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	// 分かち書き取得
	seg := t.Wakati(text)
	return seg
}

// 転置インデックスをソートした後にvbエンコードしたものを返す
// ギャップエンコード
func sortAndEncodeIndex(invertedIndex map[string][]string) map[string][]string {
	// 転置インデックスをvbエンコード
	for word, entries := range invertedIndex {
		encodedEntries, err := variablebytecode.VariableByteEncode(entries)
		if err != nil {
			log.Fatalf("Error encoding entries: %v", err)
		}
		invertedIndex[word] = encodedEntries
	}
	return invertedIndex
}
