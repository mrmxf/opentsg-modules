'''
 pip install s3fs
import s3fs
fs = s3fs.S3FileSystem(anon=True)
open a basic file at the destination and write it
hope for the correct permissions
https://s3fs.readthedocs.io/en/latest/api.html?highlight=s3filesystem#s3fs.core.S3FileSystem
'''

import s3fs
import os
ID = os.getenv('AWS_ACCESS_KEY_ID')
KEY= os.environ.get('AWS_SECRET_ACCESS_KEY')
REGION = os.getenv('AWS_DEFAULT_REGION')
print(ID, KEY,REGION)


fs = s3fs.S3FileSystem(anon=False, key=ID, secret = KEY)#client_kwargs={'keystring':ID,'secretstring':KEY})

x = fs.ls('mmh-cache/bot-tlh/dev')

fs.mkdir('mmh-cache/bot-tlh/dev/opentpg')
print(fs.exists('mmh-cache/bot-tlh/dev/opentpg'))
fs.put('clogrc/init.md','mmh-cache/bot-tlh/dev/opentpg/init.md')
print(x)
