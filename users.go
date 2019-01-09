package dropbox

import (
	"encoding/json"
)

// Users client for user accounts.
type Users struct {
	*Client
}

// NewUsers client.
func NewUsers(config *Config) *Users {
	return &Users{
		Client: &Client{
			Config: config,
		},
	}
}

// GetAccountInput request input.
type GetAccountInput struct {
	AccountID string `json:"account_id"`
}

// GetAccountOutput request output.
type GetAccountOutput struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName    string `json:"given_name"`
		Surname      string `json:"surname"`
		FamiliarName string `json:"familiar_name"`
		DisplayName  string `json:"display_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	IsTeammate    bool   `json:"is_teammate"`
}

// GetAccount returns information about a user's account.
func (c *Users) GetAccount(in *GetAccountInput) (out *GetAccountOutput, err error) {
	body, err := c.call("/users/get_account", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// GetAccountBatchInput request input. At most 300 ids may be listed.
type GetAccountBatchInput struct {
	AccountIDs []string `json:"account_ids"`
}

// GetAccountBatchOutput request output.
type GetAccountBatchOutput []*GetAccountOutput

// GetAccountBatch returns a list of information about users' accounts.
func (c *Users) GetAccountBatch(in *GetAccountBatchInput) (out GetAccountBatchOutput, err error) {
	body, err := c.call("/users/get_account_batch", in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// FullTeam represents a Dropbox team.
type FullTeam struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	SharingPolicies   TeamSharingPolicies `json:"sharing_policies"`
	OfficeAddinPolicy struct {
		Tag string `json:"tag"`
	} `json:"office_addin_policy"`
}

// TeamSharingPolicies represents the sharing policies for a Dropbox team.
type TeamSharingPolicies struct {
	SharedFolderMemberPolicy struct {
		Tag string `json:"tag"`
	} `json:"shared_folder_member_policy"`
	SharedFolderJoinPolicy struct {
		Tag string `json:"tag"`
	} `json:"shared_folder_join_policy"`
	SharedLinkCreatePolicy struct {
		Tag string `json:"tag"`
	} `json:"shared_link_create_policy"`
}

// GetCurrentAccountOutput request output.
type GetCurrentAccountOutput struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName    string `json:"given_name"`
		Surname      string `json:"surname"`
		FamiliarName string `json:"familiar_name"`
		DisplayName  string `json:"display_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	Country       string `json:"country"`
	Locale        string `json:"locale"`
	ReferralLink  string `json:"referral_link"`
	TeamMemberID  string `json:"team_member_id"`
	IsPaired      bool   `json:"is_paired"`
	AccountType   struct {
		Tag string `json:".tag"`
	} `json:"account_type"`
	RootInfo struct {
		Tag             string `json:".tag"`
		RootNamespaceID string `json:"root_namespace_id"`
		HomeNamespaceID string `json:"home_namespace_id"`
		HomePath        string `json:"home_path"`
	} `json:"root_info"`
}

// GetCurrentAccount returns information about the current user's account.
func (c *Users) GetCurrentAccount() (out *GetCurrentAccountOutput, err error) {
	body, err := c.call("/users/get_current_account", nil)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// GetSpaceUsageOutput request output.
type GetSpaceUsageOutput struct {
	Used       uint64 `json:"used"`
	Allocation struct {
		Used      uint64 `json:"used"`
		Allocated uint64 `json:"allocated"`
	} `json:"allocation"`
}

// GetSpaceUsage returns space usage information for the current user's account.
func (c *Users) GetSpaceUsage() (out *GetSpaceUsageOutput, err error) {
	body, err := c.call("/users/get_space_usage", nil)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}
