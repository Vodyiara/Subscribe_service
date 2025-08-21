package repository

import (
	"Test_project_Effective_Mobile"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TodoSub interface {
	Create(item Test_project_Effective_Mobile.Subscripe) (int, error)
	GetAllSubs() ([]Test_project_Effective_Mobile.Subscripe, error)
	Update(subId int, item Test_project_Effective_Mobile.UpdateSubInput) error
	Delete(subId int) error
	GetByIdSub(subId int) (Test_project_Effective_Mobile.Subscripe, error)
	GetTotalPrice(user_id uuid.UUID, serice_name, statr_date, end_dat string) (int, error)
}

type Repository struct {
	TodoSub
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoSub: NewTodoSubscription(db),
	}
}
