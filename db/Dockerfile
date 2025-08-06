FROM postgres:17-alpine

WORKDIR /app

# パッケージ更新＆パッケージのビルドで必要なもの
RUN apk update && apk add --no-cache \
    git \ 
    make \ 
    gcc \
    musl-dev \
    postgresql-dev

# pgvector導入（LTOを完全に無効化）
RUN git clone https://github.com/pgvector/pgvector.git /usr/src/pgvector \
    && cd /usr/src/pgvector \
    && make PG_CONFIG=/usr/local/bin/pg_config OPTFLAGS="-O2 -fno-lto" with_llvm=no \
    && make install PG_CONFIG=/usr/local/bin/pg_config with_llvm=no \
    && cd / 

# EXTENTION登録SQLをコピー
COPY ./docker-entrypoint-initdb.d/initial.sql /docker-entrypoint-initdb.d

# 使わないので削除
RUN apk del \
    git \ 
    make \ 
    gcc \
    musl-dev \
    postgresql-dev

# ファイルも使わないので削除
RUN rm -rf /usr/src/pgvector
