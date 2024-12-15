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

## ライブラリ追加方法
```
docker compose -f docker-compose.dev.yml run --rm app-dev go get github.com/nfnt/resize
```
