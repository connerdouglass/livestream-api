package services

import (
	"time"

	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

// MembershipService manages memberships between accounts and creator profiles
type MembershipService struct {
	DB *gorm.DB
}

// GetMembers gets the members of a creator profile
func (s *MembershipService) GetMembers(creatorID uint64) ([]*models.CreatorProfileMember, error) {
	var members []*models.CreatorProfileMember
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("creator_profile_id = ?", creatorID).
		Preload("Account").
		Preload("CreatorProfile").
		Find(&members).
		Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetMembershipsForAccount gets the memberships for an account
func (s *MembershipService) GetMembershipsForAccount(accountID uint64) ([]*models.CreatorProfileMember, error) {
	var members []*models.CreatorProfileMember
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("account_id = ?", accountID).
		Preload("Account").
		Preload("CreatorProfile").
		Find(&members).
		Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetCreatorProfiles gets the creator profiles that are accessible to the given account
func (s *MembershipService) GetCreatorProfiles(accountID uint64) ([]*models.CreatorProfile, error) {
	members, err := s.GetMembershipsForAccount(accountID)
	if err != nil {
		return nil, err
	}
	profiles := make([]*models.CreatorProfile, len(members))
	for i := range members {
		profiles[i] = members[i].CreatorProfile
	}
	return profiles, nil
}

// IsAccountMember checks if an account is a member of a creator profile
func (s *MembershipService) IsMember(creatorID, accountID uint64) (bool, error) {
	var count int64
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("creator_profile_id = ?", creatorID).
		Where("account_id = ?", accountID).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return (count > 0), nil
}

// AddMember adds an account as a member to a given creator profile
func (s *MembershipService) AddMember(creatorID, accountID uint64) error {

	// If the user is already a member
	isMember, err := s.IsMember(creatorID, accountID)
	if err != nil {
		return err
	}
	if isMember {
		return nil
	}

	// Create the membership
	member := models.CreatorProfileMember{
		CreatorProfileID: creatorID,
		AccountID:        accountID,
		CreatedDate:      time.Now(),
	}
	return s.DB.Create(&member).Error

}
