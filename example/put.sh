
#!/bin/sh
url='http://pttalertor.dinolai.com/users/'
user='icbruce1'

echo '
{
   "profile":{
   "account": "icbruce",
   "email":"ic09272002@mail.com"
   },
   "subscribes":[]
}' > json.txt


curl -X PUT -H "Content-Type: application/json" -T json.txt $url/$user
