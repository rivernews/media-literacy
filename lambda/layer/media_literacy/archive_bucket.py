import os
import boto3

s3 = boto3.resource('s3')

class _ArchiveBucket:
    s3 = s3
    bucket_name = os.environ.get('S3_ARCHIVE_BUCKET', '')
    bucket = None

    def __init__(self):
        try:
            self.s3.meta.client.head_bucket(Bucket=self.bucket_name)
        except self.s3.meta.client.exceptions.NoSuchBucket as e:
            raise Exception(f'Archive bucket `{self.bucket_name}` does not exist, did you pass in env var `S3_ARCHIVE_BUCKET`? {e}')

        self.bucket = self.s3.Bucket(self.bucket_name)

    def exist(self, key) -> bool:
        objects = list(self.bucket.objects.filter(Prefix=key))
        if any([obj.key == key for obj in objects]):
            return True
        else:
            return False

ArchiveBucket = _ArchiveBucket()
