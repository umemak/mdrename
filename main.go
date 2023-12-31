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
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("FileInfo is nil for path: %s", path)
		}
		if strings.HasSuffix(path, ext) {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path %v: %v\n", dir, err)
		os.Exit(1)
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
			continue
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
			// ファイル名の重複チェック
			if _, err := os.Stat(newName); err == nil {
				fmt.Println(" - exists. skipped")
				continue
			}
			// ファイル名の長さがシステムの制限を超えていないかチェック
			if len(newName) > 255 {
				fmt.Println(" - too long. skipped")
				continue
			}
			fmt.Printf(" -> %q", newName)
			err = os.Rename(file, newName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(" - renamed")
		} else {
			fmt.Println(" - skipped")
		}
	}
}
