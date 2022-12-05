package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/izaakdale/service-log/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T, log *Log){
		"append and read a record":   testAppendRead,
		"offset out of range":        testOutOfRange,
		"init with existing segment": testInitExisting,
		"reader":                     testReader,
		"truncate":                   testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := New(dir, c)
			require.NoError(t, err)
			fn(t, log)
		})
	}
}

func testAppendRead(t *testing.T, l *Log) {
	append := &api.Record{
		Value: []byte("test"),
	}
	off, err := l.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := l.Read(off)
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
}

func testOutOfRange(t *testing.T, l *Log) {
	read, err := l.Read(1)
	require.Nil(t, read)
	require.Error(t, err)
}

func testInitExisting(t *testing.T, l *Log) {
	append := &api.Record{
		Value: []byte("test"),
	}
	for i := 0; i < 3; i++ {
		_, err := l.Append(append)
		require.NoError(t, err)
	}
	require.NoError(t, l.Close())
	off, err := l.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = l.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

	n, err := New(l.Dir, l.Config)
	require.NoError(t, err)

	off, err = n.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = n.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)
}

func testReader(t *testing.T, l *Log) {
	append := &api.Record{
		Value: []byte("test"),
	}
	off, err := l.Append(append)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	reader := l.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, append.Value, read.Value)
}

func testTruncate(t *testing.T, l *Log) {
	append := &api.Record{
		Value: []byte("test"),
	}
	for i := 0; i < 3; i++ {
		_, err := l.Append(append)
		require.NoError(t, err)
	}
	err := l.Truncate(2)
	require.NoError(t, err)
	_, err = l.Read(1)
	require.Error(t, err)
}
