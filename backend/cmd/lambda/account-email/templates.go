package main

import (
	"fmt"
	"html"
)

func wrapEmail(heading, intro, tempPassword, loginURL, footnote string) string {
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background-color:#f4f5f7;font-family:Helvetica,Arial,sans-serif;">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f5f7;padding:32px 16px;">
<tr><td align="center">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:480px;background-color:#ffffff;border-radius:12px;overflow:hidden;">
<tr><td style="background-color:#1f2a44;padding:24px 32px;">
<span style="color:#ffffff;font-size:18px;font-weight:bold;">Inspirate Consulting</span>
</td></tr>
<tr><td style="padding:32px;">
<h1 style="margin:0 0 16px;font-size:20px;color:#1f2a44;">%s</h1>
<p style="margin:0 0 20px;font-size:15px;line-height:1.5;color:#3c4257;">%s</p>
<p style="margin:0 0 20px;font-size:14px;line-height:1.6;color:#3c4257;">
<strong>Step 1.</strong> Copy your temporary password below.
</p>
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="margin:0 0 24px;">
<tr><td style="background-color:#f4f5f7;border-radius:8px;padding:16px;text-align:center;">
<span style="font-family:'Courier New',monospace;font-size:18px;font-weight:bold;color:#1f2a44;letter-spacing:1px;">%s</span>
</td></tr>
</table>
<p style="margin:0 0 20px;font-size:14px;line-height:1.6;color:#3c4257;">
<strong>Step 2.</strong> Go to <a href="%s" style="color:#1f2a44;">the sign-in page</a> and paste it to log in.
</p>
<table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 0 24px;">
<tr><td style="background-color:#1f2a44;border-radius:8px;">
<a href="%s" style="display:inline-block;padding:14px 28px;font-size:15px;font-weight:bold;color:#ffffff;text-decoration:none;border-radius:8px;">Go to Sign In</a>
</td></tr>
</table>
<p style="margin:0;font-size:12px;line-height:1.5;color:#8792a2;">%s</p>
</td></tr>
</table>
</td></tr>
</table>
</body>
</html>`, heading, intro, tempPassword, loginURL, loginURL, footnote)
}

func accountCreatedEmail(name, tempPassword, loginURL string) (string, string) {
	subject := "Welcome to Inspirate Consulting — set up your account"
	intro := fmt.Sprintf("Hi %s, your counselor created an Inspirate Consulting account for you.", html.EscapeString(name))
	footnote := "If you didn't expect this email, you can ignore it."
	return subject, wrapEmail("Welcome to Inspirate", intro, html.EscapeString(tempPassword), loginURL, footnote)
}

func passwordResetEmail(name, tempPassword, loginURL string) (string, string) {
	subject := "Your Inspirate Consulting password was reset"
	intro := fmt.Sprintf("Hi %s, your counselor reset your Inspirate Consulting password.", html.EscapeString(name))
	footnote := "If you didn't request this, contact your counselor."
	return subject, wrapEmail("Password reset", intro, html.EscapeString(tempPassword), loginURL, footnote)
}
