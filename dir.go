package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getFilelist(path string, suffix string) []string {
	var imgs []string

	suf_names := strings.Split(suffix, `;`)
	log.Println(`图片后缀名包括:`, suf_names)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		for _, suf := range suf_names {
			ok := len(suf) > 0 && strings.HasSuffix(f.Name(), suf)
			if ok {
				imgs = append(imgs, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return imgs
}

func LoadImgs(path string, optSuffix string) []string {
	return getFilelist(path, optSuffix)
}
