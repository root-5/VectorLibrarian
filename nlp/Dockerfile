# glibc 問題によるため alpine が使えなかった
FROM debian:bookworm-slim

WORKDIR /nlp

# おまじない
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    wget \
    coreutils && \
    rm -rf /var/lib/apt/lists/*

# ONNX モデルとトークナイザーの環境変数を設定
ENV MODEL_NAME="sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2" \
    DOWNLOAD_DIR="/onnx_model" \
    DOWNLOAD_MODEL_PATH="onnx/model_quint8_avx2.onnx" \
    DOWNLOAD_TOKENIZER_PATH="tokenizer.json" \
    SAVED_MODEL_PATH="model.onnx" \
    SAVED_TOKENIZER_PATH="tokenizer.json" \
    MODEL_VERSION="86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d" \
    TENSOR_LENGTH=384

# ONNX モデルとトークナイザーのダウンロード（存在しない場合のみ）
RUN mkdir -p ${DOWNLOAD_DIR}
RUN wget -O "${DOWNLOAD_DIR}/${SAVED_MODEL_PATH}" \
        "https://huggingface.co/${MODEL_NAME}/resolve/${MODEL_VERSION}/${DOWNLOAD_MODEL_PATH}" && \
    wget -O "${DOWNLOAD_DIR}/${SAVED_TOKENIZER_PATH}" \
        "https://huggingface.co/${MODEL_NAME}/resolve/${MODEL_VERSION}/${DOWNLOAD_TOKENIZER_PATH}";

# ライブラリ保存場所
ENV LIBRARY_PATH="/usr/local/lib"

# libonnxruntime.so のインストール（ONNX Runtime 共有ライブラリ）
RUN wget https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz && \
    tar -xz -f onnxruntime-linux-x64-1.22.0.tgz && \
    cp -r onnxruntime-linux-x64-1.22.0/lib/* ${LIBRARY_PATH} && \
    rm -rf onnxruntime-linux-x64-1.22.0 && \
    rm -f onnxruntime-linux-x64-1.22.0.tgz

# libtokenizers.a のインストール（Tokenizers ライブラリ）
RUN wget https://github.com/daulet/tokenizers/releases/latest/download/libtokenizers.linux-amd64.tar.gz && \
    tar -xz -f libtokenizers.linux-amd64.tar.gz && \
    mv libtokenizers.a ${LIBRARY_PATH} && \
    rm -f libtokenizers.linux-amd64.tar.gz

# Go のインストール
RUN wget -O go1.24.4.linux-amd64.tar.gz https://golang.org/dl/go1.24.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz && \
    rm go1.24.4.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

# CGOを有効化（重要）
ENV CGO_ENABLED=1

# Go mod の初期化とONNX Runtime Goライブラリのインストール
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8000

CMD ["sleep", "infinity"]
