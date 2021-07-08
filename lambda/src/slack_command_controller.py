import os
import base64
import boto3
import json
import asyncio
from datetime import datetime
from slack_sdk.signature import SignatureVerifier
from media_literacy.http import HttpResponse, HttpError, BadRequestError, handle_exception
from media_literacy.logging import Logger
from media_literacy.services.slack_service import SlackService


loop = asyncio.get_event_loop()


@handle_exception
def lambda_handler(event, context):
    Logger.debug('Incoming request', event)

    if not SignatureVerifier(os.environ.get('SLACK_SIGNING_SECRET', '')).is_valid_request(
        # Validating AWS Lambda's Event Slack Request
        # https://gist.github.com/nitrocode/288bb104893698011720d108e9841b1f
        base64.b64decode(event.get('body', '')).decode("utf-8") if event.get('isBase64Encoded') else event.get('body', ''),
        event.get('headers')
    ):
        raise BadRequestError

    client = boto3.client('stepfunctions')
    step_function_submit_res = client.start_execution(
        stateMachineArn=os.environ.get('STATE_MACHINE_ARN', ''),
        name=f'media-literacy-sf-{datetime.now().strftime("%Y-%H-%M")}-{event.get("body")[-5:]}',
        input=json.dumps({
            'test': 'input'
        })
    )

    loop.run_until_complete(SlackService.send('You sent a slack command!', event))

    return {
        'message': step_function_submit_res,
    }
