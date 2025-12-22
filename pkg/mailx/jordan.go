package mailx

import (
	"fmt"
	"net/smtp"
	"net/textproto"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/jordan-wright/email"
	"gopkg.in/errgo.v2/errors"
)

type jordanClient struct {
	p *email.Pool
}

func (mc *jordanClient) Send(from, to, cc, bcc, subject, bodyType, bodyString string) error {
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
	e := &email.Email{
		To:      []string{to},
		From:    from,
		Cc:      []string{cc},
		Bcc:     []string{bcc},
		Subject: subject,
		HTML:    []byte(bodyString),
		Headers: textproto.MIMEHeader{
			consts.HeaderContentType: []string{bodyType},
		},
	}
	return mc.p.Send(e, 30*time.Second)
}

// NewJordanMailSender creates a new mail sender.
func NewJordanMailSender(host string, port int, username, password string) (Mailer, error) {
	p, err := email.NewPool(
		fmt.Sprintf("%s:%d", host, port),
		4,
		smtp.PlainAuth("", username, password, host),
	)
	if err != nil {
		return nil, err
	}
	return &jordanClient{p}, nil
}
