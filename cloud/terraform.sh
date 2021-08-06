set -o allexport

. ./local.backend.credentials.tfvars
. ./local.credentials
AWS_ACCESS_KEY_ID=${access_key}
AWS_SECRET_ACCESS_KEY=${secret_key}
AWS_DEFAULT_REGION=${region}
TF_VAR_project_name=media-literacy
TF_VAR_slack_signing_secret=${slack_signing_secret}
TF_VAR_slack_post_webhook_url=${slack_post_webhook_url}
TF_VAR_s3_archive_bucket=${s3_archive_bucket}
set +o allexport

terraform "$@"
