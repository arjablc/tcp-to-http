package response

import (
	"fmt"
	"strconv"
)

//NOTE: For future references

/*
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	const chunkSize int64 = 1024
	regNurseBytes := []byte("\r\n")
	chunkSizeHex := strconv.FormatInt(chunkSize, 16)
	// lets assume like 16 bytes chunk
	chunkSizeLine := fmt.Appendf(nil, "%s\r\n", chunkSizeHex)
	var writtenIdx int
	for i := 0; i < len(p); i += int(chunkSize) {
		n, err := w.Buf.Write(chunkSizeLine)
		if err != nil {
			return 0, fmt.Errorf("Error Writing chunk: %v", err)
		}
		writtenIdx += n
		end := min(i+int(chunkSize), len(p))
		data := append(p[i:end], regNurseBytes...)
		n, err = w.Buf.Write(data)
		if err != nil {
			return 0, fmt.Errorf("Error Writing chunk: %v", err)
		}
		writtenIdx += n
	}
	return writtenIdx, nil
}
*/

//NOTE: Actual Implementation

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)
	regNurseBytes := []byte("\r\n")
	chunkSizeHex := strconv.FormatInt(int64(chunkSize), 16)
	chunkSizeLine := fmt.Appendf(nil, "%s\r\n", chunkSizeHex)

	var writtenIdx int

	n, err := w.Buf.Write(chunkSizeLine)
	if err != nil {
		return writtenIdx, fmt.Errorf("Error Writing chunk: %v", err)
	}
	writtenIdx += n
	data := append(p, regNurseBytes...)
	n, err = w.Buf.Write(data)
	if err != nil {
		return writtenIdx, fmt.Errorf("Error Writing chunk: %v", err)
	}
	writtenIdx += n
	return writtenIdx, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	endBytes := fmt.Appendf(nil, "%s\r\n", strconv.FormatInt(0, 16))
	i, err := w.Buf.Write(endBytes)
	if err != nil {
		return 0, fmt.Errorf("Error Writing chunk Done : %v", err)
	}
	j, err := w.Buf.Write([]byte("\r\n"))
	if err != nil {
		return 0, fmt.Errorf("Error Writing chunk Done : %v", err)
	}
	return i + j, nil
}
