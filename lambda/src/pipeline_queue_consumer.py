import os
import json
import asyncio
import secrets
from datetime import datetime
from time import sleep
import boto3
from media_literacy.services.slack_service import SlackService
from media_literacy.logging import Logger

loop = asyncio.get_event_loop()

# Event shape from SQS
# https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
def lambda_handler(event, *args):
    loop.run_until_complete(SlackService.send(datetime.now(), 'Pipeline Queue Consumer invoked!', [msg.get('attributes', {}).get('ApproximateReceiveCount') for msg in event.get('Records', [])]))

    # TODO: concurrency control

    client = boto3.client('stepfunctions')
    step_function_submit_res = client.start_execution(
        stateMachineArn=os.environ.get('STATE_MACHINE_ARN', ''),
        name=f'media-literacy-sf-{datetime.now().strftime("%Y-%H-%M")}-{secrets.token_hex(nbytes=5)}',
        input=json.dumps({
            'test': 'input'
        })
    )
    loop.run_until_complete(SlackService.send(datetime.now(), 'ConsumerLambda: step func started'))

    # wait to slow down this consumer lambda
    Logger.info('Consumer sleeping to slow down')
    sleep(5)

    return
