package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"log"
)

type config struct {
	// extention to filter out
	ext string
	// min file size
	size int64
	// list files
	list bool
	del bool
	wLog io.Writer
	archive string
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	logFile := flag.String("log", "", "Log deletes to this file")
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "List files only")
	ext := flag.String("ext", "", "File extention to filter by")
	archive := flag.String("archive", "", "Archive Directory")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	var (
		f = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext: *ext,
		size: *size,
		list: *list,
		del: *del,
		wLog: f,
		archive: *archive,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error{
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)
	return filepath.Walk(root, 
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}

			if cfg.list {
				return listFile(path, out)
			}

			if cfg.archive != "" {
				if err := archiveFile(cfg.archive, root, path); err != nil {
					return err
				}
			}

			if cfg.del {
				return delFile(path, delLogger)
			}

			return listFile(path, out)
		},
	)
}
