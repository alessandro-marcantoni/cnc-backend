package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

//go:embed queries/insert_payment.sql
var insertPaymentQuery string

//go:embed queries/update_payment.sql
var updatePaymentQuery string

//go:embed queries/delete_payment.sql
var deletePaymentQuery string

type SQLPaymentRepository struct {
	db *sql.DB
}

func NewSQLPaymentRepository(db *sql.DB) *SQLPaymentRepository {
	return &SQLPaymentRepository{db: db}
}

func (r *SQLPaymentRepository) CreatePaymentForMembershipPeriod(membershipPeriodId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64] {
	var paymentId int64
	paidAt := time.Now()

	err := r.db.QueryRowContext(
		context.Background(),
		insertPaymentQuery,
		nil, // rented_facility_id
		membershipPeriodId,
		amount,
		currency,
		paidAt,
		paymentMethod,
		transactionRef,
	).Scan(&paymentId)

	if err != nil {
		return result.Err[int64](errors.RepositoryError{Description: err.Error()})
	}

	return result.Ok(paymentId)
}

func (r *SQLPaymentRepository) CreatePaymentForRentedFacility(rentedFacilityId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64] {
	var paymentId int64
	paidAt := time.Now()

	err := r.db.QueryRowContext(
		context.Background(),
		insertPaymentQuery,
		rentedFacilityId,
		nil, // membership_period_id
		amount,
		currency,
		paidAt,
		paymentMethod,
		transactionRef,
	).Scan(&paymentId)

	if err != nil {
		return result.Err[int64](errors.RepositoryError{Description: err.Error()})
	}

	return result.Ok(paymentId)
}

func (r *SQLPaymentRepository) UpdatePayment(paymentId domain.Id[payment.Payment], amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[bool] {
	execResult, err := r.db.ExecContext(
		context.Background(),
		updatePaymentQuery,
		amount,
		currency,
		paymentMethod,
		transactionRef,
		paymentId.Value,
	)

	if err != nil {
		return result.Err[bool](errors.RepositoryError{Description: err.Error()})
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return result.Err[bool](errors.RepositoryError{Description: err.Error()})
	}

	if rowsAffected == 0 {
		return result.Err[bool](errors.RepositoryError{Description: "payment not found"})
	}

	return result.Ok(true)
}

func (r *SQLPaymentRepository) DeletePayment(paymentId domain.Id[payment.Payment]) result.Result[bool] {
	execResult, err := r.db.ExecContext(
		context.Background(),
		deletePaymentQuery,
		paymentId.Value,
	)

	if err != nil {
		return result.Err[bool](errors.RepositoryError{Description: err.Error()})
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return result.Err[bool](errors.RepositoryError{Description: err.Error()})
	}

	if rowsAffected == 0 {
		return result.Err[bool](errors.RepositoryError{Description: "payment not found"})
	}

	return result.Ok(true)
}
