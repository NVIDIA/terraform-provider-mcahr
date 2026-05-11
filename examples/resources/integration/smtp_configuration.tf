resource "mcahr_integration" "smtp_integration" {
  name                  = "smtp_integration"
  service_name          = "smtp"
  smtp_host             = "smtp.example.com"
  smtp_port             = 587
  username              = "foo@bar.com"
  password              = "password"
  sender                = "foo@bar.com"
  max_emails_per_day    = 1000
  max_emails_per_second = 10
  enabled               = true
}
