package modules

import (
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/xtime"
	"strings"
	"time"
)

type fileConfig struct {
	Name string `json:"name"`
}

type baseFile struct{}

func (baseFile) EmptyConfig() interface{} {
	return &fileConfig{}
}

func formatName(c fileConfig, fileName string) string {
	_, filename := filepath.Split(fileName)
	str := xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)
	return str
}

func Copy(config interface{}, fileName string) error {
	c := config.(fileConfig)

	str := formatName(c, fileName)

	return file.Copy(fileName, str)
}

type Delete struct {
	baseFile
}

func (Delete) Name() string {
	return "delete"
}

func (Delete) Perform(_ interface{}, fileName string) error {
	return file.Delete(fileName)
}

type Move struct {
	baseFile
}

func (Move) Name() string {
	return "move"
}

func (Move) Perform(config interface{}, fileName string) error {
	c := config.(fileConfig)

	_, filename := filepath.Split(fileName)
	str := xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)
	return file.Rename(fileName, str)
}

type Decompress struct {
	baseFile
}

func (Decompress) Name() string {
	return "decompress"
}

func (Decompress) Perform(config interface{}, fileName string) error {
	c := config.(fileConfig)

	str := formatName(c, fileName)

	i, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer i.Close()

	f := flate.NewReader(i)
	defer f.Close()

	o, err := os.Create(str)
	if err != nil {
		return err
	}
	defer o.Close()

	if _, err = io.Copy(o, f); err != nil {
		return err
	}

	return nil
}

type Compress struct {
	baseFile
}

func (Compress) Name() string {
	return "compress"
}

func (Compress) Perform(config interface{}, fileName string) error {
	c := config.(fileConfig)

	str := formatName(c, fileName)

	i, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer i.Close()

	o, err := os.Create(str)
	if err != nil {
		return err
	}
	defer o.Close()

	f, err := flate.NewWriter(o, flate.BestCompression)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, i); err != nil {
		return err
	}

	return nil
}

type IsFile struct {
	baseFile
}

func (IsFile) Name() string {
	return "isfile"
}

func (IsFile) Perform(config interface{}, fileName string) error {
	if file.IsFile(fileName) {
		return nil
	}
	return fmt.Errorf("\"%v\" does not exist anymore\n", fileName)
}
