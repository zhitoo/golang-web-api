package otp

import (
	"context"
	"fmt"

	"github.com/zhitoo/golang-web-api/pkg/job"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

type SendSmsJob struct {
	job.Base
	Mobile string `json:"mobile"`
	Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
	log.Info("sending SMS", map[logging.ExtraKey]any{
		"mobile": j.Mobile,
	})

	// TODO: integrate with real SMS provider (Kavenegar, Ghasedak, etc.)
	fmt.Printf("SMS to %s: %s\n", j.Mobile, j.Body)

	return nil
}

func (j *SendSmsJob) Queue() string { return "sms" }

func init() {
	job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}
