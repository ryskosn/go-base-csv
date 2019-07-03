## これはなに

想定ユーザーは EC サイト BASE でお店を出している販売者。

受注データを CSV エクスポートする機能は標準で備わっているが、出力される csv が少し扱いにくかったので API 経由でエクスポートするクライアントを作成した。

## 使用方法

`clientID` と `clientSecret` を書いた `creds.json` を用意する

```json
{
	"clientID" : "xxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	"clientSecret" : "yyyyyyyyyyyyyyyyyyyy"
}
```

### build

`make` を使う。

```sh
# for macOS
$ make build

# for Windows
$ make buildw
```

### token 取得

実行時にオプション `--init=true` を指定する。

```sh
$ ./basecsv --init=true
```

ブラウザで認証ページが表示されるので、BASE のアカウント情報を入力して認証する。

リダイレクト先の URL に含まれている認証コードを引数として再度実行する。

```sh
$ ./basecsv --init=true 849476fuga5b4b50e7fooae0piyo2a6d88d7
token is saved to data/token.json
```

取得したトークンの情報が指定したファイル `data/token.json` に書き込まれるので、`statik` コマンドを実行する。

```sh
$ statik src=data
```

ビルドし直す。

```sh
$ make build
```
