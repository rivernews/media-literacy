from dataclasses import dataclass
from typing import Dict, Optional
import json

@dataclass
class HttpResponse:
    status_code: Optional[int] = None
    body: Optional[Dict] = None

    def build(self):
        return {
            'statusCode': self.status_code,
            'body': json.dumps(self.body)
        }

class HttpError(HttpResponse, Exception):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.body = {
            'message': repr(self)
        }

class BadRequestError(HttpError):
    def __init__(self, *args, **kwargs):
        if 'status_code' not in kwargs:
            kwargs['status_code'] = 400
        super().__init__(*args, **kwargs)
        self.status_code = 400
