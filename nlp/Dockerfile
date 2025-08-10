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

# ONNX モデルとトークナイザーのダウンロード
ENV MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2" \
    ONNX_MODEL_PATH="onnx/model_quint8_avx2.onnx" \
    TOKENIZER_PATH="tokenizer.json" \
    MODEL_VERSION="86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d" \
    DOWNLOAD_DIR="onnx_model"
RUN uv run hf download \
    ${MODEL_NAME} \
    ${ONNX_MODEL_PATH} \
    ${TOKENIZER_PATH} \
    --revision ${MODEL_VERSION} \
    --local-dir ${DOWNLOAD_DIR}

# ダウンロードキャッシュの削除
RUN rm -rf /root/.cache/huggingface

# Go のインストール
RUN curl -LsSf https://golang.org/dl/go1.24.5.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH="$PATH:/usr/local/go/bin"

# Go mod の初期化とONNX Runtime Goライブラリのインストール
COPY go.mod go.sum ./
RUN go mod tidy

# libonnxruntime.so のインストール（ONNX Runtime 共有ライブラリ）
RUN curl -LsSf https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz \
    | tar -xz -C /usr/local/lib && \
    ln -s /usr/local/lib/onnxruntime-linux-x64-1.22.0/libonnxruntime.so /usr/local/lib/libonnxruntime.so

EXPOSE 8000

CMD ["tail", "-f", "/dev/null"]
