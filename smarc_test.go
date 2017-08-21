package smarc

import (
	"log"
	"math"
	"os"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func toFloatBuffer(buf *audio.IntBuffer, bitDepth float64) *audio.FloatBuffer {
	newB := &audio.FloatBuffer{}
	newB.Data = make([]float64, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = float64(buf.Data[i]) / math.Pow(2, bitDepth)
	}
	newB.Format = &audio.Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

func toIntBuffer(buf *audio.FloatBuffer, bitDepth float64) *audio.IntBuffer {
	newB := &audio.IntBuffer{}
	newB.Data = make([]int, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = int(buf.Data[i] * math.Pow(2, bitDepth))
	}
	newB.Format = &audio.Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}
func TestResample(t *testing.T) {
	f, err := os.Open("test/drums.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := wav.NewDecoder(f)
	buf, err := w.FullPCMBuffer()
	if err != nil {
		t.Fatal(err)
	}
	outRate := 16000
	buff := toFloatBuffer(buf, float64(w.BitDepth))
	buff.Data = Resample(buff.Data, int(w.SampleRate), outRate, 0.97, 0.05, 150.0, 0.000001)
	outFile, err := os.Create("test/outfile.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()
	wr := wav.NewEncoder(outFile, outRate, int(w.BitDepth), int(w.NumChans), int(w.WavAudioFormat))
	wr.Write(toIntBuffer(buff, float64(w.BitDepth)))
	if err = wr.Close(); err != nil {
		log.Fatal(err)
	}

}
