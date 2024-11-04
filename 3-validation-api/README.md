## Password reset service

#### To set your configuration

- Create .env file
- Set up the following fields

```
  "SMTP_HOST"
	"SMTP_PORT", 587
  "SMTP_USERNAME", "api"
  "SMTP_PASSWORD", "dd148f4f71dcc620d648a25ec81e241b"
  "SMTP_FROM", "noreply@demomailtrap.com"
```

#### How to use

Send a POST request with your email to receive reset link
