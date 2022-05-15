package scrapers

import (
	"bufio"
	"fmt"
	"os"
)

func GetProxies() []string {
	readFile, err := os.Open("proxies.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	for _, line := range fileLines {
		fmt.Println(line)
	}

	fmt.Println(fileLines)
	return fileLines
}
