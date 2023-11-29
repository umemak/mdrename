package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// ディレクトリの指定はコマンドライン引数で行う
	// コマンドライン引数が指定されていない場合はカレントディレクトリを対象とする
	dir := "." // カレントディレクトリを指定する
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	ext := ".md" // 対象の拡張子を指定する

	fileList := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ext) {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	re := regexp.MustCompile(`productTitle:\s*(.+)`) // titleメタ情報を抽出する正規表現

	for _, file := range fileList {
		fmt.Printf("%q", file)
		// ファイル名に" - "が含まれている場合は処理をスキップする
		if strings.Contains(file, " - ") {
			fmt.Println(" - skipped")
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}
		match := re.FindStringSubmatch(string(content))
		if len(match) > 1 {
			title := strings.Trim(strings.TrimSpace(match[1]), "\"")
			newName := fmt.Sprintf("%s - %s.md", strings.TrimSuffix(file, ext), title)
			// ファイル名に使えない文字を置換する
			newName = strings.ReplaceAll(newName, "/", "／")
			newName = strings.ReplaceAll(newName, ":", "：")
			newName = strings.ReplaceAll(newName, "*", "＊")
			newName = strings.ReplaceAll(newName, "?", "？")
			newName = strings.ReplaceAll(newName, "\"", "”")
			newName = strings.ReplaceAll(newName, "<", "＜")
			newName = strings.ReplaceAll(newName, ">", "＞")
			newName = strings.ReplaceAll(newName, "|", "｜")

			fmt.Printf(" -> %q", newName)
			err = os.Rename(file, newName)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(" - renamed")
		} else {
			fmt.Println(" - skipped")
		}
	}
}
