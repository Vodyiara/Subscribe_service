package Test_project_Effective_Mobile

import (
	"errors"
	"github.com/google/uuid"
)

type Subscripe struct {
	Id           int       `json:"id" db:"id"`
	Service_name string    `json:"service_name" db:"service_name" binding:"required"`
	Price        int       `json:"price" db:"price" binding:"required"`
	User_id      uuid.UUID `json:"user_id" db:"user_id" binding:"required"`
	Start_date   string    `json:"start_date" db:"start_date" binding:"required"`
	End_date     *string   `json:"end_date" db:"end_date"`
}

type UpdateSubInput struct {
	Service_name *string `json:"service_name"`
	Price        *int    `json:"price"`
	Start_date   *string `json:"start_date"`
	End_date     *string `json:"end_date"`
	Done         *bool   `json:"done" `
}

func (i UpdateSubInput) Valid() error {
	if i.Service_name == nil && i.Price == nil && i.Start_date == nil && i.End_date == nil && i.Done == nil {
		return errors.New("update structur has not value")
	}
	return nil
}
