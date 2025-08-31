# glibc 問題によるため alpine が使えない
FROM debian:bookworm-slim AS builder

WORKDIR /nlp

# 必要なツールをインストール
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    coreutils \
    git && \
    rm -rf /var/lib/apt/lists/*

# ONNX モデルとトークナイザーの環境変数を設定
ENV MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2" \
    DOWNLOAD_DIR="/onnx_model" \
    DOWNLOAD_MODEL_PATH="onnx/model_quint8_avx2.onnx" \
    DOWNLOAD_TOKENIZER_PATH="tokenizer.json" \
    SAVED_MODEL_PATH="model.onnx" \
    SAVED_TOKENIZER_PATH="tokenizer.json" \
    MODEL_VERSION="86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d" \
    MAX_TOKEN_LENGTH=512 \
    OVERLAP_TOKEN_LENGTH=103 \
    TENSOR_LENGTH=384

# ONNX モデルとトークナイザーのダウンロード
RUN mkdir -p ${DOWNLOAD_DIR} && \
    curl -fsSL -o "${DOWNLOAD_DIR}/${SAVED_MODEL_PATH}" \
        "https://huggingface.co/${MODEL_NAME}/resolve/${MODEL_VERSION}/${DOWNLOAD_MODEL_PATH}" && \
    curl -fsSL -o "${DOWNLOAD_DIR}/${SAVED_TOKENIZER_PATH}" \
        "https://huggingface.co/${MODEL_NAME}/resolve/${MODEL_VERSION}/${DOWNLOAD_TOKENIZER_PATH}"

# ライブラリパスを設定
ENV LIBRARY_PATH="/libraries"
RUN mkdir -p ${LIBRARY_PATH}

# libonnxruntime.so のインストール
RUN curl -fsSL -o onnxruntime-linux-x64-1.22.0.tgz https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz && \
    tar -xz -f onnxruntime-linux-x64-1.22.0.tgz && \
    cp -r onnxruntime-linux-x64-1.22.0/lib/* ${LIBRARY_PATH} && \
    rm -rf onnxruntime-linux-x64-1.22.0 && \
    rm -f onnxruntime-linux-x64-1.22.0.tgz

# libtokenizers.a のインストール
RUN curl -fsSL -o libtokenizers.linux-amd64.tar.gz https://github.com/daulet/tokenizers/releases/latest/download/libtokenizers.linux-amd64.tar.gz && \
    tar -xz -f libtokenizers.linux-amd64.tar.gz && \
    mv libtokenizers.a ${LIBRARY_PATH} && \
    rm -f libtokenizers.linux-amd64.tar.gz

# Go のインストール
RUN curl -fsSL -o /tmp/go.tar.gz https://golang.org/dl/go1.24.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

# Go mod の初期化と依存関係のダウンロード
ENV GOPROXY=https://proxy.golang.org
COPY . .

# ライブラリのダウンロード
RUN go mod download

CMD ["go", "run", "main.go"]
