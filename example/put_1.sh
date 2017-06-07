
#!/bin/sh
url='http://pttalertor.dinolai.com/users'
user='icbruce'

echo '
{
   "enable":true, 
   "profile":{
   "account": "icbruce",
   "email":"ic09272002@mail.com"
   },
   "subscribes":[
      {
         "board":"gossiping",
         "keywords":["柯文哲","柯P","柯Ｐ"]
      },
      {
         "board":"lol",
         "keywords":["樂透"]
      }
   ]
}'  > json.txt

curl -X PUT -H "Content-Type: application/json" -T json.txt $url/$user
