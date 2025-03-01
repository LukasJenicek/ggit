package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"syscall"
)

type Entry struct {
	ctime     uint32
	ctimeNsec uint32
	mtime     uint32
	mtimeNsec uint32
	// id of the hardware device the file is stored on
	dev   uint32
	inode uint32
	mode  uint32
	uid   uint32
	gid   uint32
	// size of file
	fileSize uint32
	// 20 bytes
	OID []byte
	// length of filename
	flags uint16
	Path  []byte
}

func (e *Entry) Content() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := []any{
		e.ctime,
		e.ctimeNsec,
		e.mtime,
		e.mtimeNsec,
		e.dev,
		e.inode,
		e.mode,
		e.uid,
		e.gid,
		e.fileSize,
		e.OID,
		e.flags,
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
		ctime:     binary.BigEndian.Uint32(data[0:4]),
		ctimeNsec: binary.BigEndian.Uint32(data[4:8]),
		mtime:     binary.BigEndian.Uint32(data[8:12]),
		mtimeNsec: binary.BigEndian.Uint32(data[12:16]),
		dev:       binary.BigEndian.Uint32(data[16:20]),
		inode:     binary.BigEndian.Uint32(data[20:24]),
		mode:      binary.BigEndian.Uint32(data[24:28]),
		uid:       binary.BigEndian.Uint32(data[28:32]),
		gid:       binary.BigEndian.Uint32(data[32:36]),
		fileSize:  binary.BigEndian.Uint32(data[36:40]),
		OID:       data[40:60],
		flags:     binary.BigEndian.Uint16(data[60:62]),
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
		ctime:     uint32(stat.Ctim.Sec),
		ctimeNsec: uint32(stat.Ctim.Nsec),
		mtime:     uint32(stat.Mtim.Sec),
		mtimeNsec: uint32(stat.Mtim.Nsec),
		dev:       uint32(stat.Dev),
		inode:     uint32(stat.Ino),
		mode:      uint32(mode),
		uid:       stat.Uid,
		gid:       stat.Gid,
		fileSize:  uint32(stat.Size),
		OID:       oid,
		flags:     uint16(flags),
		Path:      []byte(pathname),
	}, nil
}
