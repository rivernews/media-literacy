import os
from enum import Enum
from typing import Union, Dict
from media_literacy.http import InternalServerError
from media_literacy.requests import AsyncRequest

ENV = os.environ.get("ENV", "production")


class SlackCommandMantra(str, Enum):
    FETCH_LANDING = '/click'
    FETCH_LANDING_STORIES = '/getallstories'

    def __str__(self):
        return self.value


class _SlackService:

    def __init__(self):
        self._webhook_url = os.environ.get('SLACK_POST_WEBHOOK_URL', '')
        if not self._webhook_url:
            raise InternalServerError(message='Slack post webhook url not configured')

    async def send(self, *messages):
        await AsyncRequest.post(self._webhook_url, data={
            "text": f'[{ENV}] ' + ' '.join([str(_message) for _message in messages])
        })

    @staticmethod
    def parse_command(request_body: dict):
        command = request_body.get('command')
        text = request_body.get('text')



SlackService = _SlackService()
