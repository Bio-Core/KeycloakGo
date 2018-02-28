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

//Action is a type of possible user actions
type Action string

var (
	//LoginAction is for user logins
	LoginAction Action = "Login"
	//LogoutAction is for user logouts
	LogoutAction Action = "Logout"
	//PageAccessAction is for any user page access
	PageAccessAction Action = "Access"
	//UploadFileAction is for users uploading files via the system
	UploadFileAction Action = "Upload File"
)
