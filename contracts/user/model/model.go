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

type User struct {
	ID                      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UnionID                 primitive.ObjectID `json:"unionID,omitempty" bson:"unionID,omitempty"`
	EmployeeID              string             `json:"employeeID,omitempty" bson:"memberID"`
	Username                string             `json:"username,omitempty" bson:"username"`
	Password                string             `json:"password,omitempty" bson:"password"`
	FirstName               string             `json:"firstName,omitempty" bson:"firstName"`
	LastName                string             `json:"lastName,omitempty" bson:"lastName"`
	MiddleName              string             `json:"middleName,omitempty" bson:"middleName,omitempty"`
	MaidenName              string             `json:"maidenName,omitempty" bson:"maidenName,omitempty"`
	CommonName              string             `json:"commonName,omitempty" bson:"commonName,omitempty"`
	Gender                  string             `json:"gender,omitempty" bson:"gender,omitempty"`
	Profile                 UserInfo           `json:"profile,omitempty" bson:"profile"`
	CreatedOn               time.Time          `json:"createdOn,omitempty" bson:"createdOn"`
	Deleted                 bool               `json:"deleted,omitempty" bson:"deleted"`
	DeletedAt               time.Time          `json:"deletedAT,omitempty" bson:"deletedAT"`
	LoggedIn                bool               `json:"loggedIn,omitempty" bson:"loggedIn"`
	Status                  string             `json:"status,omitempty" bson:"status"`
	Token                   string             `json:"token,omitempty" bson:"token"`
	DateOfBirth             time.Time          `json:"dateOfBirth,omitempty" bson:"dateOfBirth"`
	StartDate               time.Time          `json:"startDate,omitempty" bson:"startDate"`
	Location                string             `json:"location,omitempty" bson:"location"`
	UnionPosition           string             `json:"unionPosition,omitempty" bson:"unionPosition"`
	Unit                    string             `json:"unit,omitempty" bson:"unit"`
	JobTitle                string             `json:"jobTitle,omitempty" bson:"jobTitle"`
	Commitee                string             `json:"commitee,omitempty" bson:"commitee"`
	MembershipType          string             `json:"membershipType,omitempty" bson:"membershipType"`
	EmploymentType          string             `json:"employmentType,omitempty" bson:"employmentType"`
	EmploymentStatus        string             `json:"employmentStatus,omitempty" bson:"employmentStatus"`
	Permission              map[string]int64   `json:"permission,omitempty" bson:"permission"`
	MeritPoint              int                `json:"meritPoint,omitempty" bson:"meritPoint"`
	DemeritPoint            int                `json:"demeritPoint,omitempty" bson:"demeritPoint"`
	LastLoginDate           time.Time          `json:"lastLoginDate,omitempty" bson:"lastLoginDate"`
	PasswordResetKey        string             `json:"passwordResetKey,omitempty" bson:"passwordResetKey"`
	PasswordResetExpireTime time.Time          `json:"passwordResetExpireTime,omitempty" bson:"passwordResetExpireTime"`
	Tokens                  []string           `json:"tokens,omitempty" bson:"tokens"`
	IsAdmin                 bool               `json:"isAdmin,omitempty" bson:"isAdmin"`
	Signature               string             `json:"signature,omitempty" bson:"signature"`
	Device                  string             `json:"device,omitempty" bson:"device"`
	SeniorityNumber         string             `json:"seniorityNumber,omitempty" bson:"seniorityNumber"`
	PreferredLanguage       string             `json:"preferredLanguage,omitempty" bson:"preferredLanguage,omitempty"`
	EmailPassword           string             `json:"emailPassword,omitempty" bson:"emailPassword,omitempty"`
	// FamilyMembersData       []*FamilyMembers   `json:"familyMembersData,omitempty" bson:"familyMembersData,omitempty"`
	// JobLocation             []*UserLocation    `json:"jobLocation,omitempty" bson:"jobLocation,omitempty"`
	// Courses                 []*Courses         `json:"courses,omitempty" bson:"courses,omitempty"`
	Department     string   `json:"department,omitempty" bson:"department,omitempty"`
	BadgeNumber    string   `json:"badgeNumber,omitempty" bson:"badgeNumber,omitempty"`
	Classification string   `json:"classification,omitempty" bson:"classification,omitempty"`
	Zone           string   `json:"zone,omitempty" bson:"zone,omitempty"`
	Shift          string   `json:"shift,omitempty" bson:"shift,omitempty"`
	Level          int      `json:"level,omitempty" bson:"level,omitempty"`
	MemberCraft    string   `json:"member_craft,omitempty" bson:"member_craft,omitempty"`
	MemberClass    string   `json:"member_class,omitempty" bson:"member_class,omitempty"`
	Teachables     []string `json:"teachables,omitempty" bson:"teachables,omitempty"`
	// DriversLicense          DriversLicense     `json:"driversLicense,omitempty" bson:"driversLicense,omitempty"`
	CallOpOut         bool   `json:"callOpOut,omitempty" bson:"callOpOut"`
	EmailOpOut        bool   `json:"emailOpOut,omitempty" bson:"emailOpOut"`
	TextOpOut         bool   `json:"textOpOut,omitempty" bson:"textOpOut"`
	PushOpOut         bool   `json:"pushOpOut,omitempty" bson:"pushOpOut"`
	RegEmailOpOut     bool   `json:"regEmailOpOut,omitempty" bson:"regEmailOpOut"`
	UnionStatus       string `json:"unionStatus,omitempty" bson:"unionStatus,omitempty"`
	IndigenousStatus  string `json:"indigenousStatus,omitempty" bson:"indigenousStatus,omitempty"`
	DepartmentComplex string `json:"departmentComplex,omitempty" bson:"departmentComplex,omitempty"`
	Notes             string `json:"notes,omitempty" bson:"notes,omitempty"`
	StrikePay         bool   `json:"strikePay,omitempty" bson:"strikePay,omitempty"`
	ResetRequired     bool   `json:"resetRequired,omitempty" bson:"resetRequired,omitempty"`
	ZoomID            string `json:"zoomID,omitempty" bson:"zoomID,omitempty"`
	PDFunds           string `json:"PDFunds,omitempty" bson:"PDFunds,omitempty"`
	PDSent            string `json:"PDSent,omitempty" bson:"PDSent,omitempty"`
	Workshop          string `json:"workshop,omitempty" bson:"workshop,omitempty"`
	//new profile fields
	Responsibility      string `json:"responsibility,omitempty" bson:"responsibility,omitempty"`
	StewardEmail        string `json:"stewardEmail,omitempty" bson:"stewardEmail,omitempty"`
	CBEmail             string `json:"cb_email,omitempty" bson:"cb_email,omitempty"`
	ContractDate        string `json:"contractDate,omitempty" bson:"contractDate,omitempty"`
	ETFONumber          string `json:"etfo_number,omitempty" bson:"etfo_number,omitempty"`
	SCDSBNumber         string `json:"scdsb_number,omitempty" bson:"scdsb_number,omitempty"`
	Seniority           string `json:"seniority,omitempty" bson:"seniority,omitempty"`
	TimeTypeDescription string `json:"timeTypeDescription,omitempty" bson:"timeTypeDescription,omitempty"`
	SeniorityAsOf       string `json:"seniorityAsOf,omitempty" bson:"seniorityAsOf,omitempty"`
	MemberID            string `json:"MemberID,omitempty" bson:"realMemberID,omitempty"`
}

