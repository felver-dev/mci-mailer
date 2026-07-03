package dto

type EmailAddress struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"`
}

type Attachment struct {
	Filename string `json:"filename" binding:"required"`
	Content  string `json:"content" binding:"required"` // base64
	MimeType string `json:"mime_type"`
}

type SendRequest struct {
	From         *EmailAddress     `json:"from"`
	To           []EmailAddress    `json:"to" binding:"required,min=1,dive"`
	CC           []EmailAddress    `json:"cc" binding:"omitempty,dive"`
	BCC          []EmailAddress    `json:"bcc" binding:"omitempty,dive"`
	ReplyTo      *EmailAddress     `json:"reply_to"`
	Subject      string            `json:"subject"`
	Html         string            `json:"html"`
	Text         string            `json:"text"`
	TemplateName string            `json:"template"`
	Variables    map[string]any    `json:"variables"`
	Attachments  []Attachment      `json:"attachments" binding:"omitempty,dive"`
}

type SendResponse struct {
	LogID   string `json:"log_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LogFilter struct {
	Status   string `form:"status"`
	ApiKeyID string `form:"api_key_id"`
	From     string `form:"from"`
	To       string `form:"to"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=50"`
}
