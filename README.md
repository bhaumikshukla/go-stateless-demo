# go-stateless-demo
This is the sample webservice written in golang using Fiber. This explains the stateless behavior. It has REST endpoints which encrypts the text based on given key. This webservice won't store any data.

## Build 

```
cd go-stateless-demo/
go build

```

## Run the application
```
./fibr
```

This will run the application on port 8000

## Endpoints

### Encrypt the text

POST - /encrypt

```
{
    "key":"thisis32bitlongpassphraseimusing", "text": "Hello"
}
```

Sample Response:

```
{
    "enc": "DvXHa9E_CKtad5YOJUMsn0iyZvkN"
}
```

### Decrypt using Key

POST - /decrypt

```
{
    "key":"thisis32bitlongpassphraseimusing", "enc": "DvXHa9E_CKtad5YOJUMsn0iyZvkN"
}
```

Sample response
```
{
    "text": "Hello"
}
```
