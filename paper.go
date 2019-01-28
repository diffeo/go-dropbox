package dropbox

import (
	"context"
	"encoding/json"
	"time"
)

// Paper client for Dropbox Paper.
type Paper struct {
	*Client
}

// NewPaper cretes a new Dropbox Paper client.
func NewPaper(config *Config) *Paper {
	return &Paper{
		Client: &Client{
			Config: config,
		},
	}
}

const (
	// PaperDocsFilterAccessed indicates that for this paper/docs/list call we
	// only want to get documents that the user accessed.
	PaperDocsFilterAccessed = "docs_accessed"
	// PaperDocsFilterCreated indicates that for this paper/docs/list call we
	// only want to get documents that the user created.
	PaperDocsFilterCreated = "docs_created"

	// PaperDocsSortByAccessed indicates that for this paper/docs/list call we
	// want to sort the returned doc IDs in the order in which the documents
	// were accessed by the user
	PaperDocsSortByAccessed = "accessed"
	// PaperDocsSortByModified indicates that for this paper/docs/list call we
	// want to sort the returned doc IDs in the order in which the documents
	// were modified by any user
	PaperDocsSortByModified = "modified"
	// PaperDocsSortByCreated indicates that for this paper/docs/list call we
	// want to sort the returned doc IDs in the order in which the documents
	// were created
	PaperDocsSortByCreated = "created"

	// PaperDocsSortAscending indicates that for this paper/docs/list call we
	// want to return the doc IDs in ascending order for how they are sorted
	PaperDocsSortAscending = "ascending"
	// PaperDocsSortDescending indicates that for this paper/docs/list call we
	// want to return the doc IDs in decending order for how they are sorted
	PaperDocsSortDescending = "decending"
)

// PaperDocsListInput is the payload for /paper/docs/list requests.
type PaperDocsListInput struct {
	FilterBy  string `json:"filter_by"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
	Limit     int    `json:"limit"`
}

// PaperDocsListCursor is a cursor to use in /paper/docs/list/continue calls
// to continue listing papers where the previous list call left off.
type PaperDocsListCursor struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// PaperDocsListOutput is the response format for /paper/docs/list and
// /paper/docs/list/continue.
type PaperDocsListOutput struct {
	DocIDs  []string            `json:"doc_ids"`
	Cursor  PaperDocsListCursor `json:"cursor"`
	HasMore bool                `json:"has_more"`
}

// ListDocs returns the documents in a user's Dropbox Paper.
func (c *Paper) ListDocs(ctx context.Context, in *PaperDocsListInput) (out *PaperDocsListOutput, err error) {
	body, err := c.call(ctx, "/paper/docs/list", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// PaperDocsListContinueInput is the payload for /paper/docs/list/continue
// requests.
type PaperDocsListContinueInput struct {
	Cursor string `json:"cursor"`
}

// ListDocsContinue paginates using the cursor from ListDocs.
func (c *Paper) ListDocsContinue(ctx context.Context, in *PaperDocsListContinueInput) (out *PaperDocsListOutput, err error) {
	body, err := c.call(ctx, "/paper/docs/list/continue", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

const (
	// ExportFormatHTML indicates that we want to download a Dropbox Paper file
	// as HTML.
	ExportFormatHTML = "html"
	// ExportFormatMarkdown indicates that we want to download a Dropbox Paper
	// file as Markdown.
	ExportFormatMarkdown = "markdown"
)

// PaperDownloadInput is the request format for downloading a Dropbox paper,
// including the paper's ID and the format we want to download it in.
type PaperDownloadInput struct {
	DocID        string `json:"doc_id"`
	ExportFormat string `json:"export_format"`
}

// Download a Dropbox Paper.
func (c *Paper) Download(ctx context.Context, in *PaperDownloadInput) (out *DownloadOutput, err error) {
	body, l, err := c.download(ctx, "/paper/docs/download", in, nil)
	if err != nil {
		return
	}

	out = &DownloadOutput{body, l}
	return
}

// PaperGetFolderInfoInput is the request payload format for
// /paper/docs/get_folder_info requests, indicating the ID of the Dropbox Paper
// document that the caller wants to get folder information for.
type PaperGetFolderInfoInput struct {
	DocID string `json:"doc_id"`
}

// PaperGetFolderInfoOutput is the response format for
// /paper/docs/get_folder_info requests, containing information on which
// folders the requested Dropbox Paper is a part of.
type PaperGetFolderInfoOutput struct {
	FolderSharingPolicyType PaperFolderSharingPolicyType `json:"folder_sharing_policy_type"`
	Folders                 []PaperFolder                `json:"folder"`
}

// PaperFolderSharingPolicyType is the folder sharing policy for the folder
// that contains the requested Drobox Paper; can be either "team or
// "invite_only".
type PaperFolderSharingPolicyType struct {
	Tag string `json:".tag"`
}

// PaperFolder contains the ID and display name of a folder containing a
// Dropbox Paper.
type PaperFolder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetFolderInfo returns information on which folders the requested Dropbox
// Paper is part of.
func (c *Paper) GetFolderInfo(ctx context.Context, in *PaperGetFolderInfoInput) (out *PaperGetFolderInfoOutput, err error) {
	body, err := c.call(ctx, "/paper/docs/get_folder_info", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}
