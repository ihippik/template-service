//go:build integration
// +build integration

package user

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/ihippik/template-service/migrations"
)

type RepositoryTestSuite struct {
	suite.Suite

	db       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (s *RepositoryTestSuite) SetupSuite() {
	var err error

	s.pool, err = dockertest.NewPool("")
	s.NoError(err)

	s.resource, err = s.pool.Run(
		"postgres",
		"14.5-alpine",
		[]string{
			"POSTGRES_PASSWORD=pass",
			"POSTGRES_DB=template_db",
		},
	)
	s.NoError(err)

	const containerLifetime = 120

	err = s.resource.Expire(containerLifetime)
	s.NoError(err)

	dsn := "postgres://postgres:pass@localhost:" + s.resource.GetPort("5432/tcp") + "/template_db?sslmode=disable"

	err = s.pool.Retry(func() error {
		var err error

		s.db, err = sqlx.Connect(
			"postgres",
			dsn,
		)
		if err != nil {
			return err
		}

		return s.db.Ping()
	})
	s.NoError(err)

	err = migrations.Up(dsn)
	s.NoError(err)
}

func TestRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suiteTester := new(RepositoryTestSuite)
	suite.Run(t, suiteTester)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if err := s.pool.Purge(s.resource); err != nil {
		s.Fail(err.Error())
	}
}

func (s *RepositoryTestSuite) TestEndpoint() {
	logger := zap.NewNop()
	r := NewRepository(s.db)
	svc := NewService(nil, logger, r)
	endpoint := NewEndpoint(logger, svc)

	var dto = []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`)

	// create user
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(dto))
	w := httptest.NewRecorder()
	endpoint.CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(s.T(), http.StatusCreated, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(s.T(), err)

	var resp response

	err = json.Unmarshal(data, &resp)
	assert.NoError(s.T(), err)

	// get user
	getReq := httptest.NewRequest(http.MethodPost, "/v1/users/id", bytes.NewReader(dto))
	getW := httptest.NewRecorder()
	getReq = mux.SetURLVars(getReq, map[string]string{"id": resp.Data[0].ID.String()})
	endpoint.GetUser(getW, getReq)

	getRes := getW.Result()
	defer getRes.Body.Close()
	assert.Equal(s.T(), http.StatusOK, getRes.StatusCode)

	getData, err := io.ReadAll(getRes.Body)
	assert.NoError(s.T(), err)

	var getResp response

	err = json.Unmarshal(getData, &getResp)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), resp.Data[0].ID, getResp.Data[0].ID)
}
