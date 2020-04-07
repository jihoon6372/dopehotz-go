package models

// User Model
type User struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	IsStaff     bool    `json:"is_staff"`
	IsActive    bool    `json:"is_active"`
	IsSuperuser bool    `json:"is_superuser"`
	Profile     Profile `json:"profile"`
}

// TableName 유저 테이블명 반환
func (User) TableName() string {
	return "auth_user"
}

// Profile Model
type Profile struct {
	SoundcloudURL     string `json:"soundcloud_url"`
	ProfilePicture    string `json:"profile_picture"`
	Greeting          string `json:"greeting"`
	ClipsGreeting     string `json:"clips_greeting"`
	LikesGreeting     string `json:"likes_greeting"`
	Nickname          string `json:"nickname"`
	SoundcloudID      int    `json:"soundcloud_id"`
	Crew              string `json:"crew"`
	Location          string `json:"location"`
	MailingAgree      bool   `json:"mailing_agree"`
	AccessTerms       string `json:"access_terms"`
	PrivacyStatements string `json:"privacy_statements"`
	YoutubeURL        string `json:"youtube_url"`
	InstagramURL      string `json:"instagram_url"`
	FacebookURL       string `json:"facebook_url"`
	TwitterURL        string `json:"twitter_url"`
	AudiomackSlug     string `json:"audiomack_slug"`
	Soundtrack        string `json:"soundtrack"`
}

// TableName 프로필 테이블명 반환
func (Profile) TableName() string {
	return "accounts_profile"
}
