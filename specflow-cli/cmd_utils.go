package main

import (
	"os"
	"sort"
)

// readDirSorted 读取目录并按名称降序排列（最新的 journal-N.md 在前）
func readDirSorted(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}
