package razutils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
Package goUtils implements various simple functions used in my personal programs.
The package has no license, use at your discretion, no warranty is given on how good those functions work

*/

// var ext = []string{".avi", ".mkv", ".mp4", ".mov"}
var container = regexp.MustCompile(`(?i)\.(?:MKV|AVI|MP4|MOV|MPG|MPEG|FLV|F4V|SWF|WMV|MP2|MPE|MPV|OGG|M4V|M4P|AVCHD)$`)

// FileExists check if a file/directory exist at the given path.  Note that if there is an access issue the function will
// return false,error but the file might exist.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// FileParts returns the breakdown of a filepath to path, filename and extension.
// the extension includes the dot.
func FileParts(path string) (dir string, file string, ext string) {
	dir, f := filepath.Split(path)
	ext = filepath.Ext(f)
	file = f[0 : len(f)-len(ext)]
	return dir, file, ext
}

// IsVideoFile check if a given file path is a valid video file name.
// note: the check is done by extension and not by file content.
func IsVideoFile(path string) bool {
	return container.MatchString(path)
}

// IsFileExt - check if a given path has a specific extension. (case ignored)
func IsFileExt(path string, ext string) bool {
	return strings.ToLower(filepath.Ext(path)) == ext
}

// IsSrtFile check if a given file has srt extension
func IsSrtFile(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".srt"
}

// ReplaceExt replace the file extension with a new extension.
// if an empty new extension is specified the existing extension will be removed
// if the new extension does not start with a ., it will be added
func ReplaceExt(path string, newExt string) string {
	if len(path) == 0 {
		log.Fatal("bad call to ReplaceExt - path empty")
	}
	p1 := path[0 : len(path)-len(filepath.Ext(path))]
	if len(newExt) == 0 {
		return p1
	}
	if newExt[:1] != "." {
		newExt = "."
	}
	return p1 + newExt
}

// RandFileName - return a random file name with an extension as mentioned in extension and prefix as in prefix
// the function will try 50 names before giving up.  It should be assumed that those files will be deleted as repeated calls
// will cause for sure failure in the long run.
func RandFileName(path string, prefix string, ext string) string {
	for cnt := 0; cnt <= 50; cnt++ {
		r := rand.Intn(99999)
		fn := prefix + strconv.Itoa(r) + "." + ext
		if _, err := os.Stat(filepath.Join(path, fn)); errors.Is(err, os.ErrNotExist) {
			return fn
		}
	}
	log.Fatal("Cant find a random file name for path ", path, " Prefix ", prefix, " Suffix ", ext)
	return ""
}

// Abs return an abs(int)
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CopyFile - copy a file from source to destination path.
func CopyFile(src string, dst string) error {
	// Read all content of src to data
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// MoveFile - move/rename a file from source to destination path.
// currently implemented as a copy+delete.  this is not optimal as same volume rename should be quicker
// however checking this is more complex
func MoveFile(src string, dst string) error {
	// Read all content of src to data
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	err = os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}

const chunkSize = 64000

// DeepCompare Compare two files to see if content is the same.
// The files are read by chunks and the first difference cause the function to return false.
func DeepCompare(file1, file2 string) bool {
	f1s, err := os.Stat(file1)
	if err != nil {
		log.Fatal(err)
	}
	f2s, err := os.Stat(file2)
	if err != nil {
		log.Fatal(err)
	}
	if f1s.Size() != f2s.Size() {
		return false
	}
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)
		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal("can not access files in DeepCompare: ", err1, err2)
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

// GzipExtract - convert a .gz by expanding it into the original file. Source is the gz file path, dest is what the
// result filename should be
func GzipExtract(source string, dest string) error {
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	reader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer reader.Close()
	fout, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fout.Close()
	_, err = io.Copy(fout, reader)
	return err
}

// DaysSince - computer round number of days between now and specified time in the past (or future)
func DaysSince(t time.Time) int {
	return int(math.Round(math.Abs(time.Now().Sub(t).Hours()) / 24))
}
