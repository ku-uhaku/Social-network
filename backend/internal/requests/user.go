package requests

import (
	"errors"

	"kuu/internal/models"
)

// ValidateUpdateProfile ensures the privacy flag is structurally sound
func ValidateUpdateProfile(payload models.UpdateProfilePayload) []error {
	var errs []error

	if payload.IsPublic != 0 && payload.IsPublic != 1 {
		errs = append(errs, errors.New("is_public must be either 0 (private) or 1 (public)"))
	}
	// _,err:=helper.IsValidImage([]byte(p.Content))
	// if err!=nil{
	// 	errs = append(errs, errors.New("the comment image not good check the format or content"))
	// }

	return errs
}

func ValidateFollowAction(payload models.FollowActionPayload) []error {
	var errs []error
	if payload.TargetUserID <= 0 {
		errs = append(errs, errors.New("target_user_id is required and must be valid"))
	}
	return errs
}