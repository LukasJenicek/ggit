package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"syscall"
)

type Entries map[string]*Entry

func (e *Entries) SortedValues() []*Entry {
	entries := make([]*Entry, 0, len(*e))

	for _, v := range *e {
		entries = append(entries, v)
	}

	sort.Slice(entries, func(i, j int) bool {
		return string(entries[i].Path) < string(entries[j].Path)
	})

	return entries
}

type Entry struct {
	Ctime     uint32
	CtimeNsec uint32
	Mtime     uint32
	MtimeNsec uint32
	// id of the hardware device the file is stored on
	Dev   uint32
	Inode uint32
	Mode  uint32
	Uid   uint32
	Gid   uint32
	// size of file
	FileSize uint32
	// 20 bytes
	OID []byte
	// length of filename
	Flags uint16
	Path  []byte
}

func (e *Entry) Content() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := []any{
		e.Ctime,
		e.CtimeNsec,
		e.Mtime,
		e.MtimeNsec,
		e.Dev,
		e.Inode,
		e.Mode,
		e.Uid,
		e.Gid,
		e.FileSize,
		e.OID,
		e.Flags,
		e.Path,
		[]byte{0x00},
	}

	for _, d := range data {
		// git uses BigEndian for consistency across architectures
		if err := binary.Write(buf, binary.BigEndian, d); err != nil {
			return nil, fmt.Errorf("encoding entry: %w", err)
		}
	}

	// 1-8 nul bytes as necessary to pad the entry to a multiple of eight bytes
	// while keeping the name NUL-terminated.

	mod := buf.Len() % 8
	if mod > 0 {
		pad := math.Abs(float64(buf.Len()%8 - 8))
		if err := binary.Write(buf, binary.BigEndian, make([]byte, int(pad))); err != nil {
			return nil, fmt.Errorf("encoding entry: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func NewEntryFromBytes(data []byte, pathLen int) (*Entry, error) {
	length := 62 + pathLen
	if len(data) < length {
		return nil, fmt.Errorf("entry too short (%d < %d)", len(data), length)
	}

	return &Entry{
		Ctime:     binary.BigEndian.Uint32(data[0:4]),
		CtimeNsec: binary.BigEndian.Uint32(data[4:8]),
		Mtime:     binary.BigEndian.Uint32(data[8:12]),
		MtimeNsec: binary.BigEndian.Uint32(data[12:16]),
		Dev:       binary.BigEndian.Uint32(data[16:20]),
		Inode:     binary.BigEndian.Uint32(data[20:24]),
		Mode:      binary.BigEndian.Uint32(data[24:28]),
		Uid:       binary.BigEndian.Uint32(data[28:32]),
		Gid:       binary.BigEndian.Uint32(data[32:36]),
		FileSize:  binary.BigEndian.Uint32(data[36:40]),
		OID:       data[40:60],
		Flags:     binary.BigEndian.Uint16(data[60:62]),
		Path:      data[62 : 62+pathLen],
	}, nil
}

// NewEntry
// conversion from int64 to int32 is intentional. git is using int32 to support old architecture
//
//nolint:gosec
func NewEntry(pathname string, fInfo os.FileInfo, oid []byte) (*Entry, error) {
	// Get the underlying data source and type assert to syscall.Stat_t
	stat, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("not a syscall.Stat_t type")
	}

	mode := regularMode
	if fInfo.Mode().Perm()&0o100 != 0 {
		mode = executableMode
	}

	flags := len(pathname)
	if flags > maxPathSize {
		flags = maxPathSize
	}

	return &Entry{
		Ctime:     uint32(stat.Ctim.Sec),
		CtimeNsec: uint32(stat.Ctim.Nsec),
		Mtime:     uint32(stat.Mtim.Sec),
		MtimeNsec: uint32(stat.Mtim.Nsec),
		Dev:       uint32(stat.Dev),
		Inode:     uint32(stat.Ino),
		Mode:      uint32(mode),
		Uid:       stat.Uid,
		Gid:       stat.Gid,
		FileSize:  uint32(stat.Size),
		OID:       oid,
		Flags:     uint16(flags),
		Path:      []byte(pathname),
	}, nil
}
