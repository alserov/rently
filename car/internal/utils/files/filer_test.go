package files

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFiler_Save(t *testing.T) {
	f := NewFiler("images")

	file, err := os.Create("testFile")
	defer os.Remove("testFile")
	defer file.Close()
	require.NoError(t, err)

	_, err = file.Write([]byte("hello"))
	require.NoError(t, err)

	b, err := ioutil.ReadFile("testFile")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = f.Save(b, "uuid", i)
		require.NoError(t, err)
	}
	defer os.RemoveAll("images")

	_, err = os.Stat("images/uuid")
	require.NoError(t, err)
}

func TestFiler_GetImage(t *testing.T) {
	f := NewFiler("images")

	file, err := os.Create("testFile.jpg")
	defer os.Remove("testFile.jpg")
	defer file.Close()
	require.NoError(t, err)

	b, err := io.ReadAll(file)
	require.NoError(t, err)

	err = f.Save(b, "uuid", 0)
	defer os.RemoveAll("images")
	require.NoError(t, err)

	_, err = f.GetImage(fmt.Sprintf("uuid"), 0)
	require.NoError(t, err)
}
