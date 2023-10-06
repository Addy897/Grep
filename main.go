package main

import (
	"bufio"
	"sync"

	"fmt"
	"os"
	"strings"
)

var ext string = ""
var th int = 0
var MAX_TH int = 999

func search(name string, tag string, wg *sync.WaitGroup) {

	fp, err := os.Open(name)
	defer fp.Close()
	if err != nil {
		fmt.Print(err)
		th--
		(*wg).Done()
		return
	}
	c := 1
	scan := bufio.NewScanner(fp)
	scan.Split(bufio.ScanLines)
	for scan.Scan() {

		line := scan.Text()
		if strings.Contains(line, tag) {
			fmt.Printf("Found %s in %s at Line %d\n", tag, name, c)
		}

		c++
	}
	th--
	(*wg).Done()
}
func countFiles(name string, tag string, wg *sync.WaitGroup) {
	nfiles, err := os.ReadDir(name)
	if err != nil {
		//fmt.Println(err.Error())
		th--
		(*wg).Done()
		return
	}
	var files []string
	var dir []string
	for _, v := range nfiles {
		currentDir := name + v.Name() + "\\"
		if !v.IsDir() {
			if ext != "" {
				if strings.Contains(v.Name(), ext) {
					files = append(files, name+v.Name())
				}
				continue
			}
			if !strings.Contains(v.Name(), ".go") && !strings.Contains(v.Name(), ".exe") {
				files = append(files, name+v.Name())

			}

		} else {
			dir = append(dir, currentDir)

		}
	}

	if len(dir) > 0 {
		(*wg).Add(len(dir))
		for i := 0; i < len(dir); i++ {
			for th > MAX_TH {
				continue
			}
			th++
			go countFiles(dir[i], tag, wg)
		}
	}

	(*wg).Add(len(files))
	for i := 0; i < len(files); i++ {
		for th > MAX_TH {
			continue
		}
		th++
		go search(files[i], tag, wg)
	}
	th--
	(*wg).Done()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("[-]Usage: %s  [Filename] [File Name Or File Extension(Optional)]", os.Args[0])
		return
	} else {
		if len(os.Args) == 3 {
			ext = os.Args[2]
		}
		tag := os.Args[1]
		var wg sync.WaitGroup
		wg.Add(1)
		countFiles(".\\", tag, &wg)
		wg.Wait()

	}
}
