package service

import (
	"Test_project_Effective_Mobile"
	"Test_project_Effective_Mobile/pkg/repository"
	"github.com/google/uuid"
)

type TodoSub interface {
	Create(item Test_project_Effective_Mobile.Subscripe) (int, error)
	GetAllSubs() ([]Test_project_Effective_Mobile.Subscripe, error)
	Update(subId int, item Test_project_Effective_Mobile.UpdateSubInput) error
	Delete(subId int) error
	GetByIdSub(subId int) (Test_project_Effective_Mobile.Subscripe, error)
	GetTotalPrice(user_id uuid.UUID, serice_name, statr_date, end_dat string) (int, error)
}

type Service struct {
	TodoSub
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		TodoSub: NewSubService(repos.TodoSub),
	}
}
