package modules

import (
	"compress/flate"
	"encoding/json"
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

func Copy(configName, fileName string) error {
	var c fileConfig
	str, _ := file.Read(configName)
	err := json.Unmarshal([]byte(str), &c)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(fileName)
	str = xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)
	return file.Copy(fileName, str)
}

func Delete(configName, fileName string) error {
	return file.Delete(fileName)
}

func Move(configName, fileName string) error {
	var c fileConfig
	str, _ := file.Read(configName)
	err := json.Unmarshal([]byte(str), &c)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(fileName)
	str = xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)
	return file.Rename(fileName, str)
}

func Decompress(configName, fileName string) error {
	var c fileConfig
	str, _ := file.Read(configName)
	if err := json.Unmarshal([]byte(str), &c); err != nil {
		return err
	}

	_, filename := filepath.Split(fileName)
	str = xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)

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

func Compress(configName, fileName string) error {
	var c fileConfig
	str, _ := file.Read(configName)
	err := json.Unmarshal([]byte(str), &c)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(fileName)
	str = xtime.Fmt(c.Name, time.Now())
	str = strings.Replace(str, "$file", "%v", -1)
	str = fmt.Sprintf(str, filename)

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

func IsFile(configName, fileName string) error {
	if file.IsFile(fileName) {
		return nil
	}
	return fmt.Errorf("\"%v\" does not exist anymore\n", fileName)
}
