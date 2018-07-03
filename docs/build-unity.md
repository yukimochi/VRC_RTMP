## ワールドの設定
 8. -1 （全てのワールドで1つの配信を見せる場合）ワールドの Web Panel にプレーヤーの URL を設定します。

 プレーヤーは、 `addr` 変数に URL エンコードしたサーバアドレス、 `streamkey` にストリームキーを設定します。
 （ここで決定したストリームキーで配信することになります。）

 `http://<サーバのIPアドレス>/endpoint.html?addr=rtmp%3A%2F%2F<サーバのIPアドレス>%2Flive&streamkey=<ストリームキー>`

 `endpoint.html` は、 1940*1100 で表示されることを期待しています。

 ※ 転送配信側 ( `/multi` ) を視聴する場合、以下のプレーヤーの URL を設定します。

 `http://<サーバのIPアドレス>/endpoint.html?addr=rtmp%3A%2F%2F<サーバのIPアドレス>%2Fmulti&streamkey=<ストリームキー>`

 設定の参考：
 
 ![Unity設定例](./unity.png)

 8. -2（インスタンスごとに別の配信を見せたいとき） ＜略＞

 9. 色合いがおかしい件の修正。

 Web Panel では、色の発色がおかしいという現象が[見つかっています](http://uuupa.hatenablog.com/entry/2018/04/05/003936)。

 Duplicate Screen で、 [yukimochi/WebPanel-Shaders](https://github.com/yukimochi/WebPanel-Shaders
) などの代替シェーダを適用します。
