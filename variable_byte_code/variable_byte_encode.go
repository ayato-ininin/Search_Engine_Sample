package variableByteCode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

func VariableByteEncode(entries []string) ([]string, error) {
	var err error
	entryNums := make([]uint64, len(entries))
	// エントリをuint64に変換し、リストをソート
	for i, entry := range entries {
		entryNums[i], err = strconv.ParseUint(entry, 10, 64)
		if err != nil {
			fmt.Println("Error converting entry to uint64:", err)
			return nil, err
		}
	}
	// エントリをソート
	sort.Slice(entryNums, func(i, j int) bool { return entryNums[i] < entryNums[j] })

	// ギャップエンコーディングを適用
	var encodedEntries []string
	var prevEntry uint64 = 0
	for _, entryNum := range entryNums {
		// エントリをVByteエンコード
		gap := entryNum - prevEntry
		prevEntry = entryNum
		encoded := vByteEncode(gap)
		// エンコードされたエントリを16進数文字列に変換→あまり必要ないかも？
		encodedEntries = append(encodedEntries,fmt.Sprintf("%x", encoded))
	}
	return encodedEntries, nil
}

//Variable Byte Encoding（VByteエンコーディング）は、非負整数をバイト配列として効率的に格納するための方法。
//数値を7ビットのグループに分割し、最後のバイト以外のすべてのバイトに「続くバイトがある」ことを示すために最上位ビットを設定。
func vByteEncode(n uint64) []byte {
	var b bytes.Buffer
	for n >= 0x80 {  // nが128以上である限りループを続ける(0x80: 16進数で128)
		// nの下位7ビットをバッファに書き込む
		b.WriteByte(byte(n) | 0x80) // byte()は下位8ビットを取り出す, 255までしか表現できない、かつ最上位ビットを1にする
		n >>= 7  // nを7ビット右シフトして、次の7ビットを処理する準備をする(左側に0が詰められる)
	}
	b.WriteByte(byte(n))  // 最後の7ビット（もしくはそれ以下）をバッファに書き込む
	return b.Bytes()
}
