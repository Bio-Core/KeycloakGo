package keycloak

//Client Values from JSON file
type Client struct {
	Realm       string `json:"realm"`
	ID          string `json:"resource"`
	Credentials Creds  `json:"credentials"`
}

//Creds is a substruct of Keycloak
type Creds struct {
	Secret string `json:"secret"`
}

type action int

const (
	actionLogin action = iota
	actionLogout
	actionPageAccess
	actionInvalid
)

//Action is a type of possible user actions
type Action string

var (
	//ActionLogin is for user logins
	ActionLogin Action = "Login"
	//ActionLogout is for user logouts
	ActionLogout Action = "Logout"
	//ActionPageAccess is for any user page access
	ActionPageAccess Action = "Access"
)
