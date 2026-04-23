package wal

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type LogRecord struct {
	Length uint32
	LSN    uint64
	Type   uint8
	Value  int64
}

const (
	RecordTypeINCR  = 1
	RecordTypeRESET = 2
	recordHeaderLen = 4 + 8 + 1 + 8 // Length + LSN + Type + Value
)

type WAL struct {
	f   *os.File
	lsn uint64
}

func OpenWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	return &WAL{f: f, lsn: uint64(stat.Size())}, nil
}

func (w *WAL) WriteRecord(rec *LogRecord) error {
	rec.LSN = w.lsn
	rec.Length = recordHeaderLen
	buf := make([]byte, recordHeaderLen)
	binary.LittleEndian.PutUint32(buf[0:4], rec.Length)
	binary.LittleEndian.PutUint64(buf[4:12], rec.LSN)
	buf[12] = rec.Type
	binary.LittleEndian.PutUint64(buf[13:21], uint64(rec.Value))
	n, err := w.f.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return errors.New("short write")
	}
	err = w.f.Sync()
	if err != nil {
		return err
	}
	w.lsn += uint64(len(buf))
	return nil
}

func (w *WAL) Replay(apply func(LogRecord)) error {
	_, err := w.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	for {
		hdr := make([]byte, recordHeaderLen)
		n, err := io.ReadFull(w.f, hdr)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
		if n != recordHeaderLen {
			break
		}
		rec := LogRecord{
			Length: binary.LittleEndian.Uint32(hdr[0:4]),
			LSN:    binary.LittleEndian.Uint64(hdr[4:12]),
			Type:   hdr[12],
			Value:  int64(binary.LittleEndian.Uint64(hdr[13:21])),
		}
		apply(rec)
	}
	return nil
}

func (w *WAL) Close() error {
	return w.f.Close()
}
