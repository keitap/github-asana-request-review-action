module github.com/keitap/github-asana-request-review-action

go 1.24.0

replace bitbucket.org/mikehouston/asana-go => github.com/keitap/asana-go v0.0.0-20210425173123-936f596fb971

require (
	bitbucket.org/mikehouston/asana-go v0.0.0-20250102231814-14e44a300f0b
	github.com/google/go-github/v71 v71.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/oauth2 v0.27.0
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
)
