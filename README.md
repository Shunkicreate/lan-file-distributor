Dockerfileを使ってビルドして、デプロイする

# 開発環境
## ローカルでビルドする
```
docker compose -f docker-compose.dev.yml up --build
```

## ローカルで実行する
```
docker compose -f docker-compose.dev.yml up
```


# 本番環境
## ローカルでビルドする
```
docker compose -f docker-compose.prod.yml up --build
```

## ローカルで実行する
```
docker compose -f docker-compose.prod.yml up
```
