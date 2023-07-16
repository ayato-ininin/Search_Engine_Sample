package variableByteCode

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func VariableByteDecode(entries []string) ([]string, error) {
	decodedEntries := make([]string, len(entries))
	prevEntry := uint64(0)
	for i, encodeEntry := range entries {
		// encodeされた16進数エントリをバイトスライスに変換
		encoded, err := hex.DecodeString(encodeEntry)
		if err != nil {
			fmt.Println("Error decoding entry:", err)
			return nil, err
		}

		// バイトスライスをデコードして元の値を取得
		gap, err := vByteDecode(encoded)
		if err != nil {
			fmt.Println("Error decoding entry:", err)
			return nil, err
		}

		// ギャップ値を元のエントリに変換
		entry := prevEntry + gap
		prevEntry = entry

		// デコードされたエントリを文字列に変換
		decodedEntries[i] = strconv.FormatUint(entry, 10)
	}
	return decodedEntries, nil
}

// []byte{0xB8, 0x9E, 0x03}等くるので一つずつ処理
func vByteDecode(data []byte) (uint64, error) {
	var result uint64
	var shift uint
	for _, b := range data {
		// 1バイトずつ読み込んで、下位7ビットをresultに追加
		// 0x7F: 2進数で0111 1111→下位7ビットを取り出すためのマスク
		// b&0x7F: 論理積(0-0 → 0, 0-1 → 0, 1-0 → 0, 1-1 → 1)なので先頭は必ず0になる、下位7ビットが取り出せる
		// shift: 7ビットずつシフトしていく
		// |=: 論理和、+=と意味同じ
		result |= uint64(b&0x7F) << shift
		// 0x80: 2進数で1000 0000→最上位ビットを取り出すためのマスク
		// b&0x80: 論理積(0-0 → 0, 0-1 → 0, 1-0 → 0, 1-1 → 1)なので先頭は必ず0になる、最上位ビットが取り出せる
		// 0x80 == 0: 最上位ビットが0かどうかを判定
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("invalid vByte encoding")
}
