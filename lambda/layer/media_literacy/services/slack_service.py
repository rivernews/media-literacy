import os
from typing import Union, Dict
from media_literacy.http import InternalServerError
from media_literacy.requests import AsyncRequest

class _SlackService:

    def __init__(self):
        self._webhook_url = os.environ.get('SLACK_POST_WEBHOOK_URL', '')
        if not self._webhook_url:
            raise InternalServerError(message='Slack post webhook url not configured')
    
    async def send(self, *messages):
        await AsyncRequest.post(self._webhook_url, data={
            "text": ' '.join([str(_message) for _message in messages])
        })

SlackService = _SlackService()
