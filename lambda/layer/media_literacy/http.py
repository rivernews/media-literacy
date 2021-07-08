from dataclasses import dataclass
from typing import Dict, Optional
import json
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
    def decorator(event, context):
        try:
            res = HttpResponse(200, func(event, context)).build()
        except HttpError as e:
            Logger.error(e)
            res = e.build()
        except Exception as e:
            Logger.error(e)
            res = InternalServerError(str(e)).build()
        
        Logger.info('Response', res)
        
        return res

    return decorator
