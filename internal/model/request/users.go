package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type VerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserProfile struct {
	Name      string   `json:"name"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	Verify    bool     `json:"verify"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type UpdateUserProfileRequest struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
  FirstName string   `json:"first_name"`
  LastName  string   `json:"last_name"`
}

type UserResponse struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}
