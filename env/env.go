package env

import (
	"os"

	"github.com/sirupsen/logrus"
)

const Project = "democracy-tools"

func GetSmtp() string {

	const key = "SMTP"
	res := os.Getenv(key)
	if res == "" {
		res = "smtp.gmail.com"
	}
	logrus.Debugf("%s: %s", key, res)

	return res
}

func GetEmailSupport() string {

	return failIfEmpty("EMAIL_SUPPORT")
}

func GetEmailFrom() string {

	return failIfEmpty("EMAIL_FROM")
}

func GetEmailPassword() string {

	return failIfEmpty("EMAIL_PASSWORD")
}

func GetWhatsAppVerificationToken() string {

	return failIfEmpty("WHATSAPP_VERIFICATION_TOKEN")
}

func GetWhatAppToken() string {

	return failIfEmpty("WHATSAPP_TOKEN")
}

func GetWhatsAppFromPhone() string {

	return failIfEmpty("WHATSAPP_FROM_PHONE")
}

func GetWhatsAppTemplate() string {

	return failIfEmpty("WHATSAPP_TEMPLATE")
}

func GetSlackInfoUrl() string {

	return failIfEmpty("SLACK_INFO_URL")
}

func GetSlackDebugUrl() string {

	return failIfEmpty("SLACK_DEBUG_URL")
}

func failIfEmpty(key string) string {

	res := os.Getenv(key)
	if res == "" {
		logrus.Fatalf("Please, add environment variable '%s'", key)
	}
	logrus.Debugf("%s: %s", key, res)

	return res
}
