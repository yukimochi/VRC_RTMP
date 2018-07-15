## 応用
### 配信・受信 非対称ストリームキーの設定

受信に使われるストリームキーと配信に用いられるストリームキーが一致するため、放送切断時に乗っ取りを受ける可能性があります。
そのため、非対称ストリームキーを設定します。（配信元IPアドレスの制限をかけている場合を除く）

### ストリームキーごとの配信時キーを設定するため `auth.json` を変更します。
 
  - 変更前
 ````
 {
    "application": {
        "live": {
            "<配信に使いたいストリームキー1>": null,
            "<配信に使いたいストリームキー2>": null,
            "<配信に使いたいストリームキー3>": null
        },
        "multi": {
            "<転送配信に使いたいストリームキー>": null //1つまで。
        }
    }
}
 ````

  - 変更後
 ````
 {
    "application": {
        "live": {
            "<配信に使いたいストリームキー1>": "ストリームキー1の配信用キー",
            "<配信に使いたいストリームキー2>": "ストリームキー2の配信用キー",
            "<配信に使いたいストリームキー3>": "ストリームキー3の配信用キー"
        },
        "multi": {
            "<転送配信に使いたいストリームキー>": "ストリームキーの配信用キー" //1つまで。
        }
    }
}
 ````

### 配信ソフトウェアでは、次のように設定します。

 |URL|ストリームキー|
 |----|----|
 |rtmp://<サーバのIPアドレス>/live または rtmp://<サーバのIPアドレス>/multi|<配信に使いたいストリームキー>?pub_key=<配信用キー>|

※ 受信に使用するアドレスやストリームキーには変更はありません。