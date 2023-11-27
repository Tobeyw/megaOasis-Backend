package consts

const (
	Access_Token                = "AAAAAAAAAAAAAAAAAAAAAGrMfgEAAAAAaeG%2BSKENcr%2FuWvRWSCSnKXEoChE%3DQP2RKmaYDbMDCte2Oxe6XbcAToOMtfc8UQXQ45wx9Gu9sX3ofy"
	TwitterAuthorizeEndpoint    = "https://twitter.com/i/oauth2/authorize"
	TwitterResponseType         = "code"
	TwitterClientID             = "RnhqN29salVXX0hLdUxXanZQelM6MTpjaQ"
	TwitterClientScret          = "2OSLp7aEk4SO3VkjvO2FnBymfagTwmJhXIXQ0QFHhQPjxoHTXk"
	TwitterScope                = "tweet.read%20users.read" //"tweet.read%20users.read"
	TwitterRedirectURITest      = "https://megaoasis.ngd.network:8889/twitter/callback"
	TwitterRedirectURIMain      = "https://megaoasis.ngd.network:8893/twitter/callback"
	TwitterCodeChallenge        = "megaoasis"
	TwitterCodeChallengeMethod  = "plain"
	TwitterAccessTokenGrantType = "authorization_code"
	TwitterAccessTokenEndpoint  = "https://api.twitter.com/2/oauth2/token"
	TwitterGetUserInfoEndpoint  = "https://api.twitter.com/2/users/me"
	FrontEndRedirectUrlMain     = "https://megaoasis.io/account/profile"
	FrontEndRedirectUrlTest     = "http://20.24.36.189:3005/zh/account/profile"
	TwitterErrorPage            = "https://megaoasis.ngd.network:8889/twitter/error"
	NNSContractHash             = "0x50ac1c37690cc2cfc594472833cf57505d5f46de"
	DiscordAuthorizeEndpoint    = "https://discord.com/api/oauth2/authorize"
	DiscordAccessTokenEndpoint  = "https://discord.com/api/oauth2/token"
	DiscordGetUserInfoEndpoint  = "https://discord.com/api/v8/users/@me"
	DiscordRedirectURI          = "http://localhost:8888/discord/callback"
	DiscordResponseType         = "code"
	DiscordClientID             = "1176715427863343176"
	DiscordClientSecret         = "Hla1PRaBf8PCMX1nKABOfDlSkFoTL8e1"
	DiscordScope                = "identify"
	DiscordPrompt               = "consent"
	DiscordAccessTokenGrantType = "authorization_code"
	DiscordRedirectURITest      = "https://megaoasis.ngd.network:8889/discord/callback"
	DiscordRedirectURIMain      = "https://megaoasis.ngd.network:8893/discord/callback"
	DiscordErrorPage            = "https://megaoasis.ngd.network:8889/discord/error"
)
