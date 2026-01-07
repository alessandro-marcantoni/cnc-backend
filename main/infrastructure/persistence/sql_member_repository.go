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

//go:embed queries/get_all_members.sql
var getAllMembersQuery string

//go:embed queries/get_members_by_season.sql
var getMembersBySeasonQuery string

type SQLMemberRepository struct {
	db *sql.DB
}

func NewSQLMemberRepository(db *sql.DB) *SQLMemberRepository {
	return &SQLMemberRepository{db: db}
}

func (r *SQLMemberRepository) GetAllMembers() result.Result[[]m.Member] {
	rows, err := r.db.QueryContext(context.Background(), getAllMembersQuery)
	if err != nil {
		return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
	}
	defer rows.Close()

	var members []m.Member
	for rows.Next() {
		var resultRow GetAllMembersQueryResult
		err := rows.Scan(
			&resultRow.MemberID,
			&resultRow.FirstName,
			&resultRow.LastName,
			&resultRow.DateOfBirth,
			&resultRow.MembershipNumber,
			&resultRow.Season,
			&resultRow.SeasonStartsAt,
			&resultRow.SeasonEndsAt,
			&resultRow.ExclusionDeliberatedAt,
			&resultRow.AmountPaid,
			&resultRow.PaidAt,
			&resultRow.Currency,
		)
		if err != nil {
			return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
		}

		memberResult := MapToMemberFromAllMembersQuery(resultRow)
		if !memberResult.IsSuccess() {
			return result.Err[[]m.Member](memberResult.Error())
		}

		members = append(members, memberResult.Value())
	}

	if err = rows.Err(); err != nil {
		return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
	}

	return result.Ok(members)
}

func (r *SQLMemberRepository) GetMembersBySeason(year int64) result.Result[[]m.Member] {
	var members []m.Member
	rows, err := r.db.QueryContext(context.Background(), getMembersBySeasonQuery, year)
	if err != nil {
		return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
	}
	defer rows.Close()

	for rows.Next() {
		var resultRow GetMembersBySeasonQueryResult
		err := rows.Scan(
			&resultRow.MemberID,
			&resultRow.FirstName,
			&resultRow.LastName,
			&resultRow.DateOfBirth,
			&resultRow.MembershipNumber,
			&resultRow.SeasonStartsAt,
			&resultRow.SeasonEndsAt,
			&resultRow.ExclusionDeliberatedAt,
			&resultRow.AmountPaid,
			&resultRow.PaidAt,
			&resultRow.Currency,
		)
		if err != nil {
			return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
		}

		memberResult := MapToMemberFromQueryBySeason(resultRow)
		if !memberResult.IsSuccess() {
			return result.Err[[]m.Member](memberResult.Error())
		}

		members = append(members, memberResult.Value())
	}

	if err = rows.Err(); err != nil {
		return result.Err[[]m.Member](errors.RepositoryError{Description: err.Error()})
	}

	return result.Ok(members)
}

func (r *SQLMemberRepository) GetMemberById(id domain.Id[m.Member]) result.Result[m.MemberDetails] {
	var resultRow GetMemberByIdQueryResult
	err := r.db.QueryRowContext(context.Background(), getMemberByIdQuery, id.Value).Scan(
		&resultRow.MemberID,
		&resultRow.FirstName,
		&resultRow.LastName,
		&resultRow.DateOfBirth,
		&resultRow.Email,
		&resultRow.PhoneNumbers,
		&resultRow.Addresses,
		&resultRow.Memberships,
	)

	return result.MapErr(result.Bind(result.From(true, err), func(_ bool) result.Result[m.MemberDetails] {
		return MapToMemberFromMemberByIdQuery(resultRow)
	}), func(err error) error {
		if err == sql.ErrNoRows {
			return errors.NotFoundError{Description: "Member not found"}
		}
		return errors.RepositoryError{Description: err.Error()}
	})
}

func (r *SQLMemberRepository) GetMembersWhoDidNotPayForServices() []m.Member {
	// TODO: Implement query for members with unpaid service rentals
	return []m.Member{}
}

func (r *SQLMemberRepository) GetMembersWhoDidNotPayForMembership() []m.Member {
	// TODO: Implement query for members with unpaid memberships
	return []m.Member{}
}