type UserInfo struct {
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	UnionMail        string `json:"unionMail,omitempty" bson:"unionMail,omitempty"`
	ImageURL         string `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	Address          string `json:"address,omitempty" bson:"address,omitempty"`
	City             string `json:"city,omitempty" bson:"city,omitempty"`
	Province         string `json:"province,omitempty" bson:"province,omitempty"`
	PostalCode       string `json:"postalCode,omitempty" bson:"postalCode,omitempty"`
	Phone            string `json:"phone,omitempty" bson:"phone,omitempty"`
	Mobile           string `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Description      string `json:"description,omitempty" bson:"description,omitempty"`
	BannerURL        string `json:"bannerURL,omitempty" bson:"bannerURL,omitempty"`
	Fax              string `json:"fax,omitempty" bson:"fax,omitempty"`
	PresidentMessage string `json:"president_message,omitempty" bson:"president_message,omitempty"`
	// ImportantLinks   []*ImportantLinks `json:"url,omitempty" bson:"url,omitempty"`
	// WebsiteLinks     ImportantLinks    `json:"websiteUrl,omitempty" bson:"websiteUrl,omitempty"`
}

type SingleUserAuth struct {
	User  User   `json:"user,omitempty" bson:"user"`
	Token string `json:"token,omitempty" bson:"token"`
}

type Credential struct {
	UnionID  primitive.ObjectID `json:"unionID,omitempty" bson:"unionID,omitempty"`
	Username string             `json:"username,omitempty" bson:"username"`
	Password string             `json:"password,omitempty" bson:"password"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
}

type UserUploadReport struct {
	NewUsers      []*User  `json:"newUsers,omitempty"`
	UpdatedUsers  []*User  `json:"updatedUsers,omitempty"`
	ErroredUsers  []*User  `json:"erroredUsers,omitempty"`
	ErrorMessages []string `json:"errorMessages,omitempty"`
}

type UserUpdateInput struct {
	FirstName string   `json:"firstName,omitempty" bson:"firstName"`
	LastName  string   `json:"lastName,omitempty" bson:"lastName"`
	Status    string   `json:"profile,omitempty" bson:"profile"`
	Profile   UserInfo `json:"status,omitempty" bson:"status"`
}

type UserFilterInput struct {
	IsAdmin bool               `json:"isAdmin,omitempty"`
	Deleted bool               `json:"deleted,omitempty"`
	Status  string             `json:"status,omitempty"`
	UnionID primitive.ObjectID `json:"unionID,omitempty"`
}
