package storage

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

const imageBytes = "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"

func TestImageStorage(t *testing.T) {
	f := NewImageStorage()

	id, err := f.Save(context.Background(), "key1", bytes.NewReader([]byte(imageBytes)))
	require.NoError(t, err)
	require.NotEmpty(t, id)

	file, err := f.Get(context.Background(), "key1/"+id)
	require.NoError(t, err)

	content, err := io.ReadAll(file)
	require.NoError(t, err)
	file.Close()

	require.Equal(t, []byte(imageBytes), content)

	err = f.Delete(context.Background(), "key1/"+id)
	require.NoError(t, err)
}