# Synonym Dictionary
유의어를 등록하고 검색하는 API

## Architecture
It depends on elasticsearch

## Usage
Build
- docker `docker-compose up -d`
- kubernetes `kubectl apply -f deployments/app.yml`

Call
- grpc `localhost:9000/synonym`
- rest api `localhost:8080/synonym`

```
1. 입력
curl -XPOST localhost:8080/synonym/ -d'{"name":"삼성전자", "content":"삼전"}'

2. 확인
curl -XGET localhost:8080/synonym/삼전
> return
{ "name": "삼전", "tags": ["삼성전자"] }
```

API
- Create
    - Request
        - body: name, tags
    - Response
        - name, tags
- Get All
    - Request
        - body: name, tags
    - Response
        - name, tags
- Get
    - Request
        - argument: name
        - body: name, tags
    - Response
        - name, tags
- Update
    - Request
        - argument: name
        - body: name, tags
    - Response
        - name, tags
- Delete
    - Request
        - argument: name
        - body: name, tags
    - Response
        - HTTP response
