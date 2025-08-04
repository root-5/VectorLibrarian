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

# Go のインストール
RUN curl -LsSf https://golang.org/dl/go1.24.5.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH="$PATH:/usr/local/go/bin"

CMD ["tail", "-f", "/dev/null"]
