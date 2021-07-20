import os
import json
import asyncio
import secrets
from datetime import datetime
import boto3
from media_literacy.services.slack_service import SlackService

loop = asyncio.get_event_loop()

# Event shape from SQS
# https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
def lambda_handler(event, *args):
    loop.run_until_complete(SlackService.send('Pipeline Queue Consumer invoked!', event))

    # TODO: concurrency control

    client = boto3.client('stepfunctions')
    step_function_submit_res = client.start_execution(
        stateMachineArn=os.environ.get('STATE_MACHINE_ARN', ''),
        name=f'media-literacy-sf-{datetime.now().strftime("%Y-%H-%M")}-{secrets.token_hex(nbytes=5)}',
        input=json.dumps({
            'test': 'input'
        })
    )
    loop.run_until_complete(SlackService.send('ConsumerLambda: step func started', step_function_submit_res))

    return
