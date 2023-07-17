# 簡易全文検索エンジン in [Go](https://golang.org/dl/)

* 対象は直近1万件のエントリ(gitにはサンプルのみ)
* 検索語を含むエントリを返す
* 返す内容は、eid、タイトル、URL、スニペット

※書籍上、perlで書かれているものをgolangで記述

## 転置インデックス作成

MeCabを使用して形態素解析を行い、
keyValueMapをバイナリ形式で生成。

```bash
go run ./indexer
go run ./indexer -entriesDefaultPath ./indexer/10000entries/10000entries.txt -textsDefaultDir ./indexer/10000entries/texts
```

## Using GoProject

```bash
go run main.go
go run main.go -entriesDefaultPath ./indexer/10000entries/10000entries.txt -textsDefaultDir ./indexer/10000entries/text   
```

## References

This project was made from the following resources:

1. [[Web開発者のための]大規模サービス技術入門 ―データ構造、メモリ、OS、DB、サーバ/インフラ](https://www.amazon.co.jp/Web%E9%96%8B%E7%99%BA%E8%80%85%E3%81%AE%E3%81%9F%E3%82%81%E3%81%AE-%E5%A4%A7%E8%A6%8F%E6%A8%A1%E3%82%B5%E3%83%BC%E3%83%93%E3%82%B9%E6%8A%80%E8%A1%93%E5%85%A5%E9%96%80-%E2%80%95%E3%83%87%E3%83%BC%E3%82%BF%E6%A7%8B%E9%80%A0%E3%80%81%E3%83%A1%E3%83%A2%E3%83%AA%E3%80%81OS%E3%80%81DB%E3%80%81%E3%82%B5%E3%83%BC%E3%83%90-PRESS-plus%E3%82%B7%E3%83%AA%E3%83%BC%E3%82%BA/dp/4774143073)
