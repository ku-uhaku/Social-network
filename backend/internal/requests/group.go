package requests

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"kuu/internal/models"
)

func ValidateCreateGroup(payload models.CreateGroupPayload) []error {
	var errs []error
	if strings.TrimSpace(payload.Title) == "" {
		errs = append(errs, errors.New("title is required"))
		return errs
	} else if titleLen := utf8.RuneCountInString(strings.TrimSpace(payload.Title)); titleLen < 3 || titleLen > 20 {
		errs = append(errs, errors.New("title must be between 3 and 20 characters"))
		return errs
	}

	if strings.TrimSpace(payload.Description) == "" {
		errs = append(errs, errors.New("description is required"))
		return errs
	} else if utf8.RuneCountInString(strings.TrimSpace(payload.Description)) > 200 {
		errs = append(errs, errors.New("description cannot be longer than 200 characters"))
		return errs
	}

	return errs
}

func ValidateUpdateGroup(payload models.UpdateGroupPayload) []error {
	var errs []error

	if strings.TrimSpace(payload.Title) == "" {
		errs = append(errs, errors.New("title cannot be empty"))
	}

	if strings.TrimSpace(payload.Description) == "" {
		errs = append(errs, errors.New("description cannot be empty"))
	}

	return errs
}

func ValidateInviteMembers(payload models.InviteMembersPayload) []error {
	var errs []error
	if payload.GroupID <= 0 {
		errs = append(errs, errors.New("group_id is required"))
	}
	if len(payload.TargetUserIDs) == 0 {
		errs = append(errs, errors.New("target_user_ids must contain at least one user ID"))
	}
	return errs
}

func ValidateJoinRequestAction(payload models.GroupActionPayload) []error {
	var errs []error
	if payload.GroupID <= 0 {
		errs = append(errs, errors.New("group_id is required"))
	}
	if payload.TargetUserID <= 0 {
		errs = append(errs, errors.New("target_user_id is required"))
	}
	return errs
}

func ValidateCreateGroupEvent(payload models.CreateGroupEventPayload) []error {
	var errs []error
	if payload.GroupID <= 0 {
		errs = append(errs, errors.New("group_id is required"))
		return errs
	}
	if strings.TrimSpace(payload.Title) == "" {
		errs = append(errs, errors.New("title is required"))
	}
	if utf8.RuneCountInString(strings.TrimSpace(payload.Title)) > 20 {
				errs = append(errs, errors.New("title moore long "))
				return  errs
	}
	if strings.TrimSpace(payload.Description) == "" {
		errs = append(errs, errors.New("description is required"))
	}
	if utf8.RuneCountInString(strings.TrimSpace(payload.Description)) > 100 {
				errs = append(errs, errors.New("description moore long "))
				return  errs
	}
	if payload.EventTime.IsZero() {
		errs = append(errs, errors.New("event_time is required"))
		return errs
	} else if payload.EventTime.Before(time.Now()) {
		errs = append(errs, errors.New("event_time must be in the future"))
		return errs
	}
	return errs
}

func ValidateEventResponse(payload models.EventResponsePayload) []error {
	var errs []error
	if payload.EventID <= 0 {
		errs = append(errs, errors.New("event_id is required"))
	}
	if payload.Status != "going" && payload.Status != "not_going" {
		errs = append(errs, errors.New("status must be 'going' or 'not_going'"))
	}
	return errs
}