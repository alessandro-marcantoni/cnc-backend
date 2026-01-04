package persistence

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	m "github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

//go:embed queries/get_member_by_id.sql
var getMemberByIdQuery string

type SQLMemberRepository struct {
	db *sql.DB
}

func NewSQLMemberRepository(db *sql.DB) *SQLMemberRepository {
	return &SQLMemberRepository{db: db}
}

func (r *SQLMemberRepository) GetMemberById(id domain.Id[m.Member]) result.Result[m.Member] {
	var resultRow QueryResult
	err := r.db.QueryRowContext(context.Background(), getMemberByIdQuery, id.Value).Scan(
		&resultRow.MemberID,
		&resultRow.FirstName,
		&resultRow.LastName,
		&resultRow.DateOfBirth,
		&resultRow.Email,
		&resultRow.PhoneNumbers,
		&resultRow.Addresses,
		&resultRow.Memberships,
		&resultRow.RentedServices,
	)

	return result.MapErr(result.Bind(result.From(true, err), func(_ bool) result.Result[m.Member] {
		return MapToStructs(resultRow)
	}), func(err error) error {
		if err == sql.ErrNoRows {
			return errors.NotFoundError{Description: "Member not found"}
		}
		return errors.RepositoryError{Description: err.Error()}
	})
}
