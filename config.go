package dropbox

import (
	"net/http"
)

// Config for the Dropbox clients.
type Config struct {
	HTTPClient  *http.Client
	AccessToken string
	Namespace   *APIPathRoot
}

// NewConfig with the given access token.
func NewConfig(accessToken string) *Config {
	return &Config{
		HTTPClient:  http.DefaultClient,
		AccessToken: accessToken,
	}
}

// APIPathRoot is marshalled onto the Dropbox-API-Path-Root header to
// indicate which namespace to send Dropbox API requests relative to.
// doc: https://www.dropbox.com/developers/reference/namespace-guide
type APIPathRoot struct {
	Tag         string `json:".tag"`
	NamespaceID string `json:"namespace_id,omitempty"`
	Root        string `json:"root,omitempty"`
}

// HomeNamespace is a Dropbox-API-Path-Root header value used for specifying
// that Dropbox API requests should be sent relative to a user's home
// namespace.
var HomeNamespace = &APIPathRoot{Tag: "home"}

// NamespaceIDNamespace is a Dropbox-API-Path-Root header value used for
// specifying that Dropbox API requests should be sent relative to the Dropbox
// namespace passed in.
func NamespaceIDNamespace(namespaceID string) *APIPathRoot {
	return &APIPathRoot{Tag: "namespace_id", NamespaceID: namespaceID}
}

// RootNamespace is a Dropbox-API-Path-Root header value used for specifying
// that Dropbox API requests should be sent relative to the user's root Dropbox
// namespace. If a request is sent relative to a namespace that isn't a user's
// root namespace, then an error occurs.
func RootNamespace(rootNamespaceID string) *APIPathRoot {
	return &APIPathRoot{Tag: "root", Root: rootNamespaceID}
}
