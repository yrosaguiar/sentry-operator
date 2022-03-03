package sentry_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yrosaguiar/sentry-operator/pkg/sentry"
)

func TestSentry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg/sentry")
}

func setup() (*http.ServeMux, *sentry.Client) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Endpoint not found. Please check your handler's registered pattern."}`))
	})

	server := httptest.NewServer(handler)
	serverURL, _ := url.Parse(server.URL)
	client := sentry.NewClient("token", sentry.WithSentryURL(serverURL))

	return handler, client
}

func testHandler(method string, handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()

		Expect(r.Method).To(Equal(method))
		Expect(r.Header.Get("Authorization")).To(Equal("Bearer token"))
		Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))
		Expect(r.Header.Get("Accept")).To(Equal("application/json"))

		w.Header().Set("Content-Type", "application/json")
		handlerFunc(w, r)
	}
}

func newPaginationLinks() string {
	return `<https://sentry.io/api/0/previous/?&cursor=0:0:1>; rel="previous"; results="true"; cursor="0:0:1", <https://sentry.io/api/0/next/?&cursor=0:0:0>; rel="next"; results="false"; cursor="0:0:0"`
}

func newAPIError(apiErr sentry.APIError) []byte {
	defer GinkgoRecover()

	data, err := json.Marshal(apiErr)
	Expect(err).ToNot(HaveOccurred())

	return data
}

func parseTime(value string) time.Time {
	defer GinkgoRecover()

	result, err := time.Parse(time.RFC3339, value)
	Expect(err).ToNot(HaveOccurred())

	return result
}
