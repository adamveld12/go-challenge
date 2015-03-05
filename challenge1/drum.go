// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func parse(data []byte, p *Pattern) error {
	buf := bytes.NewReader(data)

	totalSize := buf.Len()
	if !isValidFile(buf) {
		return errors.New("Not a valid .splice file.")
	}

	fileSize, fileSizeErr := getSize(buf)
	if fileSizeErr != nil || fileSize > buf.Len() {
		return errors.New("File size is not correct.")
	}

	hwVersion, hwVersionErr := readHardwareVersion(buf)
	if hwVersionErr != nil {
		return hwVersionErr
	}
	p.Version = hwVersion

	tempo, readTempErr := readTempo(buf)
	if readTempErr != nil {
		return readTempErr
	}
	p.Tempo = tempo

	tracks, readTracksErr := readTracks(buf, totalSize, fileSize)
	if readTracksErr != nil {
		return readTracksErr
	}
	p.Tracks = tracks

	return nil
}

func readTracks(reader *bytes.Reader, totalSize, fileSize int) ([]Track, error) {
	var tracks []Track

	endPosition := fileSize
	position := (totalSize - fileSize) + (fileSize - reader.Len())

	for position < endPosition {
		var id int32
		binary.Read(reader, binary.LittleEndian, &id)

		channelNameSize, _ := reader.ReadByte()
		channelBytes := make([]byte, channelNameSize)
		_, err := reader.Read(channelBytes)
		if err != nil {
			return tracks, errors.New("Could not read Track name")
		}

		pattern := make([]uint32, 4)
		patternReadErr := binary.Read(reader, binary.LittleEndian, &pattern)
		if patternReadErr != nil {
			return tracks, errors.New("Could not read Track steps")
		}

		tracks = append(tracks, Track{
			id,
			string(channelBytes),
			pattern})

		position += int(21) + int(channelNameSize)
	}

	return tracks, nil
}

func readTempo(reader *bytes.Reader) (float32, error) {
	var tempo float32
	err := binary.Read(reader, binary.LittleEndian, &tempo)
	if err != nil {
		return 0.0, errors.New("Could not read tempo")
	}
	return tempo, nil
}

func readHardwareVersion(reader *bytes.Reader) (string, error) {
	versionBytes := make([]byte, 32)
	_, versionReadError := reader.Read(versionBytes)
	if versionReadError != nil {
		return "", versionReadError
	}

	versionString := string(bytes.Trim(versionBytes, "\x00"))
	if versionString == "" {
		return "", errors.New("The file version is incorrect.")
	}

	return versionString, nil
}

func isValidFile(reader *bytes.Reader) bool {
	spliceBytes := make([]byte, 6)
	_, notSpliceFileErr := reader.Read(spliceBytes)

	spliceHeader := string(spliceBytes)
	if notSpliceFileErr != nil || spliceHeader != "SPLICE" {
		return false
	}

	return true
}

func getSize(reader *bytes.Reader) (int, error) {
	var fileSize int64
	fileSizeReadError := binary.Read(reader, binary.BigEndian, &fileSize)

	if fileSizeReadError != nil {
		return 0, fileSizeReadError
	}

	return int(fileSize), nil
}
