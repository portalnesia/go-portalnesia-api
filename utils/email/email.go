package email

import util "portalnesia.com/api/utils"

const (
	EmailFromSupport     = "support"
	EmailFromNoReply     = "noreply"
	EmailFromInfo        = "info"
	TemplateTypeBasic    = "basicv2"
	TemplateTypeBirthdat = "birthday"
	TemplateTypeSecurity = "security"
	TemplateTypeRegister = "register"
	TemplateTypeForgot   = "forget"
)

type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	Content     []byte `json:"content"`
	ContentType string `json:"contentType"`
}

type IMailTemplateOptions struct {
	// With Constanta Helper
	Type string `json:"type"`
}

type ITemplateButton struct {
	// Text string with \n
	Text *string `json:"text"`
	Url  *string `json:"url"`
	// Text Only
	Label *string `json:"label"`
}
type ITemplateFooter struct {
	Email *string `json:"email"`
	Url   *string `json:"url"`
	// Text string with \n
	Extra *string `json:"extra"`
}

type ITemplateOptions struct {
	Username string `json:"username"`
	// Text string with \n
	Messages *string          `json:"messages"`
	Header   *string          `json:"header"`
	Button   *ITemplateButton `json:"button"`
	Footer   *ITemplateFooter `json:"footer"`
	ReplyTo  *bool            `json:"replyTo"`
}

type IMailOptions struct {
	// With Constanta Helper
	Email string `json:"email"`
	// With Constanta Helper
	ReplyTo *string `json:"replyTo"`
	Subject string  `json:"subject"`
	// MessageID this message is replyting
	InReplyTo   *string              `json:"inReplyTo"`
	Attachments []Attachment         `json:"attachments"`
	To          []EmailAddress       `json:"to"`
	Template    IMailTemplateOptions `json:"template"`
}

type MailingListOptions struct {
	// With Constanta Helper
	Email string `json:"email"`
	// With Constanta Helper
	ReplyTo *string `json:"replyTo"`
	Subject string  `json:"subject"`
	// MessageID this message is replyting
	InReplyTo   *string      `json:"inReplyTo"`
	Attachments []Attachment `json:"attachments"`
	To          EmailAddress `json:"to"`
	Html        string       `json:"html"`
	Id          string       `json:"id"`
	LogSlug     string       `json:"log_slug"`
}

func HelperTo(to string) []EmailAddress {
	return []EmailAddress{
		{
			Name:    to,
			Address: to,
		},
	}
}

func HelperToMany(to []string) []EmailAddress {
	var a []EmailAddress

	for _, t := range to {
		a = append(a, EmailAddress{Name: t, Address: t})
	}
	return a
}

func SendEmail(options IMailOptions) {
	go util.RabbitmqQueue(util.CHANNEL_SEND_EMAIL, options)
}
