package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanPairs_ExactMatch(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "overlay.foo.png"), []byte{}, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "map.foo.png"), []byte{}, 0644))

	pairs, warnings, err := scanPairs(dir)
	require.NoError(t, err)
	assert.Len(t, pairs, 1)
	assert.Equal(t, "foo", pairs[0].Name)
	assert.Equal(t, filepath.Join(dir, "overlay.foo.png"), pairs[0].OverlayPath)
	assert.Equal(t, filepath.Join(dir, "map.foo.png"), pairs[0].MapPath)
	assert.Equal(t, filepath.Join(dir, "source.foo.png"), pairs[0].SourcePath)
	assert.Empty(t, warnings)
}

func TestScanPairs_NoMatch(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "overlay.foo.png"), []byte{}, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "map.bar.png"), []byte{}, 0644))

	pairs, warnings, err := scanPairs(dir)
	require.NoError(t, err)
	assert.Empty(t, pairs)
	assert.Equal(t, []string{"foo"}, warnings)
}

func TestScanPairs_MultipleMatches(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, "overlay."+name+".png"), []byte{}, 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "map."+name+".png"), []byte{}, 0644))
	}

	pairs, warnings, err := scanPairs(dir)
	require.NoError(t, err)
	assert.Len(t, pairs, 3)
	assert.Empty(t, warnings)

	names := make(map[string]bool)
	for _, p := range pairs {
		names[p.Name] = true
	}
	assert.True(t, names["alpha"])
	assert.True(t, names["beta"])
	assert.True(t, names["gamma"])
}

func TestScanPairs_InvalidDir(t *testing.T) {
	pairs, warnings, err := scanPairs("/nonexistent/path")
	assert.Error(t, err)
	assert.Nil(t, pairs)
	assert.Nil(t, warnings)
}

func TestScanPairs_DottedName(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "overlay.a.b.png"), []byte{}, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "map.a.b.png"), []byte{}, 0644))

	pairs, warnings, err := scanPairs(dir)
	require.NoError(t, err)
	assert.Len(t, pairs, 1)
	assert.Equal(t, "a.b", pairs[0].Name)
	assert.Empty(t, warnings)
}

func TestScanPairs_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	pairs, warnings, err := scanPairs(dir)
	require.NoError(t, err)
	assert.Empty(t, pairs)
	assert.Empty(t, warnings)
}
