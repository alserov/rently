package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHash(t *testing.T) {
	h, err := hash("1787_1781")
	require.NoError(t, err)
	require.NotEmpty(t, h)
	fmt.Println(h)
}
