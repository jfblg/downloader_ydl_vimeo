package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type video struct {
	name        string
	filename    string
	url         string
	urlDownload string
}

func newVideo(name, url string) video {
	v := video{name: name, url: url}
	v.filename = ""
	v.urlDownload = ""
	return v
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s file\n",
		path.Base(os.Args[0]))
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	dryrun := flag.Bool("d", false, "dryrun only")
	// verbose := flag.Bool("v", false, "verbose mode")
	flag.Usage = usage
	flag.Parse()

	if *dryrun {
		fmt.Printf("Dry run mode\n")
	}

	if len(flag.Args()) < 1 {
		usage()
	}

	srcFile := flag.Args()[0]

	var vList []video
	readFileContent(srcFile, &vList)

	for _, v := range vList {
		v.filename, _ = transFile(v.name)
		v.urlDownload, _ = transURL(v.url)

		log.Printf("Downloading: %v\n", v.name)

		// download file
		runYoutubeDl(v.urlDownload)

		// rename downloaded file
		err := fileFinderAndRenamer("master.*", v.filename)
		if err != nil {
			panic(err)
		}
	}
}

func readFileContent(filename string, vList *[]video) {

	// var vListLocal []video

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var v video
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// empty line - do nothing
			continue
		} else if strings.HasPrefix(line, "#") {
			// commented line - do nothing
			continue
		} else if strings.HasPrefix(line, "http") {
			// URL
			v.url = line
		} else {
			// Video name
			v.name = line
		}
		if v.name != "" && v.url != "" {
			*vList = append(*vList, v)
			v = video{}
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%v loaded. %d URLs found for download\n", filename, len(*vList))
}

// just for finding files
func fileFinder(pattern string) {
	var rex *regexp.Regexp

	if len(pattern) == 0 {
		pattern = "master.*"
	}

	rex = regexp.MustCompile(pattern)

	dir := "."
	files := []string{}
	walk := func(fn string, fi os.FileInfo, err error) error {
		if rex.MatchString(fn) == false {
			return nil
		}
		if fi.IsDir() {
			fmt.Println(fn + string(os.PathSeparator))
		} else {
			fmt.Println(fn)
			files = append(files, fn)
		}
		return nil
	}

	filepath.Walk(dir, walk)
	fmt.Printf("Found files: %[1]d \n", len(files))
	for _, v := range files {
		fmt.Printf("%v\n", v)
	}
}

func fileFinderAndRenamer(pattern string, newFilename string) error {
	var rex *regexp.Regexp

	if len(pattern) == 0 {
		pattern = "master.*"
	}

	rex = regexp.MustCompile(pattern)

	dir := "."
	files := []string{}
	walk := func(fn string, fi os.FileInfo, err error) error {
		if rex.MatchString(fn) == false {
			return nil
		}
		if fi.IsDir() {
			fmt.Println(fn + string(os.PathSeparator))
		} else {
			fmt.Println(fn)
			files = append(files, fn)
		}
		return nil
	}

	filepath.Walk(dir, walk)
	fmt.Printf("Found files: %[1]d \n", len(files))
	if len(files) == 0 {
		log.Printf("No files to rename found. Pattern used: %+v\n", pattern)
	}
	if len(files) > 1 {
		log.Printf("More than 1 file found. Renaming just %v\n", files[0])
	}
	return os.Rename(files[0], newFilename)
}

func runYoutubeDl(url string) {
	cmd := exec.Command(
		"youtube-dl",
		"--quiet",
		"--ignore-config",
		url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File download completed")
	// fmt.Printf("command output %q\n", out.String())
}

func transFile(name string) (string, error) {
	// TODO return error if something is wrong in filename

	words := strings.Split(name, " ")
	indexes := strings.Split(words[0], ".")

	index := strings.Join(indexes, "-")
	word := strings.Join(words[1:], "_")

	filename := "yoga15_" + index + "_" + word + ".mp4"

	return filename, nil
}

func transURL(url string) (string, error) {

	splittedByDots := strings.Split(url, ".")

	if len(splittedByDots) != 4 {
		return "", errors.New("Something wrong with the URL")
	}

	if splittedByDots[3] != "json?base64_init=1" {
		return "", errors.New("Unexpexted ending of the URL")
	}

	newURL := strings.Join(splittedByDots[:3], ".") + ".mpd"

	return newURL, nil
}
