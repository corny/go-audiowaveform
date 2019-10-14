package waveform

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Waveform struct {
	Flags           uint32
	SampleRate      int32
	SamplesPerPixel int32
	Length          uint32
	Channels        int32

	reader io.Reader
}

type Sample struct {
	Minimum int8
	Maximum int8
}

const (
	Flag16bit = 0
	Flag8bit  = 1
)

func ReadWaveform(r io.Reader) (wf Waveform, err error) {
	var version int32
	binary.Read(r, binary.LittleEndian, &version)
	binary.Read(r, binary.LittleEndian, &wf.Flags)
	binary.Read(r, binary.LittleEndian, &wf.SampleRate)
	binary.Read(r, binary.LittleEndian, &wf.SamplesPerPixel)
	binary.Read(r, binary.LittleEndian, &wf.Length)

	switch version {
	case 1:

	default:
		err = fmt.Errorf("unsupported version: %d", version)
		return
	}

	wf.reader = r
	return
}

func (wf *Waveform) Samples() <-chan Sample {
	out := make(chan Sample)
	go func() {
		defer close(out)
		for {
			s := Sample{}

			if err := binary.Read(wf.reader, binary.LittleEndian, &s); err != nil {
				break
			}

			out <- s
		}
	}()
	return out
}

func (wf *Waveform) EachLine(lineCount int, callback func(min, max float32))  {

	lines := GenerateLines(wf, lineCount)
	scaleFactor := float32(1) / lines.Max()

	// Calculate and run callbacks
	for i := range lines {
		l := &lines[i]
		callback(
			scaleFactor*float32(l.SumMin)/float32(l.Count),
			scaleFactor*float32(l.SumMax)/float32(l.Count),
		)
	}
}