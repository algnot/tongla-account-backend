package util

func GetEmailContent(mapper string) string {
	if mapper == "verifyEmail" {
		return `Hello, %s

To verify your Tongla account, we need to confirm your email. Please click the following link or copy & paste into the browser:

%s/auth/verify-email?token=%s

The link is expird in 30 minutes.

	Best regards,
	Tongla
www.tongla.dev`
	}

	if mapper == "login" {
		return `Hello, %s

To login your Tongla account. Please click the following link or copy & paste into the browser:

%s/auth/login-with-token?token=%s

The link is expird in 30 minutes.

Best regards,
Tongla
www.tongla.dev`
	}

	return ""
}

func GetWebNotificationContent(mapper string) string {
	if mapper == "emailVerified" {
		return `Your email %s is verified`
	}

	if mapper == "login" {
		return `You have logged in with %s via %s device.`
	}

	if mapper == "deviceDelete" {
		return `You have deleted device id %s via %s`
	}

	return ""
}
