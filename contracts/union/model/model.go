package model

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectID = primitive.ObjectID

func MarshalObjectID(id ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprintf(w, "\"%s\"", id.Hex())
	})
}

func UnmarshalObjectID(v interface{}) (ObjectID, error) {
	str, ok := v.(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ObjectID must be a string")
	}
	return primitive.ObjectIDFromHex(str)
}

func MarshalObjectIDScalar(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprintf(w, "\"%s\"", id.Hex())
	})
}

// Union represents the main union structure
type Union struct {
	ID                   primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name                 string               `json:"name,omitempty" bson:"name,omitempty"`
	UnionID              string               `json:"unionID,omitempty" bson:"unionID,omitempty"`
	Status               int                  `json:"status,omitempty" bson:"status,omitempty"`
	CreatedOn            time.Time            `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	Information          UnionInfo            `json:"information,omitempty" bson:"information,omitempty"`
	Modules              []*string            `json:"modules,omitempty" bson:"modules,omitempty"`
	Deleted              bool                 `json:"deleted,omitempty" bson:"deleted"`
	BargainingUnits      []*string            `json:"bargainingUnits,omitempty" bson:"bargainingUnits,omitempty"`
	BannerURL            string               `json:"bannerURL,omitempty" bson:"bannerURL,omitempty"`
	DefaultPermissions   map[string]int64     `json:"defaultPermissions,omitempty" bson:"defaultPermissions,omitempty"`
	AdminPermissions     map[string]int64     `json:"adminPermissions,omitempty" bson:"adminPermissions,omitempty"`
	DefaultUser          DefaultUserInfo      `json:"defaultUser,omitempty" bson:"defaultUser,omitempty"`
	AccountManagerID     []primitive.ObjectID `json:"accountManager,omitempty" bson:"accountManager,omitempty"`
	CommunicationRepID   []primitive.ObjectID `json:"communicationRep,omitempty" bson:"communicationRep,omitempty"`
	AccountManager       []*Manager
	CommunicationRep     []*Manager
	CallDropNumber       string        `json:"callDropNumber,omitempty" bson:"callDropNumber,omitempty"`
	Domain               string        `json:"domain,omitempty" bson:"domain,omitempty"`
	BannedDomains        []*string     `json:"bannedDomains,omitempty" bson:"bannedDomains,omitempty"`
	Theme                string        `json:"theme,omitempty" bson:"theme,omitempty"`
	Twitter              string        `json:"twitter,omitempty" bson:"twitter,omitempty"`
	TwitterLinks         []*string     `json:"twitterLinks,omitempty" bson:"twitterLinks,omitempty"`
	Facebook             string        `json:"facebook,omitempty" bson:"facebook,omitempty"`
	FacebookLinks        []*string     `json:"facebookLinks,omitempty" bson:"facebookLinks,omitempty"`
	Instagram            string        `json:"instagram,omitempty" bson:"instagram,omitempty"`
	InstagramLinks       []*string     `json:"instagramLinks,omitempty" bson:"instagramLinks,omitempty"`
	FirstUser            FirstUserInfo `json:"firstUser,omitempty" bson:"firstUser,omitempty"`
	ThemeImage           string        `json:"themeImage,omitempty" bson:"themeImage,omitempty"`
	ZoomID               string        `json:"zoomID,omitempty" bson:"zoomID,omitempty"`
	HostEmail            *bool         `json:"hostEmail,omitempty" bson:"hostEmail,omitempty"`
	DefaultEmailPassword string        `json:"defaultEmailPassword,omitempty" bson:"defaultEmailPassword,omitempty"`
	DeletedAt            time.Time     `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type UnionsResponse struct {
	Unions []*Union `json:"unions"`
	Count  int      `json:"count"`
}

// UnionInfo contains detailed information about the union
type UnionInfo struct {
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	UnionMail        string `json:"unionMail,omitempty" bson:"unionMail,omitempty"`
	ImageURL         string `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	LandingPage      string `json:"landingPage,omitempty" bson:"landingPage,omitempty"`
	Address          string `json:"address,omitempty" bson:"address,omitempty"`
	City             string `json:"city,omitempty" bson:"city,omitempty"`
	Country          string `json:"country,omitempty" bson:"country,omitempty"`
	State            string `json:"state,omitempty" bson:"state,omitempty"`
	Province         string `json:"province,omitempty" bson:"province,omitempty"`
	PostalCode       string `json:"postalCode,omitempty" bson:"postalCode,omitempty"`
	ZipCode          string `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
	Phone            string `json:"phone,omitempty" bson:"phone,omitempty"`
	Mobile           string `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Description      string `json:"description,omitempty" bson:"description,omitempty"`
	BannerURL        string `json:"bannerURL,omitempty" bson:"bannerURL,omitempty"`
	Fax              string `json:"fax,omitempty" bson:"fax,omitempty"`
	PresidentMessage string `json:"president_message,omitempty" bson:"president_message,omitempty"`
}

// DefaultUserInfo represents the default user information for a union
type DefaultUserInfo struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Level    int    `json:"level,omitempty"`
}

type Manager struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName  string             `json:"firstName,omitempty" bson:"firstname,omitempty"`
	LastName   string             `json:"lastName,omitempty" bson:"lastname,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone      string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Department string             `json:"department,omitempty" bson:"department,omitempty"`
	ImageURL   string             `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	Mobile     string             `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

// FirstUserInfo contains information about the first user of the union
type FirstUserInfo struct {
	FirstName   string    `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Password    string    `json:"password,omitempty" bson:"password,omitempty"`
	Email       string    `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Position    string    `json:"position,omitempty" bson:"position,omitempty"`
	DateOfBirth time.Time `json:"dob,omitempty" bson:"dob,omitempty"`
}
type RegisterInput struct {
	Union       Union           `json:"union"`
	User        FirstUserInfo   `json:"user"`
	DefaultUser DefaultUserInfo `json:"defaultUser"`
}
