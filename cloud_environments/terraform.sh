set -e

DEPLOY_DIR=$(git rev-parse --show-toplevel)/cloud_environments/${ENV:-production}
GOLANG_SRC_DIR=$(git rev-parse --show-toplevel)/lambda_golang
PYTHON_SRC_DIR=$(git rev-parse --show-toplevel)/lambda

set -o allexport
. ${DEPLOY_DIR}/local.backend.credentials.tfvars
. ${DEPLOY_DIR}/local.credentials
AWS_ACCESS_KEY_ID=${access_key}
AWS_SECRET_ACCESS_KEY=${secret_key}
AWS_DEFAULT_REGION=${region}
TF_VAR_project_alias=media-literacy
TF_VAR_environment_name=${ENV:-}
TF_VAR_slack_signing_secret=${slack_signing_secret}
TF_VAR_slack_post_webhook_url=${slack_post_webhook_url}
set +o allexport


if (
    cd $GOLANG_SRC_DIR && \
    go build ./cmd/landing && \
    go build ./cmd/stories && \
    go build ./cmd/story && \
    cd $PYTHON_SRC_DIR && python -m compileall layer src
); then
    cd $DEPLOY_DIR

    echo "Go built success"
    echo "Launching terraform..."

    # if deploy the first time, uncomment below
    # to avoid "Invalid for_each argument" error
    # https://github.com/terraform-aws-modules/terraform-aws-step-functions/issues/20
    # terraform "$@" \
    #     -target=module.main.module.scraper_lambda \
    #     -target=module.main.module.batch_stories_fetch_parse_lambda

    terraform "$@"
else
    cd $DEPLOY_DIR

    echo "Go build failed"
    exit 1
fi
