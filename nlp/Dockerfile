FROM python:3.13-slim-bookworm

WORKDIR /nlp

# おまじない
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    git \
    curl \
    coreutils \
    && rm -rf /var/lib/apt/lists/*

# uv のインストール
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:$PATH"

# 依存関係ファイルをコピー
COPY pyproject.toml uv.lock ./

# 依存関係を同期（この段階で全パッケージがインストールされる）
RUN uv sync --frozen

# ONNX モデルとトークナイザーの環境変数を設定
ENV MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2" \
    DOWNLOAD_MODEL_PATH="onnx/model_quint8_avx2.onnx" \
    DOWNLOAD_TOKENIZER_PATH="tokenizer.json" \
    MODEL_VERSION="86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d" \
    DOWNLOAD_DIR="onnx_model" \
    SAVED_MODEL_PATH="onnx/model.onnx" \
    TENSOR_LENGTH=384

# ONNX モデルとトークナイザーのダウンロード
RUN uv run hf download \
    ${MODEL_NAME} \
    ${DOWNLOAD_MODEL_PATH} \
    ${DOWNLOAD_TOKENIZER_PATH} \
    --revision ${MODEL_VERSION} \
    --local-dir ${DOWNLOAD_DIR}

# ダウンロードキャッシュの削除
RUN rm -rf /root/.cache/huggingface

# Go のインストール
RUN curl -LsSf https://golang.org/dl/go1.24.5.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH="$PATH:/usr/local/go/bin"

# Go mod の初期化とONNX Runtime Goライブラリのインストール
COPY go.mod go.sum ./
RUN go mod download

# ライブラリ保存場所
ENV LIBRARY_PATH="/usr/local/lib"

# libonnxruntime.so のインストール（ONNX Runtime 共有ライブラリ）
RUN curl -LsSf https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz \
    | tar -xz && \
    cp -r onnxruntime-linux-x64-1.22.0/lib/* ${LIBRARY_PATH} && \
    rm -rf onnxruntime-linux-x64-1.22.0 && \
    rm -f onnxruntime-linux-x64-1.22.0.tgz

# libtokenizer.a のインストール（Rust製トークナイザーライブラリ）
RUN curl -LsSf https://github.com/daulet/tokenizers/releases/latest/download/libtokenizers.linux-amd64.tar.gz \
    | tar -xz && \
    mv libtokenizers.a ${LIBRARY_PATH} && \
    rm -f libtokenizers.linux-amd64.tar.gz

EXPOSE 8000

CMD ["tail", "-f", "/dev/null"]
