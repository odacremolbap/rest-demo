package logrus

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"

	lgrus "github.com/Sirupsen/logrus"
)

const (
	dateTimeSimpleFormat = "2006/01/02 15:04:05"
)

var (
	emptyFieldMap lgrus.FieldMap
	orderedFields []string
)

func init() {
	orderedFields = []string{"ts", "code", "method", "uri", "remote", "bytes", "elapsed"}
}

// ServerFormatter formats logs into text
type ServerFormatter struct {

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool
	// DateTimeFormat for log ts field
	DateTimeFormat string

	sync.Once
}

// NewServerFomatter creates a server logrus formatter
func NewServerFomatter(quoteEmptyFields bool, dateTimeFormat string) *ServerFormatter {
	if len(dateTimeFormat) == 0 {
		dateTimeFormat = dateTimeSimpleFormat
	}
	return &ServerFormatter{
		QuoteEmptyFields: quoteEmptyFields,
		DateTimeFormat:   dateTimeFormat,
	}
}

// Format renders a single log entry
func (f *ServerFormatter) Format(entry *lgrus.Entry) ([]byte, error) {

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.printLog(b, entry, keys)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *ServerFormatter) printLog(b *bytes.Buffer, entry *lgrus.Entry, keys []string) {

	f.appendKeyValue(b, "ts", entry.Time.UTC().Format(f.DateTimeFormat))
	for _, k := range orderedFields {
		if v, ok := entry.Data[k]; ok {
			f.appendKeyValue(b, k, v)
		}
	}

	remaining := []string{}
	for k := range entry.Data {
		skip := false
		for _, ofk := range orderedFields {
			if k == ofk {
				skip = true
				break
			}
		}
		if !skip {
			remaining = append(remaining, k)
		}
	}

	sort.Strings(remaining)
	for _, k := range remaining {
		f.appendKeyValue(b, k, entry.Data[k])
	}

	if len(entry.Message) != 0 {
		entry.Message = strings.TrimSuffix(entry.Message, "\n")
		f.appendKeyValue(b, "msg", entry.Message)
	}
}

func (f *ServerFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *ServerFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *ServerFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
