ORIGINAL_PATH=$(pwd)
REPO=$(git rev-parse --show-toplevel)

cd ${REPO}/lambda/layer/tests
pytest --log-level=DEBUG -rP
cd ${ORIGINAL_PATH}