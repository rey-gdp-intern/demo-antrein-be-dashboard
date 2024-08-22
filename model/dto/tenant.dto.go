package dto

type Tenant struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTenantRequest struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

type CreateTenantResponse struct {
	Tenant Tenant `json:"tenant"`
	Token  string `json:"token"`
}
