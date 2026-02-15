# Tryout

Local-only test harness for the provider.

## Configure Terraform to Use the Local Provider Binary

1. Build/install the provider binary into GOPATH/bin:

```bash
cd ..
go install
```

2. Run Terraform using the repo-local CLI config so Terraform finds the binary:

```bash
cp ../dev.tfrc ../dev.tfrc.local
$EDITOR ../dev.tfrc.local
export TF_CLI_CONFIG_FILE="$PWD/../dev.tfrc.local"
export TODOIST_TOKEN="..."
terraform init
terraform plan
```

To actually create the project:

```bash
terraform apply
```

## Notes

- Auth env var is `TODOIST_TOKEN` in this provider.
- Data source `todoist_projects` requires a project ID and reads a single project.
