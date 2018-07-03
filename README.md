# YUKIMOCHI RTMP-VRChat Gateway

## これはなに?
VRChat で RTMP ベースの生配信をインスタンスに対して行うための3つのアプリケーションを同梱しています。

## 同梱物

### VRChat の WebPanel で使用するための RTMP Streaming Player
 - YUKIMOCHI VRC_RTMP Player (`static/endpoint.html`)

クエリ文字列によって再生するストリーミングを選択できるので、 `Set-WebPanelURI` や javascript のページ遷移を用いて、動画の切り替えが可能。
````
location.href="http://example.jp/endpoint.html?addr=*******&streamkey=******"
````

また、 javascript の `audio_control` 関数を実行することで、オーディオの音量とパン（左右の音量割合）を変更できます。

#### audio_control (volume:Number, pan:Number)
  - volume (Number) : 音量を設定する。 `0 ~ 100` までの整数値。初期設定では `100` 。 `100` を超えるとなんかやばい音がする。
  - pan (Number) : 左右の音のふり具合を設定する。 `-100` (左に全振り) ~ `100` (右に全振り) まで設定可能。初期設定では `0` 。

  例：( `Set-WebPanelURI` や ブックマークレット)
  ````
  javascript:setTimeout(audio_control(50,0),0);
  ````
  謎： VRChat の `Set-WebPanelURI` で動かすブックマークレットは、 setTimeout でラップしないと動かないの多いよね。なんでかは知らない。 Vorlon.js のコンソールからなら普通に呼んでも使えた。

### Instance ごとに異なる StreamKey を提供する仕組み
 - YUKIMOCHI VRC_RTMP Instance Key Gen (`static/key_gen.html`)

VRChat で同一ワールドのインスタンスが複数建てられた場合に、それぞれに異なるストリームキーを付与するための仕組み。

クエリ文字列 `sid` が与えられると、日付と混ぜられてストリームキーが生成される。

（ストリームキーは、クエリ文字列 `sid` と UTC での日付で決定される。）

また、 VRC_RTMP Player に生成されたストリームキーを付与してリダイレクトする。

`sid` をインスタンスごとにランダムにする"かつ"、あとから入室したユーザにも同じ値を提供するには、 `VRC_Trigger` の `Randomize` を有効にして、大量に `sid` を埋め込んだ URL を入れておくしかない。

 設定の参考：

 ![#ランダムとは](docs/random.png)

 `Set-WebPanelURI` にエンドポイントとストリームキーを埋め込まない理由は、何らかの変更があったときに、 サーバのアドレスやストリームキーの変換ロジックをワールドの変更なしで変えられないため。

### RTMP サーバ
 - nginx w/ nginx-rtmp-module (`Dockerfile`)

nginx に nginx-rtmp-module というモジュールを追加して、 RTMP サーバを作ります。

同梱の `Dockerfile` または `docker-compose.yml` を使えば、ワンタッチで上記の2アプリを公開する Web サーバと RTMP サーバの両方が自動的に完成します。

## 使い方

 - [RTMP サーバの構築](./docs/build-server.md)
 - [RTMP サーバの更新（執筆中）](./docs/update-server.md)
 - [ワールドの設定](./docs/build-unity.md)
 - [応用:非対称ストリームキー](./docs/appendix.md)
