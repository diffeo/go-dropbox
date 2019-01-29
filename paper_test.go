package dropbox

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: This test assumes that you have 10 or more Papers in your Dropbox
// Paper directory tree.
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

func dummyPaperTitle() string { return fmt.Sprintf("go-dropbox-client-test-%d", rand.Intn(100000)) }

func setupPaperDocument(ctx context.Context, t *testing.T, c *Client) (string, func()) {
	o, err := c.Paper.Create(ctx, &PaperCreateInput{
		ImportFormat: "markdown",
		Reader:       strings.NewReader("# " + dummyPaperTitle()),
	})
	require.NoError(t, err)
	return o.DocID, func() {
		require.NoError(t, c.Paper.PermanentlyDelete(ctx, &PaperPermanentlyDeleteInput{DocID: o.DocID}))
	}
}

func TestPaper_CreateDeletePermanently(t *testing.T) {
	ctx := context.Background()
	c := client()
	title := dummyPaperTitle()

	// Create a Dropbox Paper
	o, err := c.Paper.Create(ctx, &PaperCreateInput{
		ImportFormat: "markdown",
		Reader:       strings.NewReader("# " + title),
	})
	require.NoError(t, err)
	require.NotZero(t, o.DocID)
	assert.Equal(t, title, o.Title)
	assert.NotZero(t, o.Revision)

	// We should not get an error running GetFolderInfo for this Dropbox Paper
	// since it exists
	_, err = c.Paper.GetFolderInfo(ctx, &PaperGetFolderInfoInput{DocID: o.DocID})
	require.NoError(t, err)

	// Permanently the Dropbox Paper
	require.NoError(t, c.Paper.PermanentlyDelete(
		ctx, &PaperPermanentlyDeleteInput{DocID: o.DocID}))

	// GetFolderInfo should now error since this Dropbox Paper no longer exists
	_, err = c.Paper.GetFolderInfo(ctx, &PaperGetFolderInfoInput{DocID: o.DocID})
	if assert.Error(t, err) {
		err, ok := err.(*Error)
		require.True(t, ok)
		assert.True(t, strings.HasPrefix(err.Summary, "doc_not_found/"),
			`error summary %s should have prefix doc_not_found/`)
	}
}

func TestPaper_Download(t *testing.T) {
	ctx := context.Background()
	c := client()
	fileID, deleteFile := setupPaperDocument(ctx, t, c)
	defer deleteFile()

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

	fileID, deleteFile := setupPaperDocument(ctx, t, c)
	defer deleteFile()

	out, err := c.Paper.GetFolderInfo(ctx, &PaperGetFolderInfoInput{DocID: fileID})
	require.NoError(t, err)
	assert.Zero(t, out.FolderSharingPolicyType.Tag)
	assert.Nil(t, out.Folders)
}

// NOTE: This test requires that the caller has a Dropbox Paper file that is
// nested two folders deep into their Dropbox Paper directory tree, which is
// set on the environment variable DROPBOX_PAPER_NESTED_FILE_ID. The reason why
// this cannot be automated the same way creating and testing endpoints on a
// file at the top level can't is because currently there is no create folder
// endpoint on the Dropbox Paper API currently, therefore we cannot just create
// a throwaway folder.
//
// If you are running this test and looking for an ID to use, what you can do
// is run the Dropbox /paper/docs/list endpoint with your OAuth token to get
// the list of IDs of your Dropbox Papers, then pass those Dropbox IDs to the
// Dropbox /paper/docs/get_folder_info endpoint until you find one where the
// endpoint returns two folders. Then set that file's ID as
// DROPBOX_PAPER_NESTED_FILE_ID.

func TestPaper_GetFolderInfoInsideFolders(t *testing.T) {
	ctx := context.Background()
	c := client()

	var fileID string
	if fileID = os.Getenv("DROPBOX_PAPER_NESTED_FILE_ID"); fileID == "" {
		t.Skip("No nested file ID passed in")
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

	fileID, deleteFile := setupPaperDocument(ctx, t, c)
	defer deleteFile()

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
