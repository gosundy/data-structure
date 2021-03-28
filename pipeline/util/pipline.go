package util

import "io"

type Pipeline struct {
	writer *io.PipeWriter
	reader *io.PipeReader
}

func NewPipeLine() (pipeline *Pipeline) {
	pipeline = &Pipeline{}
	reader, writer := io.Pipe()
	pipeline.writer = writer
	pipeline.reader = reader
	return pipeline
}
