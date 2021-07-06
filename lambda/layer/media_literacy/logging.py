import os
import logging

class _Logger:
    def __init__(self):
        self.loglevel = os.environ.get('LOGLEVEL', 'ERROR').upper()
        self.logger = logging.getLogger()
        numeric_level = getattr(logging, self.loglevel, None)
        if not isinstance(numeric_level, int):
            raise ValueError('Invalid log level: %s' % self.loglevel)
        self.logger.setLevel(numeric_level)
    
    def debug(self, *args):
        message = '🐞 ' + ' '.join([str(a) for a in args])
        self.logger.debug(message)
    
    def warn(self, *args):
        message = '🍊 ' + ' '.join([str(a) for a in args])
        self.logger.warning(message)
    
    def info(self, *args):
        message = '💬 ' + ' '.join([str(a) for a in args])
        self.logger.info(message)

    def error(self, *args):
        message = '🔴 ' + ' '.join([str(a) for a in args])
        self.logger.error(message)
Logger = _Logger()
