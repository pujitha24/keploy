package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.keploy.io/server/v3/pkg/models"
)

func TestPerformBodyMatch_ArrayBodies(t *testing.T) {
	h := newHTTP()
	mk := func(name, body string) *models.Mock {
		m := httpMock(name, "POST", "http://dep/batch")
		m.Spec.HTTPReq.Body = body
		return m
	}
	tests := []struct {
		name   string
		mocks  []*models.Mock
		req    string
		expect []string
	}{
		{"array with changed value", []*models.Mock{mk("a", `[{"id":1,"ts":"x"}]`)}, `[{"id":1,"ts":"y"}]`, []string{"a"}},
		{"array vs object", []*models.Mock{mk("a", `[1]`)}, `{"k":1}`, nil},
		{"object vs array", []*models.Mock{mk("a", `{"k":1}`)}, `[1]`, nil},
		{"undecodable candidate does not drop the rest", []*models.Mock{mk("bad", `not json`), mk("a", `{"k":1}`)}, `{"k":2}`, []string{"a"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := h.PerformBodyMatch(context.Background(), tc.mocks, []byte(tc.req))
			require.NoError(t, err)
			var names []string
			for _, m := range got {
				names = append(names, m.Name)
			}
			require.Equal(t, tc.expect, names)
		})
	}
}
