package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"search_engine/common"
	"search_engine/variable_byte_code"
	"strings"
)

var docInfo map[string]common.DocumentInfo
var invertedIndex map[string][]string
var entriesDefaultPath string
var textsDefaultDir string

// indexerで転置インデックスを作成後、
// このプログラムを実行すると、全文検索ができる
func main() {
	var err error

	flag.StringVar(&entriesDefaultPath, "entriesDefaultPath", "./indexer/entries_sample.txt", "エントリーデフォルトパス")
	flag.StringVar(&textsDefaultDir, "textsDefaultDir", "./indexer/texts", "textが格納されているディレクトリ")
	flag.Parse()

	// docIdをキーにドキュメント情報MAPを作成
	docInfo = common.GetDocumentInfoMap(entriesDefaultPath)

	invertedIndex, err = common.LoadInvertedIndex("./index.bin")
	if err != nil {
		log.Fatalf("failed to load index.bin: %v", err)
	}
	intro()

	// create a channel to indicat e when the program can quit
	doneChan := make(chan bool)

	// start a gorouutin to read user input and run program
	go readUserInput(os.Stdin, doneChan)

	// block until the done chan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	// say goodbye
	fmt.Println("Goodbye!")
}

func intro() {
	fmt.Println("Enter a search text.")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)

	for {
		scanner.Scan()
		stdin := scanner.Text()
		// invertedIndexからscanner.Text()に対応するデータ取得
		if (stdin == "exit" || stdin == "q") {
			doneChan <- true
			return
		}
		data := invertedIndex[stdin]
		if len(data) == 0 {
			fmt.Println("Not found %s", stdin)
			prompt()
			continue
		}
		// vbcodeでデコード
		decodedEntries, err := variableByteCode.VariableByteDecode(data)
		if err != nil {
			log.Fatalf("failed to decode %v", err)
		}
		var res []string
		// decodedEntriesはdocIdの配列なので、ループして./indexer/texts/docIdを開く(limit 5件)
		for i, docId := range decodedEntries {
			if i == 5 {
				break
			}
			file, err := os.Open(fmt.Sprintf("%s/%s", textsDefaultDir, docId))
			if err != nil {
				log.Fatalf("failed to open %s: %v", docId, err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				// stdinが含まれる行を一行のみ抽出
				if strings.Contains(scanner.Text(), stdin) {
					docInfo := docInfo[docId]
					res = append(res, fmt.Sprintf("[%s]\n%s\n%s\n%s",
					docInfo.DocID, docInfo.Title, docInfo.Url, scanner.Text()))
					break
				}
			}
		}
		fmt.Println(strings.Join(res, "\n\n-----------------\n\n"))
		prompt()
	}
}
