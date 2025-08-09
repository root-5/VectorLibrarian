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
RUN uv run hf download \
    sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2 \
    onnx/model.onnx \
    tokenizer.json \
    --revision 86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d \
    --local-dir onnx_model

EXPOSE 8000

CMD ["tail", "-f", "/dev/null"]
