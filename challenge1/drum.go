// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func parse(data []byte, p *Pattern) error {
	dataReader := bytes.NewReader(data)

	fileSize := dataReader.Len()
	if !isValidFile(dataReader) {
		return errors.New("Not a valid .splice file.")
	}

	encodedDataSize, fileSizeErr := getEncodedDataSize(dataReader)
	if fileSizeErr != nil || encodedDataSize > dataReader.Len() {
		return errors.New("Encoded data size is incorrect.")
	}

	hwVersion, hwVersionErr := readHardwareVersion(dataReader)
	if hwVersionErr != nil {
		return hwVersionErr
	}
	p.Version = hwVersion

	tempo, readTempErr := readTempo(dataReader)
	if readTempErr != nil {
		return readTempErr
	}
	p.Tempo = tempo

	tracks, readTracksErr := readTracks(dataReader, fileSize, encodedDataSize)
	if readTracksErr != nil {
		return readTracksErr
	}
	p.Tracks = tracks

	return nil
}

func readTracks(reader *bytes.Reader, fileSize, encodedDataSize int) ([]Track, error) {
	var tracks []Track

	position := (fileSize - encodedDataSize) + (encodedDataSize - reader.Len())

	for position < encodedDataSize {
		var id int32
		binary.Read(reader, binary.LittleEndian, &id)

		channelNameSize, _ := reader.ReadByte()
		channelBytes := make([]byte, channelNameSize)
		_, err := reader.Read(channelBytes)
		if err != nil {
			return []Track{}, errors.New("Could not read Track name with id " + id)
		}

		pattern := make([]uint32, 4)
		patternReadErr := binary.Read(reader, binary.LittleEndian, &pattern)
		if patternReadErr != nil {
			return []Track{}, errors.New("Could not read Track step with id " + id)
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

func getEncodedDataSize(reader *bytes.Reader) (int, error) {
	var encodedDataSize int64
	fileSizeReadError := binary.Read(reader, binary.BigEndian, &encodedDataSize)

	if fileSizeReadError != nil {
		return 0, fileSizeReadError
	}

	// int cast because it will be easier to deal
	// with throughout the code, and I feel like using
	// an int64 for the file size is a bit much
	// and we can get away with the precision loss
	return int(encodedDataSize), nil
}
