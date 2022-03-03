package controllers

import "github.com/yrosaguiar/sentry-operator/pkg/sentry"

type Sentry struct {
	Organization string
	Client       *SentryClient
}

type SentryClient struct {
	Organizations SentryOrganizations
	Projects      SentryProjects
	Teams         SentryTeams
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SentryOrganizations
type SentryOrganizations interface {
	ListProjects(organizationSlug string, opts *sentry.ListOptions) ([]sentry.Project, *sentry.Response, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SentryProjects
type SentryProjects interface {
	Update(organizationSlug, projectSlug string, params *sentry.UpdateProjectParams) (*sentry.Project, *sentry.Response, error)
	Delete(organizationSlug, projectSlug string) (*sentry.Response, error)
	ListKeys(organizationSlug, projectSlug string, opts *sentry.ListOptions) ([]sentry.ProjectKey, *sentry.Response, error)
	CreateKey(organizationSlug, projectSlug string, params *sentry.CreateProjectKeyParams) (*sentry.ProjectKey, *sentry.Response, error)
	UpdateKey(organizationSlug, projectSlug, keyID string, params *sentry.UpdateProjectKeyParams) (*sentry.ProjectKey, *sentry.Response, error)
	DeleteKey(organizationSlug, projectSlug, keyID string) (*sentry.Response, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SentryTeams
type SentryTeams interface {
	List(organizationSlug string, opts *sentry.ListOptions) ([]sentry.Team, *sentry.Response, error)
	Create(organizationSlug string, params *sentry.CreateTeamParams) (*sentry.Team, *sentry.Response, error)
	Update(organizationSlug, teamSlug string, params *sentry.UpdateTeamParams) (*sentry.Team, *sentry.Response, error)
	Delete(organizationSlug, teamSlug string) (*sentry.Response, error)
	CreateProject(organizationSlug, teamSlug string, params *sentry.CreateProjectParams) (*sentry.Project, *sentry.Response, error)
}

func removeFinalizer(finalizers []string, name string) []string {
	var result []string
	for _, item := range finalizers {
		if item == name {
			continue
		}

		result = append(result, item)
	}

	return result
}

func containsFinalizer(finalizers []string, name string) bool {
	for _, finalizer := range finalizers {
		if finalizer == name {
			return true
		}
	}

	return false
}
