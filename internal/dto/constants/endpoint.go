package constants

const (
	InvalidRequestArgs = "invalid HTTP-Request parameters"

	ApiGroupV1          = "/api/v1"
	ApiHomeUrl          = "/"
	ApiPingUrl          = "/ping"
	ApiAuthLoginUrl     = "/auth/login"
	ApiAuthRegisternUrl = "/auth/register"
	ApiRefreshTokenUrl  = "/auth/token-refresh"
	ApiProfileCreateUrl = "/profile"
	ApiProfileGetUrl    = "/profile/:id"
	ApiReportCreateUrl  = "/report"
	ApiReportGetUrl     = "/report/:id"
)

const (
	UserIdCookieKey       = "user_id"
	AccessTokenCookieKey  = "access-token"
	RefreshTokenCookieKey = "refresh-token"
)
