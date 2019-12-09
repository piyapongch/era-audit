package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main2() {
	var files []string
	var count int

	root := "/Users/pcharoen/Documents/era-audit/openstack-swift/objects"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
		count++
	}
	fmt.Printf("#%d %s\n", count, "objects")
}
