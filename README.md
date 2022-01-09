### In Progress


```
openssl req -new -x509 -subj "/CN=valkon.default.svc"  -addext "subjectAltName = DNS:valkon.default.svc" -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365
```



![Design](diagrams/arch.jpg)