package controllers_test

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	sentryv1alpha1 "github.com/yrosaguiar/sentry-operator/api/v1alpha1"
	"github.com/yrosaguiar/sentry-operator/controllers"
	"github.com/yrosaguiar/sentry-operator/controllers/controllersfakes"
	"github.com/yrosaguiar/sentry-operator/pkg/sentry"
	// +kubebuilder:scaffold:imports
)

var (
	testEnv    *envtest.Environment
	k8sClient  client.Client
	k8sManager ctrl.Manager
)

var (
	fakeSentryOrganizations *controllersfakes.FakeSentryOrganizations
	fakeSentryProjects      *controllersfakes.FakeSentryProjects
	fakeSentryTeams         *controllersfakes.FakeSentryTeams
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	log.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = sentryv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).ToNot(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	fakeSentryOrganizations = new(controllersfakes.FakeSentryOrganizations)
	fakeSentryProjects = new(controllersfakes.FakeSentryProjects)
	fakeSentryTeams = new(controllersfakes.FakeSentryTeams)

	ctrlSentry := &controllers.Sentry{
		Organization: "organization",
		Client: &controllers.SentryClient{
			Organizations: fakeSentryOrganizations,
			Projects:      fakeSentryProjects,
			Teams:         fakeSentryTeams,
		},
	}

	err = (&controllers.ProjectReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Project"),
		Scheme: k8sManager.GetScheme(),
		Sentry: ctrlSentry,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&controllers.ProjectKeyReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ProjectKey"),
		Scheme: k8sManager.GetScheme(),
		Sentry: ctrlSentry,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&controllers.TeamReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Team"),
		Scheme: k8sManager.GetScheme(),
		Sentry: ctrlSentry,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	time.Sleep(3 * time.Second)
})

var _ = AfterSuite(func() {
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func testSentryProject(id, team, name string) *sentry.Project {
	return &sentry.Project{
		DateCreated: time.Now(),
		ID:          id,
		Name:        name,
		Slug:        name,
		Team: sentry.Team{
			Slug: team,
		},
	}
}

func testSentryProjectKey(id string, projectID int, name, dsn string) *sentry.ProjectKey {
	return &sentry.ProjectKey{
		DateCreated: time.Now(),
		ID:          id,
		Name:        name,
		ProjectID:   projectID,
		DSN: sentry.ProjectKeyDSN{
			Public: dsn,
		},
	}
}

func testSentryTeam(id, name string) *sentry.Team {
	return &sentry.Team{
		DateCreated: time.Now(),
		ID:          id,
		Name:        name,
		Slug:        name,
	}
}

func newSentryResponse(statusCode int) *sentry.Response {
	return &sentry.Response{
		Response: &http.Response{
			StatusCode: statusCode,
		},
		NextPage: &sentry.Page{
			Cursor:  "0:0:0",
			Results: false,
		},
	}
}
