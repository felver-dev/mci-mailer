package dto

type CreateTemplateRequest struct {
	Name      string   `json:"name" binding:"required,alphanum"`
	Subject   string   `json:"subject" binding:"required"`
	HtmlBody  string   `json:"html_body" binding:"required"`
	TextBody  string   `json:"text_body"`
	Variables []string `json:"variables"`
}

type UpdateTemplateRequest struct {
	Subject  string `json:"subject"`
	HtmlBody string `json:"html_body"`
	TextBody string `json:"text_body"`
}

type TemplateResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Subject   string   `json:"subject"`
	HtmlBody  string   `json:"html_body"`
	TextBody  string   `json:"text_body"`
	Variables []string `json:"variables"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
