package pipeline

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"data-struct/pipeline/util"
)

type Flow struct {
}

type DataSet struct {
	shards    []*DataSetShard
	partition int
}
type DataSetShard struct {
	//write datas to task, task process it.
	pipeWriter []*util.Pipeline
	//read from pre dataset shard, already processed
	pipeReader []*util.Pipeline
}

//每一步的运行template
type Step struct {
	Function func(readers []io.Reader, writers []io.Writer) error
	Tasks    []Task
}

//task 用来连接流的上下游
type Task struct {
	step *Step
	//read from dataset shard
	pipeReaders []*util.Pipeline
	//write to next dataset shared
	pipWriter []*util.Pipeline
}
type Source interface {
	Generate(writer io.Writer) error
}

func NewDataSet(partition int) *DataSet {
	dataset := &DataSet{partition: partition}
	dataset.shards = make([]*DataSetShard, partition)
	return dataset
}

type TextSource struct {
	filePath string
	file     *os.File
	readAt   int
}

func NewTextSource(filePath string) (*TextSource, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	textSource := &TextSource{filePath: filePath, file: file}
	return textSource, nil
}
func (flow *Flow) Source(s Source) *DataSet {
	ret := NewDataSet(1)
	step := NewOneToOneStep(nil, ret)
	step.Function = func(readers []io.Reader, writers []io.Writer) error {
		return s.Generate(writers[0])
	}
	return nil
}
func (textSource *TextSource) Generate(writer io.Writer) error {
	buf := bufio.NewReader(textSource.file)
	buffer := bytes.Buffer{}
	for {
		readline, isPrefix, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			buffer.Write(readline)
			continue
		}
		buffer.Write(readline)
		_, err = buffer.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	return nil
}
func NewOneToOneStep(from *DataSet, to *DataSet) *Step {
	return nil
}
