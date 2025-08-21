package service

import (
	"Test_project_Effective_Mobile"
	"Test_project_Effective_Mobile/pkg/repository"
	"github.com/google/uuid"
)

type SubService struct {
	repo repository.TodoSub
}

func NewSubService(repo repository.TodoSub) *SubService {
	return &SubService{repo: repo}
}

func (s *SubService) Create(item Test_project_Effective_Mobile.Subscripe) (int, error) {
	return s.repo.Create(item)
}

func (s *SubService) GetAllSubs() ([]Test_project_Effective_Mobile.Subscripe, error) {
	return s.repo.GetAllSubs()
}

func (s *SubService) Update(subId int, item Test_project_Effective_Mobile.UpdateSubInput) error {
	return s.repo.Update(subId, item)
}

func (s *SubService) Delete(subId int) error {
	return s.repo.Delete(subId)
}

func (s *SubService) GetByIdSub(subId int) (Test_project_Effective_Mobile.Subscripe, error) {
	return s.repo.GetByIdSub(subId)
}

func (s *SubService) GetTotalPrice(user_id uuid.UUID, serice_name, statr_date, end_dat string) (int, error) {
	return s.repo.GetTotalPrice(user_id, serice_name, statr_date, end_dat)
}
