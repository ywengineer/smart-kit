package mailx

import (
	"errors"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type gomailClient struct {
	client gomail.SendCloser
}

func (mc *gomailClient) Send(from, to, cc, bcc, subject, bodyType, bodyString string) error {
	var rec []string
	if utilk.ValidMail(from) == false {
		return errors.New("missing mail's sender")
	}
	if utilk.ValidMail(to) == false {
		return errors.New("missing mail's to")
	}
	rec = append(rec, to)
	//
	if len(cc) > 0 {
		if utilk.ValidMail(cc) == false {
			return errors.New("unknown mail's cc. %s" + cc)
		}
		rec = append(rec, cc)
	}

	//
	if len(bcc) > 0 {
		if utilk.ValidMail(bcc) == false {
			return errors.New("unknown mail's bcc. %s" + bcc)
		}
		rec = append(rec, bcc)
	}
	//
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetAddressHeader("Cc", cc, cc)
	}
	if len(bcc) > 0 {
		m.SetAddressHeader("Bcc", bcc, bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, bodyString)
	// Send the email to Bob, Cora and Dan.
	if err := mc.client.Send(from, rec, m); err != nil {
		return err
	}
	return nil
}

// NewGoMailSender creates a new mail sender.
func NewGoMailSender(host string, port int, username, password string) (Mailer, error) {
	d := gomail.NewDialer(host, port, username, password)
	if c, err := d.Dial(); err != nil {
		return nil, err
	} else {
		return &gomailClient{client: c}, nil
	}
}

func DirectSendMail(host string, port int, username, password string,
	from, to, cc, bcc string, subject, bodyType, bodyString string) {
	if utilk.ValidMail(from) == false {
		logk.Error("missing mail's sender")
		return
	}
	if utilk.ValidMail(to) == false {
		logk.Error("missing mail's to")
		return
	}
	if len(cc) > 0 && utilk.ValidMail(cc) == false {
		logk.Error("unknown mail's cc. %s", zap.String("cc", cc))
		return
	}
	if len(bcc) > 0 && utilk.ValidMail(bcc) == false {
		logk.Error("unknown mail's bcc. %s", zap.String("bcc", bcc))
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetAddressHeader("Cc", cc, cc)
	}
	if len(bcc) > 0 {
		m.SetAddressHeader("Bcc", bcc, bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, bodyString)
	//m.Attach("/home/Alex/lolcat.jpg")
	d := gomail.NewDialer(host, port, username, password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logk.Error("send mail failed, %v", zap.Error(err))
	}
}
