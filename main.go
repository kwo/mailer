package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Host     string `required:"true"`
	Port     int    `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

func (z *Config) Address() string {
	return fmt.Sprintf("%s:%d", z.Host, z.Port)
}

type Envelope struct {
	From    *string
	To      *[]string
	Subject *string
}

func (z *Envelope) Validate() error {
	var errs []string
	if z.Subject == nil || *z.Subject == "" {
		errs = append(errs, "missing subject")
	}
	if z.From == nil || *z.From == "" {
		errs = append(errs, "missing from address")
	}
	if z.To == nil || len(*z.To) == 0 {
		errs = append(errs, "missing to address")
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func main() {

	cfg := &Config{}
	if err := envconfig.Process("mailer", cfg); err != nil {
		log.Fatalf("invalid environment configuration: %s\n", err)
	}

	envelope := &Envelope{}
	envelope.Subject = flag.String("subject", "", "subject")
	envelope.From = flag.String("from", "", "from address") // "Jordan Wright <test@gmail.com>"
	envelope.To = flag.StringSlice("to", nil, "to address")
	flag.Parse()
	if err := envelope.Validate(); err != nil {
		log.Fatalf("invalid envelope configuration: %s\n", err)
	}

	msgbody, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("cannot read message from stdin: %s", err)
	}

	e := email.NewEmail()
	e.From = *envelope.From
	e.To = *envelope.To
	e.Subject = *envelope.Subject
	e.Text = msgbody

	if err := e.SendWithStartTLS(
		cfg.Address(),
		smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host),
		&tls.Config{ServerName: cfg.Host},
	); err != nil {
		log.Fatalf("cannot send email: %s\n", err)
	}

}
