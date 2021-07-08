from media_literacy.services.slack_service import SlackService


def test_async_request_basic_2(event_loop):
    res = event_loop.run_until_complete(SlackService.send('Pytest OK'))
    