set -e

REPO_DIR=$(git rev-parse --show-toplevel)
ENV=dev sh "${REPO_DIR}/cloud_environments/terraform.sh" "$@"
