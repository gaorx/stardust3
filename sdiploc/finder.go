package sdiploc

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gaorx/stardust3/sdcsv"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdparse"
)

type Finder struct {
	recs []Record
	locs map[int64]Loc
}

// 中国行政区划: administrative division
//
// 省级:
//	 省: province
//	 自治区:autonomous region
//	 直辖市: municipality/centrally-administered municipality
//	 特别行政区: special administrative region/SAR
//
// 市级:
//	 地区: prefecture
//	 自治州: autonomous prefecture
//	 地级市: prefecture-level city
//	 盟: league
//
// 区县级:
//	 县: county
//	 自治县: autonomous county
//	 县级市: county-level city
//	 市辖区: district
//	 旗: banner
//	 自治旗: autonomous banner
//	 林区: forestry area
//	 特区: special district

type Loc struct {
	Id       int64
	Country  string // 国家
	Province string // 省/自治区/直辖市
	City     string // 地级市
	District string // 区/县
}

type Record struct {
	Start, End uint32
	Id         int64
}

func (loc *Loc) Split() (string, string, string, string) {
	if loc == nil {
		return "", "", "", ""
	}
	return loc.Country, loc.Province, loc.City, loc.District
}

func (loc *Loc) Short() string {
	if loc == nil {
		return ""
	}
	s := loc.Province + loc.City
	if loc.Country != "中国" {
		s = loc.Country + s
	}
	return s
}

func (loc *Loc) String() string {
	if loc == nil {
		return ""
	}
	return fmt.Sprintf("%d,%s,%s,%s,%s", loc.Id, loc.Country, loc.Province, loc.City, loc.District)
}

func New(recs []Record, locs []Loc) (*Finder, error) {
	recs1 := make([]Record, 0, len(recs))
	locs1 := make(map[int64]Loc, len(locs))
	for _, rec := range recs {
		recs1 = append(recs1, rec)
	}
	for _, loc := range locs {
		locs1[loc.Id] = loc
	}
	sort.Sort(orderByIp(recs1)) // 确认按照IP段进行排序
	return &Finder{
		recs: recs1,
		locs: locs1,
	}, nil
}

func Load(recReader, locReader io.Reader) (*Finder, error) {
	if recReader == nil {
		return nil, sderr.New("nil record reader")
	}
	if locReader == nil {
		return nil, sderr.New("nil location reader")
	}

	// records
	var recs []Record
	recReader1, err := sdcsv.NewReader(recReader, nil)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	var procErr error
	err = recReader1.EachRecord(func(_ int, rec []string) sdcsv.HandlerResult {
		var rec1 Record
		procErr = parseRecord(rec, &rec1)
		if procErr != nil {
			return sdcsv.Stop
		}
		recs = append(recs, rec1)
		return sdcsv.Continue
	})
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	if procErr != nil {
		return nil, sderr.WithStack(procErr)
	}

	// locations
	var locs []Loc
	locReader1, err := sdcsv.NewReader(locReader, nil)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	err = locReader1.EachRecord(func(_ int, rec []string) sdcsv.HandlerResult {
		var loc1 Loc
		procErr = parseLoc(rec, &loc1)
		if procErr != nil {
			return sdcsv.Stop
		}
		locs = append(locs, loc1)
		return sdcsv.Continue
	})
	if err != nil {
		return nil, sderr.WithStack(err)
	}

	// new
	return New(recs, locs)
}

func LoadFile(recFilename, locFilename string) (*Finder, error) {
	recReader, err := os.Open(recFilename)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() {
		_ = recReader.Close()
	}()
	locReader, err := os.Open(locFilename)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() {
		_ = locReader.Close()
	}()
	r, err := Load(recReader, locReader)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func (f *Finder) RecordNum() int {
	return len(f.recs)
}

func (f *Finder) LocNum() int {
	return len(f.locs)
}

func (f *Finder) Loc(ip string) *Loc {
	if ip == "" {
		return nil
	}
	ipInt, err := ip2int(ip)
	if err != nil {
		return nil
	}
	return f.LocInt(ipInt)
}

func (f *Finder) LocInt(ip uint32) *Loc {
	left, right := 0, len(f.recs)
	for {
		if left > right {
			break
		}
		middle := (left + right) / 2
		rec := &f.recs[middle]
		cr := compareIpAndRecord(ip, rec)
		if cr == 0 {
			loc, ok := f.locs[rec.Id]
			if !ok || loc.Id <= 0 {
				return nil
			} else {
				return &loc
			}
		} else if cr > 0 {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}
	return nil
}

func ip2int(ip string) (uint32, error) {
	ss := strings.Split(ip, ".")
	if len(ss) != 4 {
		return 0, sderr.New("parse ip error: " + ip)
	}
	ipBytes := make([]byte, 0, 4)
	for _, s := range ss {
		ipSeg, err := strconv.Atoi(s)
		if err != nil {
			return 0, sderr.New("parse ip error: " + ip)
		}
		if ipSeg < 0 || ipSeg > 255 {
			return 0, sderr.New("parse ip error: " + ip)
		}
		ipBytes = append(ipBytes, byte(ipSeg))
	}
	return binary.BigEndian.Uint32(ipBytes), nil
}

func int2ip(nn uint32) string {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func parseRecord(rec []string, to *Record) error {
	if len(rec) != 3 {
		return sderr.New("record size error")
	}
	start, err := ip2int(rec[0])
	if err != nil {
		return sderr.WithStack(err)
	}
	end, err := ip2int(rec[1])
	if err != nil {
		return sderr.WithStack(err)
	}
	locId, err := sdparse.Int64(rec[2])
	if err != nil {
		return sderr.WithStack(err)
	}
	if start > end {
		start, end = end, start
	}
	to.Start, to.End, to.Id = start, end, locId
	return nil
}

func parseLoc(rec []string, to *Loc) error {
	if len(rec) != 5 {
		return sderr.New("location size error")
	}
	locId, err := sdparse.Int64(rec[0])
	if err != nil {
		return sderr.WithStack(err)
	}
	to.Id = locId
	to.Country = rec[1]
	to.Province = rec[2]
	to.City = rec[3]
	to.District = rec[4]
	return nil
}

func compareRecord(rec1, rec2 *Record) int {
	if rec1 == rec2 {
		return 0
	}
	if rec1.Start <= rec1.End && rec2.Start <= rec2.End && rec1.End < rec2.Start {
		return -1
	}
	if rec1.Start <= rec1.End && rec2.Start <= rec2.End && rec2.End < rec1.Start {
		return 1
	}
	panic(sderr.Newf(
		"compare record error (%s->%s %s->%s)",
		int2ip(rec1.Start),
		int2ip(rec1.End),
		int2ip(rec2.Start),
		int2ip(rec2.End),
	))
}

func compareIpAndRecord(ip uint32, rec *Record) int {
	if ip >= rec.Start && ip <= rec.End {
		return 0
	} else if ip > rec.End {
		return 1
	} else if ip < rec.Start {
		return -1
	}
	panic(sderr.New("run here error"))
}

type orderByIp []Record

func (recs orderByIp) Len() int {
	return len(recs)
}
func (recs orderByIp) Swap(i, j int) {
	recs[i], recs[j] = recs[j], recs[i]
}
func (recs orderByIp) Less(i, j int) bool {
	return compareRecord(&recs[i], &recs[j]) < 0
}
