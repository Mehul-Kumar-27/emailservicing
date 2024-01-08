package handellers

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailServer struct {
	Mail *Mail
}

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          []string
	Subject     string
	Attachments []string
	Data        interface{}
	DataMap     map[string]interface{}
}

func CreateMail() *Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mail := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}
	return &mail
}

func (m *Mail) SendMail(msg *Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]interface{}{
		"message": msg.Data,
	}
	msg.DataMap = data

	formattedMessage, err := m.buildHTMLmessage(msg)
	if err != nil {
		log.Println(err)
		log.Println("1")
		return err
	}

	plainMessage, err := m.buildPlainMessage(msg)
	if err != nil {
		log.Println(err)
		log.Println("2")
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncription(m.Encryption)
	server.SendTimeout = 10 * time.Second
	server.ConnectTimeout = 10 * time.Second
	server.KeepAlive = false

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		log.Println("3")
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To...).
		SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	if err := email.Send(smtpClient); err != nil {
		log.Println(err)
		log.Println("4")
		return err
	}

	return nil
}

func (m *Mail) buildHTMLmessage(msg *Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-template").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formatted, err := m.buildInLineCss(formattedMessage)
	if err != nil {
		return "", err
	}

	return formatted, nil

}

func (m *Mail) buildInLineCss(formattedMessage string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(formattedMessage, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) buildPlainMessage(msg *Message) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)

	}
	log.Println(dir)
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	log.Println(msg.DataMap)

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	return formattedMessage, nil
}

func (m *Mail) getEncription(s string) mail.Encryption {
	switch s {
	case "ssl":
		return mail.EncryptionSSL
	case "tls":
		return mail.EncryptionTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionTLS
	}
}
