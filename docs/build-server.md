## RTMP サーバの構築
 1. Linux ベースのサーバを用意し、git, [Docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/) および [Docker-Compose](https://docs.docker.com/compose/install/#install-compose) をインストールします。
 2. ホームディレクトリで以下の通りのコマンドを実行します。
 ````
git clone https://github.com/yukimochi/VRC_RTMP.git //コードのダウンロード
cd VRC_RTMP //ディレクトリの移動
sudo docker-compose build //nginxのセットアップをします。
sudo docker-compose up -d //サーバが起動します。
 ````

 ※ サーバの停止コマンド `sudo docker-compose stop`

 ※ サーバの起動ログ表示コマンド `sudo docker-compose logs` //正常起動した場合何も表示されない。

 ※ httpのアクセスログの場所 `./logs/access.log`
 
 ※ rtmpとエラーのログの場所 `./logs/error.log` //rtmp のログが error の方に出力される原因は不明。

 3. `http://<サーバのIPアドレス>/key_gen.html` にアクセスすると、`static/key_gen.html` を見ることができるか確認する。
 4. `rtmp://<サーバのIPアドレス>/live` に RTMP 配信（OBS-Studio など使用）できるか確認する。
 5. -1 セキュリティを確保するため、`nginx.conf` を変更します。

  - 変更前
 ````nginx.conf
rtmp {
  server {
    listen 1935;

    application live {
      live on;
      record off;
      allow publish all; //誰でも、配信が可能になっている。（上が優先度大）
      deny publish all; //この設定は、無視される。
      allow play all;

      on_publish http://auth:3000/auth;
    }

    # Comment out to enable relay to YouTube Live etc...
    #application multi {
    #  live on;
    #  record off;
    #  allow publish all;
    #  deny publish all;
    #  allow play all;

    #  on_publish http://auth:3000/auth;

    #  push rtmp://a.rtmp.youtube.com/live*/***********
    #  push rtmp://live-tyo.twitch.tv/app/***********
    }
  }
}
 ````

  - 変更後
 ````nginx.conf 
rtmp {
  server {
    listen 1935;

    application live {
      live on;
      record off;
      allow publish <配信PCのグローバルIPアドレス>; //自分のIPアドレス以外の配信は、許可されない。
      deny publish all; //上の条件に満たない場合、この条件によって配信が禁止される。
      allow play all;
    }

    # Comment out to enable relay to YouTube Live etc...
    #application multi {
    #  live on;
    #  record off;
    #  allow publish all;
    #  deny publish all;
    #  allow play all;

    #  on_publish http://auth:3000/auth;

    #  push rtmp://a.rtmp.youtube.com/live*/***********
    #  push rtmp://live-tyo.twitch.tv/app/***********
    }
  }
}

 ````

 ※ 設定変更後は、 `sudo docker-compose stop` と `sudo docker-compose start` で必ず再起動します。

 5. -2 配信状況の統計情報を表示できるようにしたい場合、コメントアウトされた以下の部分を修正します。
  - 変更前
 ````
    # Comment out to enable statistics.
    #location /statistics/view {
    #  rtmp_stat all;
    #  rtmp_stat_stylesheet /stat.xsl;
    #}
 ````

  - 変更後
 ````
    # Comment out to enable statistics.
    location /statistics/view {
      rtmp_stat all;
      rtmp_stat_stylesheet /stat.xsl;
    }
 ````

 設定適応後は、 `http://<サーバのIPアドレス>/statistics/view` にアクセスすると、以下のような統計情報を確認できます。

 ![RTMP Statisstics](./stat.png)

 5. -3 YouTube Live など他の配信サービスに配信を転送したい場合、 `rtmp://<サーバのIPアドレス>/multi` へ配信するようにします。 それに伴い、以下の通りに設定を変更します。

  - 変更前
 ```` 
    # Comment out to enable relay to YouTube Live etc...
    #application multi {
    #  live on;
    #  record off;
    #  allow publish all;
    #  deny publish all;
    #  allow play all;

    #  on_publish http://auth:3000/auth;

    #  push rtmp://a.rtmp.youtube.com/live*/***********
    #  push rtmp://live-tyo.twitch.tv/app/***********
    }
 ````

  - 変更後
 ````
    # Comment out to enable relay to YouTube Live etc...
    application multi {
      live on;
      record off;
      allow publish all;
      deny publish all;
      allow play all;

      on_publish http://auth:3000/auth;

      push rtmp://a.rtmp.youtube.com/live*/*********** //あなたの配信したいプラットフォームのRTMPアドレスを指定します。
      push rtmp://live-tyo.twitch.tv/app/*********** //あなたの配信したいプラットフォームのRTMPアドレスを指定します。
    }
 ````

 6. 配信に利用できるストリームキーを設定するため `auth.json` を変更します。
 
  - 変更前
 ````
{
    "application": {
        "live": {
            "stream": null,
            "sub_stream": null
        },
        "multi": {
            "multi_stream": null
        }
    }
}
 ````

  - 変更後
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

 7. 配信ソフトウェアでは、次のように設定します。

 |URL|ストリームキー|
 |----|----|
 |rtmp://<サーバのIPアドレス>/live または rtmp://<サーバのIPアドレス>/multi|<配信に使いたいストリームキー>|
