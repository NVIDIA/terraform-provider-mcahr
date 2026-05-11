resource "mcahr_smtp_subscription" "smtp_subscription" {
  name             = "smtp_subscription"
  integration_name = "existing_smtp_integration"
  recipients       = ["foo@bar.com", "bar@baz.com"]
  filters = [
    {
      "type" : "TRIGGER",
      "category" : "ACTION",
      "status" : "EXECUTING"
    },
    {
      "type" : "TRIGGER",
      "category" : "ALARM",
      "status" : "EXECUTING"
    },
    {
      "type" : "TRIGGER",
      "category" : "BOT",
      "status" : "EXECUTING"
    },
    {
      "type" : "TRIGGER",
      "category" : "TIME_TRIGGER",
      "status" : "EXECUTING"
    }
  ]
  enabled = true
}
