package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func processDir(dir string, writeOnly, commentsOnly bool) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			processFile(path, writeOnly, commentsOnly)
		}

		return nil
	})

	if err != nil {
		errChan <- err
	}
}

func processFile(filename string, writeOnly, commentsOnly bool) {
	var (
		fw  *os.File
		err error
	)

	_, err = os.Stat(filename)
	if err != nil {
		errChan <- err
		return
	}

	if writeOnly {
		handleWriteFile(fw, filename, commentsOnly)
	} else {
		handleReadFile(fw, filename, commentsOnly)
	}
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = dst.Close()
	}()

	return io.Copy(dst, src)
}

func handleReadFile(fw *os.File, filename string, commentsOnly bool) {
	// 将每行的处理结果打在终端上
	fw = os.Stdout
	err := SpacingFile(filename, fw, commentsOnly)
	if err != nil {
		errChan <- err
		return
	}
}

func handleWriteFile(fw *os.File, filename string, commentsOnly bool) {
	// 将当前文件备份
	backFilename := "/var/tmp/" + filepath.Base(filename)
	_, err := copyFile(backFilename, filename)
	if err != nil {
		errChan <- err
		return
	}

	// 新建一个空文件, 用于临时存放处理完成后的文件
	newFilename := "/var/tmp/" + filepath.Base(filename) + ".readable"
	fw, err = os.Create(newFilename)
	if err != nil {
		errChan <- err
		return
	}
	defer func() {
		_ = fw.Close()
	}()

	// 处理文件, 添加空格
	err = SpacingFile(filename, fw, commentsOnly)
	if err != nil {
		errChan <- err
		return
	}

	// 将新的临时文件重命名为当前文件名
	err = os.Rename(newFilename, filename)
	if err != nil {
		errChan <- err
		return
	}

	// 删除备份文件
	_ = os.Remove(backFilename)
}

// SpacingFile reads the file named by filename, performs paranoid text
// spacing on its contents and writes the processed content to w.
// A successful call returns err == nil.
func SpacingFile(filename string, w io.Writer, commentsOnly bool) error {
	fr, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = fr.Close()
	}()

	br := bufio.NewReader(fr)
	bw := bufio.NewWriter(w)

	for {
		line, err := br.ReadString('\n')
		if err == nil {
			_, _ = fmt.Fprint(bw, Spacing(line, commentsOnly))
		} else {
			if err == io.EOF {
				_, _ = fmt.Fprint(bw, Spacing(line, commentsOnly))
				break
			}
			return err
		}
	}
	defer func() {
		_ = bw.Flush()
	}()

	return nil
}
