import os
import json
import base64
import aiohttp
from slack_sdk.signature import SignatureVerifier
from media_literacy.http import HttpResponse, HttpError, BadRequestError


def lambda_handler(event, context):
    try:
        print('Hello there!', event, context)

        if not SignatureVerifier(os.environ.get('SLACL_SIGNING_SECRET', '')).is_valid_request(
            base64.b64decode(event.get('body', '')).decode("utf-8") if event.get('isBase64Encoded') else event.get('body', ''),
            event.get('headers')
        ):
            raise BadRequestError

    except HttpError as e:
        return e.build()

    return HttpResponse(200, {
            'message': 'hello',
            'event': event
        }).build()
