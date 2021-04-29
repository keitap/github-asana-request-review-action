package githubasana

import (
	"fmt"

	"bitbucket.org/mikehouston/asana-go"
)

type Account struct {
	Name         string
	AsanaUserGID string
	GitHubLogin  string
}

func NewNoAsanaAccount(githubLogin string) *Account {
	return &Account{
		Name:         githubLogin,
		AsanaUserGID: "",
		GitHubLogin:  githubLogin,
	}
}

func NewAccount(client *asana.Client, userGID string, githubLogin string) (*Account, error) {
	u := &asana.User{ID: userGID}

	err := u.Fetch(client)
	if err != nil {
		return nil, err
	}

	return &Account{
		Name:         u.Name,
		AsanaUserGID: u.ID,
		GitHubLogin:  githubLogin,
	}, nil
}

func (u *Account) GetUserPermalink() string {
	if u.AsanaUserGID == "" {
		return fmt.Sprintf(`<a href="https://github.com/%s">%s</a>`, u.GitHubLogin, u.Name)
	}

	return fmt.Sprintf(`<a data-asana-gid="%s"/>`, u.AsanaUserGID)
}

func (u *Account) String() string {
	if u.AsanaUserGID == "" {
		return fmt.Sprintf("@%s", u.GitHubLogin)
	}

	return fmt.Sprintf("%s (@%s)", u.Name, u.GitHubLogin)
}
