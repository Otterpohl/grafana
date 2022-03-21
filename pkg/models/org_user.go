package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Typed errors
var (
	ErrLastOrgAdmin        = errors.New("cannot remove last organization admin")
	ErrOrgUserNotFound     = errors.New("cannot find the organization user")
	ErrOrgUserAlreadyAdded = errors.New("user is already added to organization")
)

type RoleType string

const (
	ROLE_EMPTY  RoleType = ""
	ROLE_VIEWER RoleType = "Viewer"
	ROLE_EDITOR RoleType = "Editor"
	ROLE_ADMIN  RoleType = "Admin"
)

func (r RoleType) IsValid() bool {
	return r == ROLE_VIEWER || r == ROLE_ADMIN || r == ROLE_EDITOR || r == ROLE_EMPTY
}

func (r RoleType) Includes(other RoleType) bool {
	switch r {
	case ROLE_ADMIN:
		return true
	case ROLE_EDITOR:
		return other != ROLE_ADMIN
	case ROLE_VIEWER:
		return other == ROLE_EMPTY || r == other
	default:
		return r == other
	}
}

func (r RoleType) Children() []RoleType {
	switch r {
	case ROLE_ADMIN:
		return []RoleType{ROLE_EDITOR, ROLE_VIEWER, ROLE_EMPTY}
	case ROLE_EDITOR:
		return []RoleType{ROLE_VIEWER, ROLE_EMPTY}
	case ROLE_VIEWER:
		return []RoleType{ROLE_EMPTY}
	default:
		return nil
	}
}

func (r *RoleType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*r = RoleType(str)

	if !r.IsValid() {
		// TODO Restore this if we end up having ROLE_EMPTY != ""
		// if (*r) != "" {
		return fmt.Errorf("JSON validation error: invalid role value: %s", *r)
		// }

		// *r = ROLE_VIEWER
	}

	return nil
}

type OrgUser struct {
	Id      int64
	OrgId   int64
	UserId  int64
	Role    RoleType
	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type RemoveOrgUserCommand struct {
	UserId                   int64
	OrgId                    int64
	ShouldDeleteOrphanedUser bool
	UserWasDeleted           bool
}

type AddOrgUserCommand struct {
	LoginOrEmail string   `json:"loginOrEmail" binding:"Required"`
	Role         RoleType `json:"role" binding:"Required"`

	OrgId  int64 `json:"-"`
	UserId int64 `json:"-"`
}

type UpdateOrgUserCommand struct {
	Role RoleType `json:"role" binding:"Required"`

	OrgId  int64 `json:"-"`
	UserId int64 `json:"-"`
}

// ----------------------
// QUERIES

type GetOrgUsersQuery struct {
	UserID int64
	OrgId  int64
	Query  string
	Limit  int

	User   *SignedInUser
	Result []*OrgUserDTO
}

type SearchOrgUsersQuery struct {
	OrgID int64
	Query string
	Page  int
	Limit int

	User   *SignedInUser
	Result SearchOrgUsersQueryResult
}

type SearchOrgUsersQueryResult struct {
	TotalCount int64         `json:"totalCount"`
	OrgUsers   []*OrgUserDTO `json:"OrgUsers"`
	Page       int           `json:"page"`
	PerPage    int           `json:"perPage"`
}

// ----------------------
// Projections and DTOs

type OrgUserDTO struct {
	OrgId         int64           `json:"orgId"`
	UserId        int64           `json:"userId"`
	Email         string          `json:"email"`
	Name          string          `json:"name"`
	AvatarUrl     string          `json:"avatarUrl"`
	Login         string          `json:"login"`
	Role          string          `json:"role"`
	LastSeenAt    time.Time       `json:"lastSeenAt"`
	Updated       time.Time       `json:"-"`
	Created       time.Time       `json:"-"`
	LastSeenAtAge string          `json:"lastSeenAtAge"`
	AccessControl map[string]bool `json:"accessControl,omitempty"`
}
