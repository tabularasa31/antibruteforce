package repo_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
	"github.com/tabularasa31/antibruteforce/pkg/postgres"
	"golang.org/x/net/context"
)

var (
	pg       *postgres.Postgres
	listrepo *repo.ListRepo
)

type req struct {
	subnet string
	color  string
}

type testCase struct {
	description  string
	input        req
	expectedBoo  bool
	expectedMess string
	expectedErr  error
}

var (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	port     = "5433"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	maxConn  = 25
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	dsn = fmt.Sprintf(dsn, user, password, port, db)
	if err = pool.Retry(
		func() error {
			pg, err = postgres.New(&config.Config{
				Postgres: config.Postgres{Dsn: dsn, PoolMax: maxConn},
			})
			if err != nil {
				return err
			}
			return pg.Pool.Ping(context.Background())
		}); err != nil {
		log.Errorf("Could not connect to docker: %s", err.Error())
	}

	listrepo = repo.NewListRepo(pg)

	err = listrepo.Drop()
	if err != nil {
		panic(err)
	}

	err = listrepo.Up()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Errorf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestListRepo_SaveToList(t *testing.T) {
	testCases := []testCase{
		{
			description: "success test",
			input: req{
				subnet: "168.0.1.0/24",
				color:  "white",
			},
			expectedBoo:  true,
			expectedMess: "",
			expectedErr:  nil,
		},
		{
			description: "empty subnet",
			input: req{
				subnet: "",
				color:  "white",
			},
			expectedBoo:  false,
			expectedMess: "",
		},
		{
			description: "try add subnet same as already exists",
			input: req{
				subnet: "168.0.1.0/24",
				color:  "white",
			},
			expectedBoo:  false,
			expectedMess: "given subnet 168.0.1.0/24 already in whitelist",
		},
		{
			description: "list conflict",
			input: req{
				subnet: "168.0.1.0/24",
				color:  "black",
			},
			expectedBoo:  false,
			expectedMess: "list conflict: can't add given subnet 168.0.1.0/24 in blacklist because it is already in whitelist",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			boo, mess, _ := listrepo.SaveToList(context.Background(), tc.input.subnet, tc.input.color)
			assert.Equal(t, tc.expectedMess, mess)
			assert.Equal(t, tc.expectedBoo, boo)
		})
	}
}
