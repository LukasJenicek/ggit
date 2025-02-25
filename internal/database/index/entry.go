package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
)

type Entry struct {
	ctime     int32
	ctimeNsec int32
	mtime     int32
	mtimeNsec int32
	// id of the hardware device the file is stored on
	dev   int32
	inode int32
	mode  int32
	uid   uint32
	gid   uint32
	// length of relative file path
	fileSize int32
	oid      []byte
	flags    int16
	path     []byte
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
		e.oid,
		e.flags,
		e.path,
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
	pad := buf.Len() % 8
	if err := binary.Write(buf, binary.BigEndian, make([]byte, pad)); err != nil {
		return nil, fmt.Errorf("encoding entry: %w", err)
	}

	return buf.Bytes(), nil
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
		ctime:     int32(stat.Ctim.Sec),
		ctimeNsec: int32(stat.Ctim.Nsec),
		mtime:     int32(stat.Mtim.Sec),
		mtimeNsec: int32(stat.Mtim.Nsec),
		dev:       int32(stat.Dev),
		inode:     int32(stat.Ino),
		mode:      int32(mode),
		uid:       stat.Uid,
		gid:       stat.Gid,
		fileSize:  int32(stat.Size),
		oid:       oid,
		flags:     int16(flags),
		path:      []byte(pathname),
	}, nil
}
