package dto

import "user-service/internal/domain"

func ToUserStatusResponse(user *domain.User) UserStatusResponse {
	return UserStatusResponse{
		ID:         user.ID,
		Username:   user.Username,
		Balance:    user.Balance,
		ReferrerID: user.ReferrerID,
	}
}

func ToLeaderboardUserDTO(user *domain.User, rank int) LeaderboardUserDTO {
	return LeaderboardUserDTO{
		ID:       user.ID,
		Username: user.Username,
		Balance:  user.Balance,
		Rank:     rank,
	}
}
