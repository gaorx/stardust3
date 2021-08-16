package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gaorx/stardust3/sdcompress"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdfile"
	"github.com/gaorx/stardust3/sdload"
)

func main() {
	var jarUrl string
	flag.StringVar(&jarUrl, "jar", "", "The url of 'com.xiaomi.ad:ip-utils_2.10'")
	flag.Parse()
	if jarUrl == "" {
		printUsage()
		return
	}

	// 加载com.xiaomi.ad:ip-utils_2.10的jar包
	jarBytes, err := sdload.Bytes(jarUrl)
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}

	// 从jar包中读取所需数据
	recordBytes, err := readRecords(jarBytes)
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}
	locBytes, err := readLocs(jarBytes)
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}

	// 检测目标目录是否存在
	sdIpLocDir, err := filepath.Abs("./sdiploc")
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}
	if !sdfile.IsDir(sdIpLocDir) {
		panic(sderr.Newf("not found source directory '%s'", sdIpLocDir))
	}

	err = ioutil.WriteFile(filepath.Join(sdIpLocDir, "rec.csv"), recordBytes, 0600)
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}
	err = ioutil.WriteFile(filepath.Join(sdIpLocDir, "loc.csv"), locBytes, 0600)
	if err != nil {
		panic(sderr.WithStack(err))
		return
	}
}

func printUsage() {
	fmt.Println("go run scripts/update_iploc_data.go -jar=<url/file for ip-utils.jar>")
}

func readRecords(jarBytes []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(jarBytes), int64(len(jarBytes)))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	var gzipRecordBytes []byte = nil
	for _, f := range zipReader.File {
		if strings.HasPrefix(f.Name, "ipv4_district_") {
			f1, err := f.Open()
			if err != nil {
				return nil, sderr.WithStack(err)
			}
			gzipRecordBytes, err = ioutil.ReadAll(f1)
			if err != nil {
				_ = f1.Close()
				return nil, sderr.WithStack(err)
			}
			_ = f1.Close()
			break
		}
	}
	if gzipRecordBytes == nil {
		return nil, sderr.New("not found ip_district_xxxxxxxx.csv.gz")
	}
	recordBytes, err := sdcompress.Ungzip(gzipRecordBytes)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return recordBytes, nil
}

func readLocs(jarBytes []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(jarBytes), int64(len(jarBytes)))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	var gzipLocBytes []byte = nil
	for _, f := range zipReader.File {
		if strings.HasPrefix(f.Name, "admin_hierarchy_") {
			f1, err := f.Open()
			if err != nil {
				return nil, sderr.WithStack(err)
			}
			gzipLocBytes, err = ioutil.ReadAll(f1)
			if err != nil {
				_ = f1.Close()
				return nil, sderr.WithStack(err)
			}
			_ = f1.Close()
			break
		}
	}
	if gzipLocBytes == nil {
		return nil, sderr.New("admin_hierarchy_xxxxxxxx.csv")
	}
	return gzipLocBytes, nil
}
