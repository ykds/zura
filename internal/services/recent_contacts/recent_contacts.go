package recent_contacts

import "zura/internal/entity"

func NewRecentContactsService(recentContactsEntity entity.RecentContactsEntity) RecentContactsService {
	return &recentContactsService{
		recentContactsEntity: recentContactsEntity,
	}
}

type RecentContactsService interface {
}

type recentContactsService struct {
	recentContactsEntity entity.RecentContactsEntity
}
