# leaf

### api
```
GET http://baidu.com
header k=v
body {"key":"value"}
ret 200

GET http://baidu.com
header k=v,v2
body @a.json
ret 200

GET http://baidu.com
header k=v,v2 k=v3
body @a.json
ret 200
```
### commmand line 
```
leaf api -f api.txt
```