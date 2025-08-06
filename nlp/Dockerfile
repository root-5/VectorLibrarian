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

# uvicornのインストール
RUN uv tool install "uvicorn[standard]"

CMD ["tail", "-f", "/dev/null"]
