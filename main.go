package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const FMT_ID_LINER_PCM = 1

type WavHdr struct {
	Size           uint32
	FmtSize        uint32
	FmtId          uint16
	ChanelNum      uint16
	SamplingRateHz uint32
	ByteRate       uint32
	BlkSize        uint16
	BitPerSample   uint16
	ExtSize        uint16
	ExtBin         []uint8

	WavesSize uint32
	WavesBin  []uint8
}

func main() {
	filepath := "sample.wav"
	fi, err := os.Stat(filepath)
	if err != nil {
		os.Exit(-1)
	}

	f, err := os.Open(filepath)
	if err != nil {
		os.Exit(-1)
	}
	defer f.Close()

	bin := make([]byte, fi.Size())
	f.Read(bin)

	wavHdr := WavHdr{}
	pos := 0
	hRiff := string(bin[pos:4])
	if hRiff != "RIFF" {
		panic("Not RIFF")
	}
	pos += 4

	wavHdr.Size, err = FromLeToUInt32(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 4
	fmt.Println(wavHdr.Size)

	hWave := string(bin[pos : pos+4])
	if hWave != "WAVE" {
		panic("Not WAVE")
	}
	pos += 4

	hFmt := string(bin[pos : pos+4])
	if hFmt != "fmt " {
		panic("Not fmt ")
	}
	pos += 4

	wavHdr.FmtSize, err = FromLeToUInt32(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 4

	wavHdr.FmtId, err = FromLeToUInt16(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 2

	wavHdr.ChanelNum, err = FromLeToUInt16(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 2

	wavHdr.SamplingRateHz, err = FromLeToUInt32(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 4

	wavHdr.ByteRate, err = FromLeToUInt32(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 4

	wavHdr.BlkSize, err = FromLeToUInt16(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 2

	wavHdr.BitPerSample, err = FromLeToUInt16(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 2

	fmt.Println(wavHdr.FmtSize)
	fmt.Println(wavHdr.FmtId)
	fmt.Println(wavHdr.ChanelNum)
	fmt.Println(wavHdr.SamplingRateHz)
	fmt.Println(wavHdr.ByteRate)
	fmt.Println(wavHdr.BlkSize)
	fmt.Println(wavHdr.BitPerSample)

	if wavHdr.FmtId == FMT_ID_LINER_PCM {
		// no ext header
		fmt.Println("liner PCM format")
	} else {
		fmt.Println("not liner PCM format")
		wavHdr.ExtSize, err = FromLeToUInt16(bin[pos:])
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		wavHdr.ExtBin = bin[pos : pos+int(wavHdr.ExtSize)]
		pos += int(wavHdr.ExtSize)
	}
	hData := string(bin[pos : pos+4])
	if hData != "data" {
		panic("Not fmt ")
	}
	pos += 4

	wavHdr.WavesSize, err = FromLeToUInt32(bin[pos:])
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pos += 4

	wavHdr.WavesBin = bin[pos : pos+int(wavHdr.WavesSize)]
	pos += int(wavHdr.WavesSize)
	fmt.Println(wavHdr.WavesSize)
	fmt.Println(pos)

}

func FromLeToUInt16(buf []byte) (uint16, error) {
	var val uint16
	if len(buf) < 2 {
		return val, errors.New("buf is too short")
	}

	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func FromLeToUInt32(buf []byte) (uint32, error) {
	var val uint32
	if len(buf) < 4 {
		return val, errors.New("buf is too short")
	}

	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}
