package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type repo struct {
	sql  string
	err  error
	rows *sqlmock.Rows
}

func TestRepository_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type param struct {
		id uuid.UUID
	}

	type args struct {
		p    param
		repo repo
	}

	columns := []string{"id", "first_name", "last_name", "birthday", "created_at", "updated_at"}

	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr error
	}{
		{
			name: "success",
			args: args{
				p: param{
					id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				},
				repo: repo{
					sql: prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users WHERE id=$1`),
					err: nil,
					rows: sqlmock.NewRows(columns).AddRow(
						"ccae37ea-d41e-4371-a3a3-89203b9e2608",
						"Elon",
						"Musk",
						"1971-06-28",
						time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					),
				},
			},
			want: &User{
				ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  "1971-06-28",
				CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				p: param{
					id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				},
				repo: repo{
					sql:  prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users WHERE id=$1`),
					err:  sql.ErrNoRows,
					rows: sqlmock.NewRows(columns),
				},
			},
			want:    nil,
			wantErr: errors.New("not exists"),
		},
		{
			name: "some err",
			args: args{
				p: param{
					id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				},
				repo: repo{
					sql:  prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users WHERE id=$1`),
					err:  errors.New("some err"),
					rows: sqlmock.NewRows(columns),
				},
			},
			want:    nil,
			wantErr: errors.New("exec: some err"),
		},
	}

	r := Repository{
		db: sqlxDB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(tt.args.repo.sql).
				WithArgs(&tt.args.p.id).
				WillReturnRows(tt.args.repo.rows).
				WillReturnError(tt.args.repo.err)

			got, err := r.Get(context.Background(), tt.args.p.id)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type param struct {
		id uuid.UUID
	}

	type args struct {
		p    param
		repo repo
	}

	columns := []string{"id", "first_name", "last_name", "birthday", "created_at", "updated_at"}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				p: param{
					id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				},
				repo: repo{
					sql: prepareSQL(`DELETE FROM users WHERE id=$1`),
					err: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "some err",
			args: args{
				p: param{
					id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				},
				repo: repo{
					sql:  prepareSQL(`DELETE FROM users WHERE id=$1`),
					err:  errors.New("some err"),
					rows: sqlmock.NewRows(columns),
				},
			},
			wantErr: errors.New("exec: some err"),
		},
	}

	r := Repository{
		db: sqlxDB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(tt.args.repo.sql).
				WithArgs(&tt.args.p.id).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.args.repo.err)

			err := r.Delete(context.Background(), tt.args.p.id)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type param struct {
		user *User
	}

	type args struct {
		p    param
		repo repo
	}

	columns := []string{"id", "first_name", "last_name", "birthday", "created_at", "updated_at"}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				p: param{
					user: &User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
				repo: repo{
					sql: prepareSQL(`UPDATE users SET first_name=$1, last_name=$2, birthday=$3, updated_at=$4 WHERE id=$5`),
					err: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "some err",
			args: args{
				p: param{
					user: &User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
				repo: repo{
					sql:  prepareSQL(`UPDATE users SET first_name=$1, last_name=$2, birthday=$3, updated_at=$4 WHERE id=$5`),
					err:  errors.New("some err"),
					rows: sqlmock.NewRows(columns),
				},
			},
			wantErr: errors.New("exec: some err"),
		},
	}

	r := Repository{
		db: sqlxDB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(tt.args.repo.sql).
				WithArgs(
					&tt.args.p.user.FirstName,
					&tt.args.p.user.LastName,
					&tt.args.p.user.Birthday,
					&tt.args.p.user.UpdatedAt,
					&tt.args.p.user.ID,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.args.repo.err)

			err := r.Update(context.Background(), tt.args.p.user)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}
		})
	}
}

func TestRepository_Created(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type param struct {
		user *User
	}

	type args struct {
		p    param
		repo repo
	}

	columns := []string{"id", "first_name", "last_name", "birthday", "created_at", "updated_at"}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				p: param{
					user: &User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
				repo: repo{
					sql: prepareSQL(`INSERT INTO users (id, first_name, last_name, birthday, created_at) VALUES($1, $2, $3, $4, $5)`),
					err: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "some err",
			args: args{
				p: param{
					user: &User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
				repo: repo{
					sql:  prepareSQL(`INSERT INTO users (id, first_name, last_name, birthday, created_at) VALUES($1, $2, $3, $4, $5)`),
					err:  errors.New("some err"),
					rows: sqlmock.NewRows(columns),
				},
			},
			wantErr: errors.New("exec: some err"),
		},
	}

	r := Repository{
		db: sqlxDB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(tt.args.repo.sql).
				WithArgs(
					&tt.args.p.user.ID,
					&tt.args.p.user.FirstName,
					&tt.args.p.user.LastName,
					&tt.args.p.user.Birthday,
					&tt.args.p.user.CreatedAt,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.args.repo.err)

			err := r.Create(context.Background(), tt.args.p.user)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type args struct {
		repo repo
	}

	columns := []string{"id", "first_name", "last_name", "birthday", "created_at", "updated_at"}

	tests := []struct {
		name    string
		args    args
		want    []*User
		wantErr error
	}{
		{
			name: "success",
			args: args{
				repo: repo{
					sql: prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users`),
					err: nil,
					rows: sqlmock.NewRows(columns).AddRow(
						"ccae37ea-d41e-4371-a3a3-89203b9e2608",
						"Elon",
						"Musk",
						"1971-06-28",
						time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					),
				},
			},
			want: []*User{
				{
					ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					FirstName: "Elon",
					LastName:  "Musk",
					Birthday:  "1971-06-28",
					CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: nil,
		},
		{
			name: "scan err",
			args: args{
				repo: repo{
					sql: prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users`),
					err: nil,
					rows: sqlmock.NewRows(columns).AddRow(
						"ccae37ea-d41e-4371-a3a3-89203b9e2608",
						"Elon",
						"Musk",
						"1971-06-28",
						time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						123,
					),
				},
			},
			want:    nil,
			wantErr: errors.New("scan: sql: Scan error on column index 5, name \"updated_at\": unsupported Scan, storing driver.Value type int64 into type *time.Time"),
		},
		{
			name: "some err",
			args: args{
				repo: repo{
					sql:  prepareSQL(`SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users`),
					err:  errors.New("some err"),
					rows: sqlmock.NewRows(columns),
				},
			},
			want:    nil,
			wantErr: errors.New("query: some err"),
		},
	}

	r := Repository{
		db: sqlxDB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(tt.args.repo.sql).
				WillReturnRows(tt.args.repo.rows).
				WillReturnError(tt.args.repo.err)

			got, err := r.List(context.Background())
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func prepareSQL(sql string) string {
	replacer := strings.NewReplacer("$", "\\$", "(", "\\(", ")", "\\)")
	return replacer.Replace(sql)
}

func TestNewRepository(t *testing.T) {
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			name: "success",
			args: args{
				db: &sqlx.DB{},
			},
			want: &Repository{
				db: &sqlx.DB{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRepository(tt.args.db), "NewRepository(%v)", tt.args.db)
		})
	}
}
