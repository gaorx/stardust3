package sdcsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/gaorx/stardust3/sdcoll/sdstrcoll"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
)

type Reader struct {
	reader *csv.Reader
	fields []string
}

type Options struct {
	Header           bool
	Fields           []string
	Comma            rune
	Comment          rune
	FieldsPerRecord  int
	LazyQuotes       bool
	TrimLeadingSpace bool
	ReuseRecord      bool
}

type HandlerResult int

const (
	Stop     HandlerResult = 0
	Continue HandlerResult = 1
)

func NewReader(r io.Reader, opts *Options) (*Reader, error) {
	if r == nil {
		return nil, sderr.New("nil reader")
	}
	r1 := csv.NewReader(r)
	var fields []string
	if opts != nil {
		if opts.Comma != 0 {
			r1.Comma = opts.Comma
		}
		if opts.Comma != 0 {
			r1.Comment = opts.Comment
		}
		if opts.FieldsPerRecord > 0 {
			r1.FieldsPerRecord = opts.FieldsPerRecord
		}
		r1.LazyQuotes = opts.LazyQuotes
		r1.TrimLeadingSpace = opts.TrimLeadingSpace
		r1.ReuseRecord = opts.ReuseRecord
	}
	if opts != nil && len(opts.Fields) > 0 {
		fields = sdstrcoll.Slice(opts.Fields).Copy()
	}
	if opts != nil && opts.Header {
		header, err := r1.Read()
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		header = sdstrcoll.Slice(header).Copy()
		if len(fields) == 0 {
			fields = header
		}
	}
	return &Reader{
		reader: r1,
		fields: fields,
	}, nil
}

func NewReaderData(b []byte, opts *Options) (*Reader, error) {
	if b == nil {
		return nil, sderr.New("nil data")
	}
	return NewReader(bytes.NewReader(b), opts)
}

func NewReaderText(s string, opts *Options) (*Reader, error) {
	return NewReader(strings.NewReader(s), opts)
}

func NewReaderFile(filename string, opts *Options) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return NewReader(f, opts)
}

func (r *Reader) Fields() []string {
	return r.fields
}

func (r *Reader) Read() ([]string, error) {
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return rec, nil
}

func (r *Reader) ReadAll() ([][]string, error) {
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return recs, nil
}

func (r *Reader) ReadMap() (map[string]string, error) {
	if len(r.fields) == 0 {
		return nil, sderr.New("no field")
	}
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return makeMap(r.fields, rec), nil
}

func (r *Reader) ReadMapAll() ([]map[string]string, error) {
	if len(r.fields) == 0 {
		return nil, sderr.New("no field")
	}
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	maps := make([]map[string]string, 0, len(recs))
	for _, rec := range recs {
		maps = append(maps, makeMap(r.fields, rec))
	}
	return maps, nil
}

func (r *Reader) ReadObject() (sdjson.Object, error) {
	if len(r.fields) == 0 {
		return nil, sderr.New("no field")
	}
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return makeObject(r.fields, rec), nil
}

func (r *Reader) ReadObjectAll() ([]sdjson.Object, error) {
	if len(r.fields) == 0 {
		return nil, sderr.New("no field")
	}
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	objs := make([]sdjson.Object, 0, len(recs))
	for _, rec := range recs {
		objs = append(objs, makeObject(r.fields, rec))
	}
	return objs, nil
}

func (r *Reader) EachRecord(h func(recNo int, rec []string) HandlerResult) error {
	if h == nil {
		return sderr.New("nil handler")
	}
	recNo := 0
	for {
		rec, err := r.Read()
		if err != nil {
			if sderr.Cause(err) == io.EOF {
				break
			} else {
				return sderr.WithStack(err)
			}
		}
		hr := h(recNo, rec)
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func (r *Reader) EachMap(h func(recNo int, rec map[string]string) HandlerResult) error {
	if h == nil {
		return sderr.New("nil handler")
	}
	if len(r.fields) == 0 {
		return sderr.New("no field")
	}
	recNo := 0
	for {
		rec, err := r.Read()
		if err != nil {
			if sderr.Cause(err) == io.EOF {
				break
			} else {
				return sderr.WithStack(err)
			}
		}
		hr := h(recNo, makeMap(r.fields, rec))
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func (r *Reader) EachObject(h func(recNo int, rec sdjson.Object) HandlerResult) error {
	if h == nil {
		return sderr.New("nil handler")
	}
	if len(r.fields) == 0 {
		return sderr.New("no field")
	}
	recNo := 0
	for {
		rec, err := r.Read()
		if err != nil {
			if sderr.Cause(err) == io.EOF {
				break
			} else {
				return sderr.WithStack(err)
			}
		}
		hr := h(recNo, makeObject(r.fields, rec))
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func makeMap(fields, record []string) map[string]string {
	fieldNum, valNum := len(fields), len(record)
	m := make(map[string]string, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := fields[i]
		v := ""
		if i < valNum {
			v = record[i]
		}
		m[field] = v
	}
	return m
}

func makeObject(fields, record []string) sdjson.Object {
	fieldNum, valNum := len(fields), len(record)
	o := make(sdjson.Object, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := fields[i]
		v := ""
		if i < valNum {
			v = record[i]
		}
		o[field] = v
	}
	return o
}
