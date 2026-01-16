package payment

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type PaymentManagementService struct {
	repository PaymentRepository
}

func NewPaymentManagementService(repository PaymentRepository) *PaymentManagementService {
	return &PaymentManagementService{repository: repository}
}

func (this PaymentManagementService) CreatePaymentForMembershipPeriod(membershipPeriodId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64] {
	return this.repository.CreatePaymentForMembershipPeriod(membershipPeriodId, amount, currency, paymentMethod, transactionRef)
}

func (this PaymentManagementService) CreatePaymentForRentedFacility(rentedFacilityId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64] {
	return this.repository.CreatePaymentForRentedFacility(rentedFacilityId, amount, currency, paymentMethod, transactionRef)
}

func (this PaymentManagementService) UpdatePayment(paymentId domain.Id[Payment], amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[bool] {
	return this.repository.UpdatePayment(paymentId, amount, currency, paymentMethod, transactionRef)
}

func (this PaymentManagementService) DeletePayment(paymentId domain.Id[Payment]) result.Result[bool] {
	return this.repository.DeletePayment(paymentId)
}
