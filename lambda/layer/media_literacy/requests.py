import aiohttp
from media_literacy.logging import Logger

# Response Object Reference
# https://docs.aiohttp.org/en/stable/client_reference.html#response-object

class _AsyncRequest:
    def __init__(self, headers=None):
        self._session = None
    
    async def post(self, url, params=None, data=None, headers=None, raise_for_status=True):
        if not self._session:
            self._session = aiohttp.ClientSession()

        async with self._session.post(url, params=None, json=data, headers=headers, raise_for_status=raise_for_status) as async_res:
            Logger.debug(async_res.method, async_res.status, async_res.url, await async_res.text())
            if async_res.headers.get('content-type') == 'application/json':
                return await async_res.json()
            else:
                return await async_res.text()


AsyncRequest = _AsyncRequest()