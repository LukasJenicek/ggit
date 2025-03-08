package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"syscall"
)

type Parents struct {
	mu      sync.Mutex
	parents map[string][]*Entry
}

func NewParents() *Parents {
	return &Parents{
		parents: make(map[string][]*Entry),
	}
}

func (e *Parents) Add(key string, value *Entry) {
	e.mu.Lock()
	defer e.mu.Unlock()

	values, ok := e.parents[key]
	if !ok {
		values = make([]*Entry, 0)
	}

	values = append(values, value)

	e.parents[key] = values
}

func (e *Parents) Get(key string) ([]*Entry, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	parents, ok := e.parents[key]

	return parents, ok
}

func (e *Parents) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.parents, key)
}

type Entries struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

func NewEntries() *Entries {
	return &Entries{
		entries: make(map[string]*Entry),
	}
}

func (e *Entries) Add(key string, value *Entry) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.entries[key] = value
}

func (e *Entries) Get(key string) (*Entry, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry, ok := e.entries[key]

	return entry, ok
}

func (e *Entries) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.entries, key)
}

func (e *Entries) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return len(e.entries)
}

func (e *Entries) SortedValues() []*Entry {
	e.mu.Lock()
	defer e.mu.Unlock()

	entries := make([]*Entry, 0, len(e.entries))

	for _, v := range e.entries {
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
	UID   uint32
	GID   uint32
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
		e.UID,
		e.GID,
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
		UID:       binary.BigEndian.Uint32(data[28:32]),
		GID:       binary.BigEndian.Uint32(data[32:36]),
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
		UID:       stat.Uid,
		GID:       stat.Gid,
		FileSize:  uint32(stat.Size),
		OID:       oid,
		Flags:     uint16(flags),
		Path:      []byte(pathname),
	}, nil
}
