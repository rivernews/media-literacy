import os
import boto3
import asyncio
import secrets
import json
from datetime import datetime
from slack_sdk.signature import SignatureVerifier
from media_literacy.http import BadRequestError, handle_exception, APIGatewayRequest
from media_literacy.logging import Logger
from media_literacy.services.slack_service import SlackService, SlackCommandMantra
from media_literacy.archive_bucket import ArchiveBucket


loop = asyncio.get_event_loop()
sqs = boto3.resource('sqs')
sfn = boto3.client('stepfunctions')
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

    command = request.body.get('command', '')
    command = command[0] if command and isinstance(command, list) else command
    res = None
    if command.startswith(SlackCommandMantra.FETCH_LANDING):
        # Based on
        # https://boto3.amazonaws.com/v1/documentation/api/latest/guide/sqs.html#sending-messages
        queue = sqs.get_queue_by_name(QueueName=PIPELINE_QUEUE_NAME)
        res = queue.send_message(MessageBody=str(request.body), MessageGroupId=f'{PIPELINE_QUEUE_NAME}')

    elif command.startswith(SlackCommandMantra.FETCH_LANDING_STORIES):
        text = request.body.get('text', '')
        landing_html_key = text[0] if text and isinstance(text, list) else text
        if not ArchiveBucket.exist(landing_html_key):
            raise BadRequestError(f'Landing html does not exist at s3://{ArchiveBucket.bucket_name}/{landing_html_key}')

        res = sfn.start_execution(
            stateMachineArn=os.environ['BATCH_STORIES_SFN_ARN'],
            name=f'media-literacy-sf-batch-stories-{datetime.now().strftime("%Y-%H-%M")}-{secrets.token_hex(nbytes=5)}',
            input=json.dumps({
                'landingURL': landing_html_key
            })
        )

    loop.run_until_complete(SlackService.send(f'You sent a slack command. Processed response: {res}'))

    return {
        'message': 'OK',
        'res': str(res)
    }
