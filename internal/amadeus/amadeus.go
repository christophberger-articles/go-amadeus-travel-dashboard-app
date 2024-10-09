package amadeus

// tokenResponse contains either a valid access token
// or an error that occurred while fetching the token
type tokenResponse struct {
	Token string
	Err   error
}

// Client is a client for the Amadeus API.
// It takes care of refreshing the access token regularly
// in the background while serving the currently valid
// token to clients
type Client struct {
	baseURL     string
	accessToken chan tokenResponse
}

// Create a new client and start the token refreshing goroutine.
func New() *Client {
	c := &Client{
		baseURL:     "https://test.api.amadeus.com/v1",
		accessToken: make(chan tokenResponse),
	}
	go c.refreshToken()
	return c
}

// AuthResponse contains the unmarshaled response from the Amadeus
// authorization API.
// It is a blend of the success and error response, so that we can
// unmarshal the response first and check for success or error later.
type AuthResponse struct {
	AuthSuccessResponse
	AuthErrorResponse
}

type AuthSuccessResponse struct {
	Type            string `json:"type"`
	Username        string `json:"username"`
	ApplicationName string `json:"application_name"`
	ClientID        string `json:"client_id"`
	TokenType       string `json:"token_type"`
	AccessToken     string `json:"access_token"`
	ExpiresIn       int    `json:"expires_in"`
	State           string `json:"state"`
	Scope           string `json:"scope"`
}

type AuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Code             int    `json:"code"`
	Title            string `json:"title"`
}

