from media_literacy.requests import AsyncRequest

def test_async_request_basic_2(event_loop):
    res = event_loop.run_until_complete(AsyncRequest.post('https://media-literacy.api.shaungc.com/slack/command', raise_for_status=False))
    assert res is not None