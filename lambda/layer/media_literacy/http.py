from dataclasses import dataclass
from typing import Union, Dict, Optional
import json
import base64
from media_literacy.logging import Logger
from media_literacy.serializer import JSONResponseEncoder

@dataclass
class HttpResponse:
    status_code: Optional[int] = None
    body: Optional[Dict] = None

    def build(self):
        return {
            'statusCode': self.status_code,
            'headers': {
                'Content-Type': 'application/json'
            },
            'body': json.dumps(self.body, cls=JSONResponseEncoder)
        }

class HttpError(HttpResponse, Exception):
    def __init__(self, message=None, status_code=None):
        super().__init__(status_code=status_code, body={
            'message': repr(self) if not message else message
        })

class BadRequestError(HttpError):
    def __init__(self, message=None, status_code=None):
        super().__init__(message=message, status_code=400 if not status_code else status_code)

class InternalServerError(HttpError):
    def __init__(self, message=None, status_code=None):
        super().__init__(message=message, status_code=500 if not status_code else status_code)

def handle_exception(func):
    def decorator(*args, **kwargs):
        try:
            res = HttpResponse(200, func(*args, **kwargs)).build()
        except HttpError as e:
            Logger.error(e)
            res = e.build()
        except Exception as e:
            Logger.error(e)
            res = InternalServerError(str(e)).build()
        
        Logger.info('Response', res)
        
        return res

    return decorator


class APIGatewayRequest:
    # the raw body string from event
    _body: str

    # the parsed body based on content-type; if empty then parse to empty dict {}; if cannot parse then default to string as-is from event body
    body: Union[dict, str]

    def __init__(self, event: dict):
        self._event = event
        self.headers = event.get('headers', {})

        # parse body
        content_type = self.headers.get('content-type', {})
        self._body = event.get('body', '')
        if not self._body:
            self.body = {}
        elif content_type == 'application/x-www-form-urlencoded':
            self.body = base64.b64decode(self._body).decode("utf-8") if event.get('isBase64Encoded') else self._body
        elif content_type == 'application/json':
            self.body = json.loads(event.get('body', ''))
        else:
            self.body = self._body

    @staticmethod
    def use_request(func):
        def decorator(event, context):
            request = APIGatewayRequest(event)
            return func(request, context)
        return decorator
    
    def __str__(self):
        return str({
            **self._event,
            'parsedBody': self.body,
        })
