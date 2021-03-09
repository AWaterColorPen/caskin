package example

import "github.com/awatercolorpen/caskin"

// Stage example Stage for easy testing
type Stage struct {
	Caskin         *caskin.Caskin  // caskin instance on stage
	Options        *caskin.Options // caskin options on stage
	Domain         *Domain         // a domain on stage
	SuperadminUser *User           // superadmin user on stage
	AdminUser      *User           // a domain admin user on stage
	MemberUser     *User           // a domain member user on stage
}
