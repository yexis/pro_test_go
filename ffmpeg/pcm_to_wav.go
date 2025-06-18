package ffmpeg

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func PcmToWav(pcmPath, wavPath string, sampleRate, channels, bitDepth int) error {
	// 打开PCM文件
	pcmFile, err := os.Open(pcmPath)
	if err != nil {
		return err
	}
	defer pcmFile.Close()

	// 获取PCM数据大小
	pcmStat, _ := pcmFile.Stat()
	dataLen := uint32(pcmStat.Size())

	// 计算一些WAV参数
	byteRate := uint32(sampleRate * channels * bitDepth / 8)
	blockAlign := uint16(channels * bitDepth / 8)

	// 打开WAV输出文件
	wavFile, err := os.Create(wavPath)
	if err != nil {
		return err
	}
	defer wavFile.Close()

	// 写入 WAV Header
	// ChunkID "RIFF"
	wavFile.Write([]byte("RIFF"))
	// ChunkSize = 36 + SubChunk2Size
	binary.Write(wavFile, binary.LittleEndian, uint32(36+dataLen))
	// Format "WAVE"
	wavFile.Write([]byte("WAVE"))

	// Subchunk1ID "fmt "
	wavFile.Write([]byte("fmt "))
	binary.Write(wavFile, binary.LittleEndian, uint32(16))         // Subchunk1Size
	binary.Write(wavFile, binary.LittleEndian, uint16(1))          // AudioFormat PCM = 1
	binary.Write(wavFile, binary.LittleEndian, uint16(channels))   // NumChannels
	binary.Write(wavFile, binary.LittleEndian, uint32(sampleRate)) // SampleRate
	binary.Write(wavFile, binary.LittleEndian, byteRate)           // ByteRate
	binary.Write(wavFile, binary.LittleEndian, blockAlign)         // BlockAlign
	binary.Write(wavFile, binary.LittleEndian, uint16(bitDepth))   // BitsPerSample

	// Subchunk2ID "data"
	wavFile.Write([]byte("data"))
	binary.Write(wavFile, binary.LittleEndian, dataLen)

	// 复制PCM数据
	_, err = io.Copy(wavFile, pcmFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := PcmToWav("input.pcm", "output.wav", 16000, 1, 16)
	if err != nil {
		fmt.Println("转换失败:", err)
	} else {
		fmt.Println("转换成功")
	}
}
