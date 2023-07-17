package singgeo

import (
	"bytes"
	"io"
	"os"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/rw"
)

type GeoSiteReader struct {
	reader       io.ReadSeeker
	domainIndex  map[string]int
	domainLength map[string]int
}

func LoadSite(data []byte) (reader *GeoSiteReader, codes []string, err error) {
	reader = &GeoSiteReader{
		reader: bytes.NewReader(data),
	}
	codes, err = loadSite(reader)
	return
}

func LoadSiteFromFile(path string) (*GeoSiteReader, []string, error) {
	content, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	reader := &GeoSiteReader{
		reader: content,
	}
	codes, err := loadSite(reader)
	return reader, codes, err
}

func loadSite(reader *GeoSiteReader) (codes []string, err error) {
	err = reader.readMetadata()
	if err != nil {
		common.Close(reader.reader)
		return nil, err
	}
	codes = make([]string, 0, len(reader.domainIndex))
	for code := range reader.domainIndex {
		codes = append(codes, code)
	}
	return codes, nil
}

func (r *GeoSiteReader) readMetadata() error {
	version, err := rw.ReadByte(r.reader)
	if err != nil {
		return err
	}
	if version != 0 {
		return E.New("unknown version")
	}
	entryLength, err := rw.ReadUVariant(r.reader)
	if err != nil {
		return err
	}
	keys := make([]string, entryLength)
	domainIndex := make(map[string]int)
	domainLength := make(map[string]int)
	for i := 0; i < int(entryLength); i++ {
		var (
			code       string
			codeIndex  uint64
			codeLength uint64
		)
		code, err = rw.ReadVString(r.reader)
		if err != nil {
			return err
		}
		keys[i] = code
		codeIndex, err = rw.ReadUVariant(r.reader)
		if err != nil {
			return err
		}
		codeLength, err = rw.ReadUVariant(r.reader)
		if err != nil {
			return err
		}
		domainIndex[code] = int(codeIndex)
		domainLength[code] = int(codeLength)
	}
	r.domainIndex = domainIndex
	r.domainLength = domainLength
	return nil
}

func (r *GeoSiteReader) Read(code string) ([]Item, error) {
	index, exists := r.domainIndex[code]
	if !exists {
		return nil, E.New("code ", code, " not exists!")
	}
	_, err := r.reader.Seek(int64(index), io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	counter := rw.ReadCounter{Reader: r.reader}
	domain := make([]Item, r.domainLength[code])
	for i := range domain {
		var (
			item Item
			err  error
		)
		item.Type, err = rw.ReadByte(&counter)
		if err != nil {
			return nil, err
		}
		item.Value, err = rw.ReadVString(&counter)
		if err != nil {
			return nil, err
		}
		domain[i] = item
	}
	_, err = r.reader.Seek(int64(-index)-counter.Count(), io.SeekCurrent)
	return domain, err
}
