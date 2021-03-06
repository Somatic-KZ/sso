package api

type SignUPRequest struct {
	Email      string  `json:"email" validate:"required,email,unique_email"`
	Phone      string  `json:"phone" validate:"required,is_phone,unique_phones,unique_phone"`
	Password   string  `json:"password" validate:"required,gte=6"`
	FirstName  string  `json:"firstname,omitempty"`
	LastName   string  `json:"lastname,omitempty"`
	Patronymic string  `json:"patronymic,omitempty"`
	Lang       string  `json:"lang,omitempty" validation:"oneof='ru en kk'"`
	Sex        *string `json:"sex,omitempty" validation:"oneof='male female unknown'"`
	BirthDate  *string `json:"birth_date,omitempty" validation:"datetime=2006-01-02"`
	IIN        *int    `json:"iin,omitempty" validate:"omitempty,iin"`
	Region     *string `json:"region,omitempty"`
	City       *string `json:"city,omitempty"`
	Street     *string `json:"street,omitempty"`
	Corpus     *string `json:"corpus,omitempty"`
	House      *string `json:"house,omitempty"`
	Apartment  *string `json:"apartment,omitempty"`
	Zipcode    *int    `json:"zipcode,omitempty"`
}

type SignUPFastRequest struct {
	Email      string  `json:"email" validate:"required,email,unique_email"`
	Phone      string  `json:"phone" validate:"required,is_phone,unique_phone,unique_phones"`
	FirstName  string  `json:"firstname,omitempty"`
	LastName   string  `json:"lastname,omitempty"`
	Patronymic string  `json:"patronymic,omitempty"`
	Lang       string  `json:"lang,omitempty" validation:"oneof='ru en kk'"`
	Sex        *string `json:"sex,omitempty" validation:"oneof='male female unknown'"`
	BirthDate  *string `json:"birth_date,omitempty" validation:"datetime=2006-01-02"`
	IIN        *int    `json:"iin,omitempty" validate:"omitempty,iin"`
	Region     *string `json:"region,omitempty"`
	City       *string `json:"city,omitempty"`
	Street     *string `json:"street,omitempty"`
	Corpus     *string `json:"corpus,omitempty"`
	House      *string `json:"house,omitempty"`
	Apartment  *string `json:"apartment,omitempty"`
	Zipcode    *int    `json:"zipcode,omitempty"`
}

type SignUPResponse struct {
	TDID   string `json:"tdid"`
	Status string `json:"status"`
}

// SignUPOrganizationResponse ???????????? ?????????? ???? ???????????? ??????????????????????
// ???????????????????????? ????????.
type SignUPOrganizationResponse struct {
	TDID   string `json:"tdid"`
	Status string `json:"status"`
}

type NewJWTTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Status       string `json:"status"`
}

// SignInByEmailRequest ?????????????? ?????????????????? ?????? ???????????? ?????????? ?? ????????????
// ????????????????????????.
type SignInByEmailRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

// SignInByPhoneRequest ?????????????? ?????????????????? ?????? ???????????? ???????????????? ?? ????????????
// ????????????????????????.
type SignInByPhoneRequest struct {
	Password string `json:"password" validate:"required,gte=6"`
	Phone    string `json:"phone" validate:"required,is_phone"`
}

// SignInByLoginRequest ?????????????? ?????????????????? ?????? ???????????? ?????????? ?? ????????????
// ????????????????????????.
type SignInByLoginRequest struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type SignInFastRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
	Email string `json:"email" validate:"required,email"`
	TDID  string `json:"tdid"`
}

type SignInFastResponse struct {
	MaskedPhone string `json:"maskedPhone"`
	Status      string `json:"status"`
}

type SignInFastTokenRequest struct {
	TDID  string `json:"tdid,omitempty" validate:"required_without=Phone"`
	Phone string `json:"phone,omitempty" validate:"required_without=TDID,omitempty,is_phone"`
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
	TDID  string `json:"tdid"`
}

type VerifyPhoneNumRequest struct {
	TDID  string `json:"tdid,omitempty" validate:"required_without=Phone"`
	Phone string `json:"phone,omitempty" validate:"required_without=TDID,omitempty,is_phone"`
	Token string `json:"token"`
}

// RecoveryByPhoneRequest ???????????????? ???????????? ???? ?????????????? ???????????????????? ??????????????????
// ???????????????????????????? ?????????????? ?? ???????????????? ???? ???????????? ????????????????.
type RecoveryByPhoneRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
}

type RecoveryByPhoneValidateRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
	Token string `json:"token"`
}

type RecoveryByPhoneNewPasswordRequest struct {
	Phone    string `json:"phone" validate:"required,is_phone"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// DelegateTokenRequest
// ???????????? ???? ?????????????????? ?????????????????????????? ???????????? ?? ????????????????????????
// ???????????????? ???????????????????????? (????????????????????) ???? ???????????? ????????????????.
type DelegateTokenRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"` // ?????????? ???????????????? ????????????????????
}

// SignInDelegateTokenRequest
// ?????????????????? ?????????????????????????? ???????????? ?????? ????????????????
// ???? ???????????? ???????????????? ?????????????? ?? OTP ???????? ???? CMC
type SignInDelegateTokenRequest struct {
	Phone      string `json:"phone,omitempty" validate:"required,omitempty,is_phone"` // ?????????? ???????????????? ????????????????????
	OTP        string `json:"otp"`                                                    // ?????????????????????? ?????? ???? SMS
	MerchantID string `json:"merchant_id"`                                            // ID ????????????????
}

// NewDelegateJWTTokenResponse
// ???????????????????????? ?????????? ???? ???????????????? ???????????????? ??????????????????
type NewDelegateJWTTokenResponse struct {
	DelegateToken string `json:"delegate_token"` // ???????????????????????? ??????????
	Status        string `json:"status"`         // ???????????? ???????????????? ?????????????????? ????????????
}
