package dropbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharing_CreateSharedLink(t *testing.T) {
	c := client()
	out, err := c.Sharing.CreateSharedLink(ctx, &CreateSharedLinkInput{
		Path:     "/hello.txt",
		ShortURL: true,
	})

	assert.NoError(t, err, "error sharing file")
	assert.Equal(t, "/hello.txt", out.Path)
}

func TestSharing_ListSharedFolder(t *testing.T) {
	c := client()
	out, err := c.Sharing.ListSharedFolders(ctx, &ListSharedFolderInput{
		Limit: 1,
	})

	shared := out.Entries

	assert.NoError(t, err, "listing shared folders")
	assert.NotEmpty(t, out.Entries, "output should be non-empty")

	for out.Cursor != "" {
		out, err = c.Sharing.ListSharedFoldersContinue(ctx, &ListSharedFolderContinueInput{
			Cursor: out.Cursor,
		})

		shared = append(shared, out.Entries...)

		assert.NoError(t, err, "listing shared folders")
		assert.NotEmpty(t, out.Entries, "output should be non-empty")
	}

	for _, sharedFolder := range shared {
		out, err := c.Sharing.ListSharedFolderMembers(ctx, &ListSharedFolderMembersInput{
			SharedFolderID: sharedFolder.SharedFolderID,
			Limit:          1,
		})

		assert.NoError(t, err, "listing shared folder members")
		assert.Equal(t, 1, len(out.Users)+len(out.Groups)+len(out.Invitees), "there should be 1 item present")

		for out.Cursor != "" {
			out, err = c.Sharing.ListSharedFolderMembersContinue(ctx, &ListSharedMembersContinueInput{
				Cursor: out.Cursor,
			})

			assert.NoError(t, err, "listing shared folder members")
			assert.Equal(t, 1, len(out.Users)+len(out.Groups)+len(out.Invitees), "there should be 1 item present")
		}
	}
}

func TestSharing_ListSharedFile(t *testing.T) {
	c := client()
	out, err := c.Sharing.ListSharedFileMembers(ctx, &ListSharedFileMembersInput{
		File:             "/hello.txt",
		IncludeInherited: true,
		Limit:            1,
	})

	assert.NoError(t, err, "getting shared user list")
	assert.NotEmpty(t, out.Users, "output should be non-empty")
	assert.NotEmpty(t, out.Cursor, "cursor should be non-empty for complete test")

	out, err = c.Sharing.ListSharedFileMembersContinue(ctx, &ListSharedMembersContinueInput{
		Cursor: out.Cursor,
	})

	assert.NoError(t, err, "continuing shared user list")
	assert.NotZero(t, len(out.Users)+len(out.Groups)+len(out.Invitees), "continuing list should be non-empty")
}
