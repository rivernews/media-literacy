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
def lambda_handler(*args, **kwargs):
    loop.run_until_complete(async_lambda_handler(*args, **kwargs))

async def async_lambda_handler(event, *args):
    await SlackService.send(datetime.now(), 'Pipeline Queue Consumer invoked!')

    asyncio.create_task(SlackService.send(
        datetime.now(), 'Consumer event:\n ``` \n ', event, ' \n ``` \n '
    ))

    # TODO: concurrency control

    for msg in event.get('Records', []):
        await SlackService.send(
            datetime.now(),
            'Processing msg.',
            'ApproximateReceiveCount:',
            msg.get('attributes', {}).get('ApproximateReceiveCount')
        )

        client = boto3.client('stepfunctions')
        step_function_submit_res = client.start_execution(
            stateMachineArn=os.environ.get('STATE_MACHINE_ARN', ''),
            name=f'media-literacy-sf-{datetime.now().strftime("%Y-%H-%M")}-{secrets.token_hex(nbytes=5)}',
            input=json.dumps({
                'test': 'input'
            })
        )
        await SlackService.send(datetime.now(), 'ConsumerLambda: step func started')

    return
