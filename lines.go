package waveform

type Line struct {
	Count   int
	SumMin  int
	SumMax  int
	PeakMin int
	PeakMax int
}

type Lines []Line

func (l *Line) add(s *Sample) {
	l.Count++
	l.SumMin += int(s.Minimum)
	l.SumMax += int(s.Maximum)

	if int(s.Minimum) < l.PeakMin {
		l.PeakMin = int(s.Minimum)
	}
	if int(s.Maximum) > l.PeakMax {
		l.PeakMax = int(s.Maximum)
	}
}

func GenerateLines(wf *Waveform, linesCount int) Lines {
	samplesPerLine := float64(wf.Length) / float64(linesCount)

	lines := make([]Line, linesCount)
	currentLine := 0
	sampleCount := 0

	// Calculate lines
	for sample := range wf.Samples() {
		sampleCount++

		lines[currentLine].add(&sample)

		if int(samplesPerLine*float64(currentLine+1)) < sampleCount {
			currentLine++
		}
	}

	return lines
}

func (lines Lines) Max() (result float32) {
	for i := range lines {
		l := &lines[i]
		min := float32(-l.SumMin) / float32(l.Count)
		max := float32(l.SumMax) / float32(l.Count)

		if result < min {
			result = min
		}
		if result < max {
			result = max
		}
	}

	return
}
