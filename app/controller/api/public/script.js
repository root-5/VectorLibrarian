// index.html
if (window.location.pathname === '/') {
  document.getElementById('search-input').addEventListener('keydown', function (event) {
    if (event.key === 'Enter') {
      document.getElementById('search-button').click();
    }
  });

  document.getElementById('search-button').addEventListener('click', async function () {
    const query = document.getElementById('search-input').value;
    const resultsDiv = document.getElementById('search-results');
    const loadingIndicator = document.getElementById('loading-indicator');
    resultsDiv.innerHTML = ''; // Clear previous results

    if (query.trim() === '') {
      resultsDiv.innerHTML = '<p>検索キーワードを入力してください。</p>';
      return;
    }

    loadingIndicator.style.display = 'block'; // Show loading indicator

    try {
      const response = await fetch(`/search?q=${query}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const searchResults = await response.json();

      if (searchResults.length > 0) {
        searchResults.forEach((result) => {
          const resultItem = document.createElement('div');
          let description;
          if (result.description || result.description == '--') {
            description = result.description;
          } else {
            description = result.keywords;
          }
          resultItem.classList.add('result-item');
          resultItem.innerHTML = `
                        <h3><a href="https://${result.domain}${result.path}" target="_blank" rel="noopener noreferrer">${result.title}</a></h3>
                        <p>${result.description}</p>
                    `;
          resultsDiv.appendChild(resultItem);
        });
      } else {
        resultsDiv.innerHTML = '<p>検索結果が見つかりませんでした。</p>';
      }
    } catch (error) {
      console.error('検索中にエラーが発生しました:', error);
      resultsDiv.innerHTML = '<p>検索中にエラーが発生しました。もう一度お試しください。</p>';
    } finally {
      loadingIndicator.style.display = 'none'; // Hide loading indicator
    }
  });
}

// chat.html
if (window.location.pathname === '/chat') {
  const chatMessages = document.getElementById('chat-messages');
  const chatInput = document.getElementById('chat-input');
  const sendButton = document.getElementById('send-button');

  // Markdownを簡易的にHTMLに変換する関数
  function markdownToHtml(text) {
    // エスケープ処理
    const escapeHtml = (str) => {
      const div = document.createElement('div');
      div.textContent = str;
      return div.innerHTML;
    };

    let html = text;

    // 行ごとに処理
    const lines = html.split('\n');
    const processedLines = [];
    let listStack = []; // ネストレベルを管理するスタック

    for (let i = 0; i < lines.length; i++) {
      let line = lines[i];
      const trimmedLine = line.trim();

      // 空行の処理
      if (trimmedLine === '') {
        // リスト中の空行はリストを閉じる
        while (listStack.length > 0) {
          processedLines.push('</ul>');
          listStack.pop();
        }
        continue;
      }

      // 区切り線（---）をスキップ
      if (trimmedLine === '---' || trimmedLine.match(/^-{3,}$/)) {
        while (listStack.length > 0) {
          processedLines.push('</ul>');
          listStack.pop();
        }
        continue;
      }

      // 見出し (# ## ### ####)
      if (line.match(/^(#{1,4})\s+(.+)$/)) {
        const match = line.match(/^(#{1,4})\s+(.+)$/);
        const level = match[1].length;
        const content = escapeHtml(match[2]);
        while (listStack.length > 0) {
          processedLines.push('</ul>');
          listStack.pop();
        }
        processedLines.push(`<h${level}>${content}</h${level}>`);
      }
      // 箇条書き (- または * で始まる、インデント対応)
      else if (line.match(/^(\s*)([\-\*])\s+(.+)$/)) {
        const match = line.match(/^(\s*)([\-\*])\s+(.+)$/);
        const indent = match[1].length;
        let content = match[3];

        // インデントレベルを計算（4スペース = 1レベル）
        const currentLevel = Math.floor(indent / 4);

        // 太字処理
        content = content.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
        content = content.replace(/__(.+?)__/g, '<strong>$1</strong>');
        content = escapeHtml(content)
          .replace(/&lt;strong&gt;/g, '<strong>')
          .replace(/&lt;\/strong&gt;/g, '</strong>');

        // リストのネストレベル調整
        while (listStack.length > currentLevel + 1) {
          processedLines.push('</ul>');
          listStack.pop();
        }

        if (listStack.length === currentLevel) {
          processedLines.push('<ul>');
          listStack.push(currentLevel);
        }

        processedLines.push(`<li>${content}</li>`);
      }
      // 通常のテキスト
      else {
        while (listStack.length > 0) {
          processedLines.push('</ul>');
          listStack.pop();
        }
        // 太字処理 (**text** または __text__)
        line = line.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
        line = line.replace(/__(.+?)__/g, '<strong>$1</strong>');
        line = escapeHtml(line)
          .replace(/&lt;strong&gt;/g, '<strong>')
          .replace(/&lt;\/strong&gt;/g, '</strong>');
        processedLines.push(`<p>${line}</p>`);
      }
    }

    // 最後にリストが閉じられていない場合
    while (listStack.length > 0) {
      processedLines.push('</ul>');
      listStack.pop();
    }

    // 改行なしで結合
    return processedLines.join('');
  }

  // ストリーミングメッセージを作成
  function createStreamingMessage() {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message', 'assistant-message');

    const roleDiv = document.createElement('div');
    roleDiv.classList.add('message-role');
    roleDiv.textContent = 'アシスタント';

    const contentDiv = document.createElement('div');
    contentDiv.classList.add('message-content');
    contentDiv.id = 'streaming-content';

    messageDiv.appendChild(roleDiv);
    messageDiv.appendChild(contentDiv);
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;

    return { messageDiv, contentDiv };
  }

  // ソース情報を追加
  function addSourcesToMessage(messageDiv, sources) {
    const sourcesDiv = document.createElement('div');
    sourcesDiv.classList.add('sources');

    const sourcesTitle = document.createElement('div');
    sourcesTitle.classList.add('sources-title');
    sourcesTitle.textContent = '参照元:';
    sourcesDiv.appendChild(sourcesTitle);

    sources.forEach((source, index) => {
      const sourceItem = document.createElement('div');
      sourceItem.classList.add('source-item');

      const sourceLink = document.createElement('a');
      sourceLink.classList.add('source-link');
      sourceLink.href = `https://${source.domain}${source.path}`;
      sourceLink.target = '_blank';
      sourceLink.rel = 'noopener noreferrer';
      sourceLink.textContent = `${index + 1}. ${source.title}`;

      sourceItem.appendChild(sourceLink);
      sourcesDiv.appendChild(sourceItem);
    });

    messageDiv.appendChild(sourcesDiv);
  }

  // ユーザーメッセージを追加
  function addUserMessage(content) {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message', 'user-message');

    const roleDiv = document.createElement('div');
    roleDiv.classList.add('message-role');
    roleDiv.textContent = 'あなた';

    const contentDiv = document.createElement('div');
    contentDiv.classList.add('message-content');
    contentDiv.textContent = content;

    messageDiv.appendChild(roleDiv);
    messageDiv.appendChild(contentDiv);
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  // メッセージを送信する関数（ストリーミング対応）
  async function sendMessage() {
    const query = chatInput.value.trim();
    if (query === '') return;

    // ユーザーメッセージを表示
    addUserMessage(query);
    chatInput.value = '';

    // ボタンを無効化
    sendButton.disabled = true;
    chatInput.disabled = true;

    // ストリーミングメッセージを作成
    const { messageDiv, contentDiv } = createStreamingMessage();
    let fullContent = '';
    let sources = null;

    try {
      const response = await fetch(`/rag_search?q=${encodeURIComponent(query)}`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');

        // 最後の行は不完全な可能性があるので保持
        buffer = lines.pop() || '';

        for (const line of lines) {
          if (line.trim() === '') continue;

          if (line.startsWith('data: ')) {
            const data = line.slice(6).trim();

            try {
              const parsed = JSON.parse(data);

              // オブジェクトの場合（sources, error, done）
              if (typeof parsed === 'object' && parsed !== null) {
                if (parsed.type === 'sources') {
                  sources = parsed.data.sources;
                } else if (parsed.type === 'error') {
                  contentDiv.textContent = `エラー: ${parsed.message}`;
                } else if (parsed.type === 'done') {
                  // ストリーミング完了
                  if (sources) {
                    addSourcesToMessage(messageDiv, sources);
                  }
                }
              } else {
                // 文字列の場合（テキストチャンク）
                fullContent += parsed;
                // MarkdownをHTMLに変換して表示
                contentDiv.innerHTML = markdownToHtml(fullContent);
                chatMessages.scrollTop = chatMessages.scrollHeight;
              }
            } catch (e) {
              console.log('Parse error for data:', data, e);
            }
          }
        }
      }
    } catch (error) {
      console.error('エラーが発生しました:', error);
      contentDiv.textContent = '回答の生成中にエラーが発生しました。もう一度お試しください。';
    } finally {
      // ボタンを有効化
      sendButton.disabled = false;
      chatInput.disabled = false;
      chatInput.focus();
    }
  }

  // イベントリスナーを設定
  sendButton.addEventListener('click', sendMessage);
  chatInput.addEventListener('keydown', function (event) {
    if (event.key === 'Enter') {
      sendMessage();
    }
  });

  // ページ読み込み時にフォーカス
  window.addEventListener('load', () => {
    chatInput.focus();
  });
}
