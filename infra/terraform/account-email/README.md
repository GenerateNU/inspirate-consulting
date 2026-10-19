# Account email infra

Provisions the SQS queue + Lambda that deliver account-creation and
password-reset emails (see `backend/cmd/lambda/account-email`). This is
separate from the existing `inspirate-email-notification-queue` /
`email-notification` lambda (task due-date reminders) — different message
shape, different templates, left untouched.

Unlike that existing lambda (manual console upload, see its README), this one
is managed by Terraform end to end.

## Build the Lambda package

Lambda runs on Amazon Linux; the binary must be cross-compiled and named
`bootstrap`, matching the `arm64` architecture set in `main.tf`:

```sh
cd backend
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o ../infra/terraform/account-email/build/bootstrap ./cmd/lambda/account-email
```

Terraform zips that binary itself (`data.archive_file`), so there's no manual
`zip` step — just make sure the binary exists at `lambda_zip_path`
(`build/bootstrap` by default) before `terraform plan`/`apply`.

## Apply

```sh
cd infra/terraform/account-email
terraform init
terraform plan \
  -var="sender_email=no-reply@inspirate.org" \
  -var="app_url=https://inspirate.org" \
  -var="claim_token_secret=$CLAIM_TOKEN_SECRET"
terraform apply # same -var flags
```

`sender_email` must already be a verified SES identity. SES starts in sandbox
mode, so recipient addresses need to be verified too until production access
is granted — same constraint the existing lambda's README calls out.

`claim_token_secret` must match the `CLAIM_TOKEN_SECRET` the backend API uses
to verify claim-password tokens (`backend/internal/claimtoken`) — the Lambda
signs tokens with this secret, the backend verifies them with the same one.
Keep it in a `terraform.tfvars` file (gitignored) or pass it from your
secrets manager of choice; don't commit it.

No remote state backend is configured — this applies with local state until
one is wired up. Treat `terraform.tfstate` as sensitive (it contains the
Lambda's environment variables, including the secret above) and never commit
it; the root `.gitignore` already excludes `*.tfstate*`.

## After applying

Set the backend's `SQS_ACCOUNT_QUEUE_URL` env var (see
`backend/internal/config/sqs.go`) to this module's `queue_url` output, and
redeploy the backend so `CreateUser` and the password-reset endpoint start
enqueueing to it (see PR 3 / PR 4).

## Redeploying after a code change

Rebuild the binary (above), then:

```sh
terraform apply -var="sender_email=..." -var="app_url=..." -var="claim_token_secret=..."
```

`source_code_hash` picks up the change and updates the function in place.
