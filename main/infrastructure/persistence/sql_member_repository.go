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

//go:embed queries/insert_member.sql
var insertMemberQuery string

//go:embed queries/insert_phone_number.sql
var insertPhoneNumberQuery string

//go:embed queries/insert_address.sql
var insertAddressQuery string

//go:embed queries/get_next_membership_number.sql
var getNextMembershipNumberQuery string

//go:embed queries/insert_membership.sql
var insertMembershipQuery string

//go:embed queries/insert_membership_period.sql
var insertMembershipPeriodQuery string

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
			&resultRow.MembershipStatus,
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
			&resultRow.MembershipStatus,
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

func (r *SQLMemberRepository) GetMemberById(id domain.Id[m.Member], season string) result.Result[m.MemberDetails] {
	var resultRow GetMemberByIdQueryResult
	err := r.db.QueryRowContext(context.Background(), getMemberByIdQuery, id.Value, season).Scan(
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

func (r *SQLMemberRepository) CreateMember(user m.User, createMembership bool, seasonId *int64, price *float64) result.Result[m.MemberDetails] {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	// 1. Insert member
	var memberId int64
	err = tx.QueryRowContext(ctx, insertMemberQuery,
		user.FirstName,
		user.LastName,
		user.BirthDate,
		user.Email.Value,
	).Scan(&memberId)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert member: " + err.Error()})
	}

	// 2. Insert phone numbers
	for _, phone := range user.PhoneNumbers {
		prefix := ""
		if phone.Prefix != nil {
			prefix = *phone.Prefix
		}
		fullNumber := prefix + phone.Number
		_, err = tx.ExecContext(ctx, insertPhoneNumberQuery, memberId, fullNumber, nil)
		if err != nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert phone number: " + err.Error()})
		}
	}

	// 3. Insert addresses
	for _, address := range user.Addresses {
		_, err = tx.ExecContext(ctx, insertAddressQuery,
			memberId,
			address.Country,
			address.City,
			address.Street,
			address.Number,
		)
		if err != nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert address: " + err.Error()})
		}
	}

	// Get next membership number
	var nextMembershipNumber int64
	err = tx.QueryRowContext(ctx, getNextMembershipNumberQuery).Scan(&nextMembershipNumber)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to get next membership number: " + err.Error()})
	}

	// Insert membership
	var membershipId int64
	err = tx.QueryRowContext(ctx, insertMembershipQuery,
		memberId,
		nextMembershipNumber,
	).Scan(&membershipId)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert membership: " + err.Error()})
	}

	// 4. Create membership if requested
	if createMembership {
		// Validate that seasonId is provided
		if seasonId == nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "seasonId is required when createMembership is true"})
		}

		// Get season dates
		var seasonStartsAt, seasonEndsAt string
		err = tx.QueryRowContext(ctx, "SELECT starts_at, ends_at FROM seasons WHERE id = $1", *seasonId).
			Scan(&seasonStartsAt, &seasonEndsAt)
		if err != nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to get season dates: " + err.Error()})
		}

		// Determine the price to use (custom price or suggested price)
		membershipPrice := m.SuggestedMembershipPrice
		if price != nil {
			membershipPrice = *price
		}

		// Insert membership period with status_id = 1 (ACTIVE)
		_, err = tx.ExecContext(ctx, insertMembershipPeriodQuery,
			membershipId,
			seasonStartsAt,
			seasonEndsAt,
			1, // status_id for ACTIVE
			*seasonId,
			membershipPrice,
		)
		if err != nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert membership period: " + err.Error()})
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	// Fetch and return the created member details
	if createMembership && seasonId != nil {
		// We need to get the season code to query the member
		var seasonCode string
		err = r.db.QueryRowContext(ctx, "SELECT code FROM seasons WHERE id = $1", *seasonId).Scan(&seasonCode)
		if err != nil {
			return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to get season code: " + err.Error()})
		}
		return r.GetMemberById(domain.Id[m.Member]{Value: memberId}, seasonCode)
	}

	// If no membership was created, return member details with empty memberships
	return result.Ok(m.MemberDetails{
		User: m.User{
			Id:           domain.Id[m.User]{Value: memberId},
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			BirthDate:    user.BirthDate,
			Email:        user.Email,
			Addresses:    user.Addresses,
			PhoneNumbers: user.PhoneNumbers,
		},
		Memberships: []m.Membership{},
	})
}

func (r *SQLMemberRepository) AddMembership(memberId domain.Id[m.Member], seasonId int64, seasonStartsAt string, seasonEndsAt string, price float64) result.Result[m.MemberDetails] {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	// Get the membership_id for this member
	var membershipId int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM memberships WHERE member_id = $1", memberId.Value).Scan(&membershipId)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to get membership id: " + err.Error()})
	}

	// Insert membership period with status_id = 1 (ACTIVE)
	_, err = tx.ExecContext(ctx, insertMembershipPeriodQuery,
		membershipId,
		seasonStartsAt,
		seasonEndsAt,
		1, // status_id for ACTIVE
		seasonId,
		price,
	)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to insert membership period: " + err.Error()})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	// Fetch season code to query the member
	var seasonCode string
	err = r.db.QueryRowContext(ctx, "SELECT code FROM seasons WHERE id = $1", seasonId).Scan(&seasonCode)
	if err != nil {
		return result.Err[m.MemberDetails](errors.RepositoryError{Description: "failed to get season code: " + err.Error()})
	}

	// Fetch and return the updated member details
	return r.GetMemberById(memberId, seasonCode)
}
