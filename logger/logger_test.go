package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"

	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleLogEntry struct {
	Code    int    `json:"code"`
	Label   string `json:"label"`
	Service string `json:"service"`
	Level   string `json:"level"`
}

type simpleLogEntryWithOneStackTrace struct {
	Code       int    `json:"code"`
	Label      string `json:"label"`
	Service    string `json:"service"`
	Level      string `json:"level"`
	StackTrace string `json:"ErrorStack_0"`
}

type logEntryWithTextField struct {
	Code      int    `json:"code"`
	Label     string `json:"label"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	FieldName string `json:"Field name"`
}

type logEntryWithStructField struct {
	Code    int    `json:"code"`
	Label   string `json:"label"`
	Service string `json:"service"`
	Level   string `json:"level"`
	Field   Field  `json:"field"`
}

type Field struct {
	Value string `json:"value"`
}

func TestParseLevel(t *testing.T) {
	vals := map[string]Level{
		"debug": DEBUG,
		"DEBUG": DEBUG,
		"Debug": DEBUG,
		"error": ERROR,
		"info":  INFO,
	}

	for str := range vals {
		lvl, err := ParseLevel(str)

		assert.Nil(t, err)
		assert.EqualValues(t, *lvl, vals[str])
	}

	_, err := ParseLevel("notalevel")

	assert.NotNil(t, err)
}

func Test_StackTrace_Write_Debug(t *testing.T) {

	var b bytes.Buffer
	buffer := &b
	writer := bufio.NewWriter(buffer)

	Init(DEBUG, writer, false)

	NewLogEntryFactory("test").MakeEntry(DEBUG, 1, "test")().Write(errors.Wrap(errors.New("inner"), "outer"))

	writer.Flush()

	actual := simpleLogEntryWithOneStackTrace{}
	err := json.Unmarshal(b.Bytes(), &actual)
	assert.NoError(t, err)

	if len(actual.StackTrace) == 0 {
		t.Fail()
	}
}

func Test_StackTrace_NoWrite_Debug(t *testing.T) {
	var b bytes.Buffer
	buffer := &b
	writer := bufio.NewWriter(buffer)

	Init(DEBUG, writer, false)
	SetStackTraceLevels(ERROR)

	NewLogEntryFactory("test").MakeEntry(DEBUG, 1, "test")().Write(errors.Wrap(errors.New("inner"), "outer"))

	writer.Flush()

	actual := simpleLogEntryWithOneStackTrace{}
	err := json.Unmarshal(b.Bytes(), &actual)
	assert.NoError(t, err)

	if len(actual.StackTrace) != 0 {
		t.Fail()
	}
}

func Test_Simple_Write(t *testing.T) {

	var b bytes.Buffer
	buffer := &b
	writer := bufio.NewWriter(buffer)

	Init(DEBUG, writer, false)

	NewLogEntryFactory("test").MakeEntry(DEBUG, 1, "test")().Write()

	writer.Flush()

	expected := simpleLogEntry{
		Label:   "test",
		Service: "test",
		Level:   "debug",
		Code:    1,
	}

	actual := simpleLogEntry{}
	err := json.Unmarshal(b.Bytes(), &actual)
	assert.NoError(t, err)

	if expected != actual {
		t.Fail()
	}
}

func Test_WithTextField_Write(t *testing.T) {

	var b bytes.Buffer
	buffer := &b
	writer := bufio.NewWriter(buffer)

	Init(DEBUG, writer, false)

	NewLogEntryFactory("test").
		MakeEntry(DEBUG, 1, "test")().
		WithField("Field name", "Field val").
		Write()

	writer.Flush()

	expected := logEntryWithTextField{
		Label:     "test",
		Service:   "test",
		Level:     "debug",
		Code:      1,
		FieldName: "Field val",
	}

	actual := logEntryWithTextField{}
	err := json.Unmarshal(b.Bytes(), &actual)
	assert.NoError(t, err)

	if expected != actual {
		t.Fail()
	}
}

func Test_WithStructField_Write(t *testing.T) {

	var b bytes.Buffer
	buffer := &b
	writer := bufio.NewWriter(buffer)

	Init(DEBUG, writer, false)

	NewLogEntryFactory("test").
		MakeEntry(DEBUG, 1, "test")().
		WithField("Field", Field{Value: "Field val"}).
		Write()

	writer.Flush()

	expected := logEntryWithStructField{
		Label:   "test",
		Service: "test",
		Level:   "debug",
		Code:    1,
		Field:   Field{Value: "Field val"},
	}

	actual := logEntryWithStructField{}
	err := json.Unmarshal(b.Bytes(), &actual)
	assert.NoError(t, err)

	if expected != actual {
		t.Fail()
	}
}
