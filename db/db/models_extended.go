package db

import (
	"time"

	"github.com/rs/xid"
)

type Profile struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" example:"2024-04-06T06:20:10.749615Z"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:"2024-04-06T06:20:10.749615Z"`
	Email     string     `json:"email,omitempty"     example:"johndoe@me.com"`
	Username  string     `json:"username,omitempty"  example:"johndoe"`
	Firstname string     `json:"firstname,omitempty" example:"John"`
	Lastname  string     `json:"lastname,omitempty"  example:"Doe"`
	ID        []byte     `json:"id,omitempty"        example:"ZhzG-Vbn1dT42OPB"`
	Verified  bool       `json:"verified"            example:"false"`
	Blocked   bool       `json:"blocked"             example:"false"`
} // @name Profile

func (user *User) Profile() (*Profile, error) {
	userID, err := xid.FromBytes(user.ID)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:        userID.Bytes(),
		Email:     user.Email,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Verified:  user.Verified,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
		Blocked:   user.Blocked,
	}, nil
}
