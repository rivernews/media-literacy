set -e

ENV=dev

set -o allexport

. ./local.backend.credentials.tfvars
. ./local.credentials
AWS_ACCESS_KEY_ID=${access_key}
AWS_SECRET_ACCESS_KEY=${secret_key}
AWS_DEFAULT_REGION=${region}
TF_VAR_project_alias=media-literacy
TF_VAR_environment_name=${ENV}
TF_VAR_slack_signing_secret=${slack_signing_secret}
TF_VAR_slack_post_webhook_url=${slack_post_webhook_url}
set +o allexport


DEPLOY_DIR=$(git rev-parse --show-toplevel)/cloud_environments/${ENV}
SCRAPER_SRC_DIR=$(git rev-parse --show-toplevel)/scraper_lambda

cd $SCRAPER_SRC_DIR

if go build -o main; then
    cd $DEPLOY_DIR

    # if deploy the first time, uncomment below
    # to avoid "Invalid for_each argument" error
    # https://github.com/terraform-aws-modules/terraform-aws-step-functions/issues/20
    # terraform "$@" -target=module.main.module.scraper_lambda

    terraform "$@"
else
    echo "Go build failed"
    cd $DEPLOY_DIR
    exit 1
fi
