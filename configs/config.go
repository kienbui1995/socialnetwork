package configs

//Const config system
const (
	neo4jURL = "Bolt://neo4j:tlis2016@localhost:7687"
	URLDB    = "http://neo4j:madawg00@localhost:7474/db/data/"
	APIPort  = 8080
)

//PrivacyType uint
type PrivacyType uint

// Const privacy
const (
	Public           = 1
	ShareToFollowers = 2
	Private          = 3
)

// ErrorCode Table
const (
	APIEcSuccess                           = 1   //	Success
	APIEcNoExistObject                     = 2   //	No exist this object.
	APIEcParam                             = 100 //	Invalid parameter
	APIEcParamMissingField                 = 101 //	Missing a few fields.
	APIEcParamUserID                       = 110 //	Invalid user id
	APIEcParamUserField                    = 111 //	Invalid user info field
	APIEcParamEmail                        = 113 //	Invalid email
	APIEcParamFieldList                    = 115 //	Invalid field list
	APIEcParamPhotoID                      = 121 //	Invalid photo id
	APIEcParamTitle                        = 142 //	Invalid title
	APIEcParamAccessToken                  = 190 //	Invalid OAuth 2.0 Access Token
	APIEcPermission                        = 200 //	Permissions error
	APIEcPermissionUser                    = 210 //	User not visible
	APIEcPermissionPhoto                   = 221 //	Photo not visible
	APIEcPermissionMessage                 = 230 //	Permissions disallow message to user
	APIEcEdit                              = 300 //	Edit failure
	APIEcEditUserData                      = 310 //	User data edit failure
	APIEcUsersCreateInvalidEmail           = 370 //	The email address you provided is not a valid email address
	APIEcUsersCreateExistingEmail          = 371 //	The email address you provided belongs to an existing account
	APIEcUsersCreateBirthday               = 372 //	The birthday provided is not valid
	APIEcUsersCreatePassword               = 373 //	The password provided is too short or weak
	APIEcUsersRegisterInvalidCredential    = 374 //	The login credential you provided is invalid.
	APIEcUsersRegisterConfFailure          = 375 //	Failed to send confirmation message to the specified login credential.
	APIEcUsersRegisterExisting             = 376 //	The login credential you provided belongs to an existing account
	APIEcUsersRegisterDefaultError         = 377 //	Sorry, we were unable to process your registration.
	APIEcUsersRegisterPasswordBlank        = 378 //	Your password cannot be blank. Please try another.
	APIEcUsersRegisterPasswordInvalidChars = 379 //	Your password contains invalid characters. Please try another.
	APIEcUsersRegisterPasswordShort        = 380 //	Your password must be at least 6 characters long. Please try another.
	APIEcUsersRegisterPasswordWeak         = 381 //	Your password should be more secure. Please try another.
	APIEcUsersRegisterUsernameError        = 382 //	Please enter a valid username.
	APIEcUsersRegisterMissingInput         = 383 //	You must fill in all of the fields.
	APIEcUsersRegisterIncompleteBday       = 384 //	You must indicate your full birthday to register.
	APIEcUsersRegisterInvalidEmail         = 385 //	Please enter a valid email address.
	APIEcUsersRegisterEmailDisabled        = 386 //	The email address you entered has been disabled. Please contact disabled@facebook.com with any questions.
	APIEcUsersRegisterAddUserFailed        = 387 //	There was an error with your registration. Please try registering again.
	APIEcUsersRegisterNoGender             = 388 //	Please select either Male or Female.
	APIEcAuthEmail                         = 400 //	Invalid email address
	APIEcAuthLogin                         = 401 //	Invalid username or password
	APIEcAuthMissingToken                  = 404 //	Missing token.
	APIEcAuthInvalidToken                  = 405 //	Invalid token.
	APIEcAuthNoExistToken                  = 406 //	No exist token.
	APIEcAuthCheckToken                    = 407 //	Error in checking token.
	APIEcAuthGenerateToken                 = 408 //	Error in generate token.
	APIEcAuthNoExistUser                   = 409 //	No exist user.
	APIEcAuthNoExistFacebook               = 410 //	No exist account with this facebook.
	APIEcAuthInvalidFacebookToken          = 411 //	Error in checking token.
	APIEcAuthWrongPassword                 = 412 //	Error in login: Wrong password.
	APIEcAuthNoExistEmail                  = 413 //
	APIEcAuthWrongRecoveryCode             = 414 //	Error in recover password: Wrong recovery code.
	APIEcMesgNoBody                        = 501 //	Missing message body
)

// TypePost const
const (
	Post           = 0
	PostStatus     = 1
	PostPhoto      = 2
	PostLink       = 3
	PostGroup      = 4
	PostSharePost  = 5
	PostSharePage  = 6
	PostShareGroup = 7
)

//FCMToken struct
const (
	FCMToken = "AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V"
)
