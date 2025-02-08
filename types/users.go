package types

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type DeleteUserRequestBody struct {
	Password string `json:"password" validate:"required,min=8"`
}

type UpdatePasswordRequestBody struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
