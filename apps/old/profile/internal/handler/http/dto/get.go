package httphandlerdto

import "mandacode.com/accounts/profile/internal/domain/model"

type GetProfileResponse struct {
	*model.Profile
}
