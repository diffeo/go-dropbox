package dropbox

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaper_List(t *testing.T) {
	ctx := context.Background()
	c := client()

	out, err := c.Paper.ListDocs(ctx, &PaperDocsListInput{
		SortBy:    "modified",
		SortOrder: "ascending",
		Limit:     10,
	})

	require.NoError(t, err)
	assert.Equal(t, 10, len(out.DocIDs))
	assert.NotZero(t, out.Cursor.Value)
	assert.False(t, out.Cursor.Expiration.IsZero())
}

func TestPaper_Download(t *testing.T) {
	ctx := context.Background()
	c := client()
	var fileID string
	if fileID = os.Getenv("DROPBOX_PAPER_TOP_LEVEL_FILE_ID"); fileID == "" {
		t.Skip("No top-level file ID passed in")
	}

	out, err := c.Paper.Download(ctx, &PaperDownloadInput{
		DocID:        fileID,
		ExportFormat: "markdown",
	})
	require.NoError(t, err)

	defer out.Body.Close()
	content, err := ioutil.ReadAll(out.Body)
	require.NoError(t, err)
	assert.Equal(t, out.Length, int64(len(content)))
}

func TestPaper_GetFolderInfoTopLevel(t *testing.T) {
	ctx := context.Background()
	c := client()

	var fileID string
	if fileID = os.Getenv("DROPBOX_PAPER_TOP_LEVEL_FILE_ID"); fileID == "" {
		t.Skip("No top-level file ID passed in")
	}

	out, err := c.Paper.GetFolderInfo(ctx, &PaperGetFolderInfoInput{DocID: fileID})
	require.NoError(t, err)
	assert.Zero(t, out.FolderSharingPolicyType.Tag)
	assert.Nil(t, out.Folders)
}

func TestPaper_GetFolderInfoInsideFolders(t *testing.T) {
	ctx := context.Background()
	c := client()

	var fileID string
	if fileID = os.Getenv("DROPBOX_PAPER_IN_FOLDERS_FILE_ID"); fileID == "" {
		t.Skip("No top-level file ID passed in")
	}

	out, err := c.Paper.GetFolderInfo(ctx, &PaperGetFolderInfoInput{DocID: fileID})
	require.NoError(t, err)
	assert.Equal(t, "invite_only", out.FolderSharingPolicyType.Tag)
	require.Len(t, out.Folders, 2)
	for _, f := range out.Folders {
		assert.NotZero(t, f.ID)
		assert.NotZero(t, f.Name)
	}
}

func TestPaper_AlphaGetMetadata(t *testing.T) {
	ctx := context.Background()
	c := client()

	var fileID string
	if fileID = os.Getenv("DROPBOX_PAPER_TOP_LEVEL_FILE_ID"); fileID == "" {
		t.Skip("No top-level file ID passed in")
	}

	out, err := c.Paper.AlphaGetMetadata(ctx, &PaperGetMetadataInput{DocID: fileID})
	require.NoError(t, err)
	assert.Equal(t, fileID, out.DocID)
	assert.NotZero(t, out.Owner)
	assert.NotZero(t, out.Title)
	assert.False(t, out.CreatedDate.IsZero())
	assert.Equal(t, "active", out.Status.Tag)
	assert.False(t, out.LastUpdatedDate.IsZero())
	assert.NotZero(t, out.Revision)
	assert.NotZero(t, out.LastEditor)
}
