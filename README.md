# YUKIMOCHI RTMP-VRChat Gateway

## これはなに?
かつて VRChat に存在した `WebPanel` というオブジェクトを対象にした RTMP サーバ構築用のアプリケーションセットです。
現在は、簡易 RTMP サーバ構築用スクリプトとして利用することができます。ただしメンテナンスはされません。
[Docker Hub](https://hub.docker.com/r/yukimochi/vrc_rtmp/) には、 nginx 1.15.5 で構築された amd64, arm, aarch64 向けのイメージが永続的に提供されます。

## 生配信を VRChat で再生したいのですがどうすればよいですか。
以下のサービス（アプリケーション）と著名なビデオストリーミングサービスを併用して実現します。
ただし、遅延は 5~15 秒にも及ぶと思われます。"超"低遅延ストリーミングの実現には、専用のソリューションをご検討ください。

[YUKIMOCHI VRChat HLS Bridge](https://github.com/yukimochi/VRC_HLS)

## かつての Readme はどこにありますか？（私は過去の時間からこのページを見ています。）
[こちら](./docs/readme.md) から参照できます。
