twitter-favorites-api
===============

特定のTwitterアカウントのお気に入りに入っている画像のURLをランダムで返すだけのAPI。
GAE/Goで動かします。

## Usage

※雑な使い方説明なので。気が向いたら整理します。

### Setup

```
make init
make install
```

### app.yaml

1. app_example.yamlをコピーしてapp.yamlを作成
2. applicationを自分のアプリ名にする
3. TWITTER_API_KEYとTWITTER_API_SECRETを設定する
4. TWITTER_IDを設定（お気に入りを取り込むやつ）
5. CACHE_COUNTを設定（何件までオンメモリにキャッシュするかの設定）

### ローカル起動

```
make serve
```

### GAE/Goデプロイ

```
make deploy
```
