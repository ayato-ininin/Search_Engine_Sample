package common

import (
	"encoding/gob"
	"os"
)

// LoadInvertedIndex バイナリ形式で保存された転置インデックスを読み込む
func LoadInvertedIndex(path string) (map[string][]string, error) {
	var invertedIndex map[string][]string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)

	if err := dec.Decode(&invertedIndex); err != nil {
		return nil, err
	}

	return invertedIndex, nil
}
