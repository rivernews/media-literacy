import os
import boto3
import asyncio
from slack_sdk.signature import SignatureVerifier
from media_literacy.http import BadRequestError, handle_exception, APIGatewayRequest
from media_literacy.logging import Logger
from media_literacy.services.slack_service import SlackService


loop = asyncio.get_event_loop()
PIPELINE_QUEUE_NAME = os.environ.get('PIPELINE_QUEUE_NAME', '')


@handle_exception
@APIGatewayRequest.use_request
def lambda_handler(request: APIGatewayRequest, context):
    Logger.debug('Incoming request', request)

    if not SignatureVerifier(os.environ.get('SLACK_SIGNING_SECRET', '')).is_valid_request(
        # Validating AWS Lambda's Event Slack Request
        # https://gist.github.com/nitrocode/288bb104893698011720d108e9841b1f
        request._body,
        request.headers
    ):
        raise BadRequestError

    # Based on
    # https://boto3.amazonaws.com/v1/documentation/api/latest/guide/sqs.html#sending-messages
    sqs = boto3.resource('sqs')
    queue = sqs.get_queue_by_name(QueueName=PIPELINE_QUEUE_NAME)
    # only using a single `MessageGroupId` for this queue - does not intend to use for multiple FIFO orderings in one queue
    response = queue.send_message(MessageBody=str(request.body), MessageGroupId=PIPELINE_QUEUE_NAME)

    loop.run_until_complete(SlackService.send('You sent a slack command!'))

    return {
        'message': 'OK'
    }
