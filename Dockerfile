FROM alpine:3.8

ENV NGINX_VERSION 1.15.5
ENV NGINX_RTMP_MODULE_VERSION 1.2.1

RUN adduser -s /sbin/nologin -D -H nginx

RUN apk add --update --no-cache \
    ca-certificates \
    build-base \
    libressl \
    libressl-dev \
    pcre-dev && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

RUN mkdir /tmp/build && \
    cd /tmp/build && \
    wget -O nginx-${NGINX_VERSION}.tar.gz https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
    tar -zxf nginx-${NGINX_VERSION}.tar.gz && \
    rm -f nginx-${NGINX_VERSION}.tar.gz && \
    wget -O nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz https://github.com/arut/nginx-rtmp-module/archive/v${NGINX_RTMP_MODULE_VERSION}.tar.gz && \
    tar -zxf nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz && \
    rm -f nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz && \
    cd /tmp/build/nginx-${NGINX_VERSION} && \    
    ./configure \
    --sbin-path=/usr/local/sbin/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --http-log-path=/var/log/nginx/access.log \
    --error-log-path=/var/log/nginx/error.log \
    --pid-path=/var/run/nginx/nginx.pid \
    --lock-path=/var/lock/nginx/nginx.lock \
    --http-client-body-temp-path=/tmp/http-client-body \
    --user=nginx --group=nginx \
    --with-http_ssl_module \
    --with-http_gzip_static_module \
    --without-http_access_module \
    --without-http_auth_basic_module \
    --without-http_autoindex_module \
    --without-http_browser_module \
    --without-http_empty_gif_module \
    --without-http_fastcgi_module \
    --without-http_geo_module \
    --without-http_gzip_module \
    --without-http_limit_conn_module \
    --without-http_limit_req_module \
    --without-http_map_module \
    --without-http_memcached_module \
    --without-http_proxy_module \
    --without-http_referer_module \
    --without-http_scgi_module \
    --without-http_split_clients_module \
    --without-http_upstream_hash_module \
    --without-http_upstream_ip_hash_module \
    --without-http_upstream_keepalive_module \
    --without-http_upstream_least_conn_module \
    --without-http_upstream_zone_module \
    --without-http_userid_module \
    --without-http_uwsgi_module \
    --without-http-cache \
    --without-mail_imap_module \
    --without-mail_pop3_module \
    --without-mail_smtp_module \
    --without-stream_access_module \
    --without-stream_limit_conn_module \
    --without-stream_upstream_hash_module \
    --without-stream_upstream_least_conn_module \
    --without-stream_upstream_zone_module \
    --with-threads \
    --add-module=/tmp/build/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION} && \
    make -j $(getconf _NPROCESSORS_ONLN) && \
    make install && \
    mkdir /var/lock/nginx /tmp/http-client-body && \
    rm -rf /tmp/build && \
    apk del --purge build-base libressl-dev && \
    rm -rf /var/cache/apk/*

EXPOSE 80 1935

CMD [ "nginx" ]
