terraform {
  required_providers {
    todoist = {
      # This provider is not published on the public registry; use TF_CLI_CONFIG_FILE + dev_overrides.
      source = "andreaswwilson/todoist"
    }
  }
}

provider "todoist" {
  # Set via env var: export TODOIST_TOKEN="..."
}

resource "todoist_projects" "example" {
  name  = "tf-managed-example"
  color = "charcoal"
}

output "project_id" {
  value = todoist_projects.example.id
}
