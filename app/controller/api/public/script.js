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
            console.log('Received data:', data); // デバッグログ

            try {
              const parsed = JSON.parse(data);
              console.log('Parsed:', parsed, 'Type:', typeof parsed); // デバッグログ

              // オブジェクトの場合（sources, error, done）
              if (typeof parsed === 'object' && parsed !== null) {
                if (parsed.type === 'sources') {
                  sources = parsed.data.sources;
                  console.log('Sources received:', sources);
                } else if (parsed.type === 'error') {
                  contentDiv.textContent = `エラー: ${parsed.message}`;
                } else if (parsed.type === 'done') {
                  // ストリーミング完了
                  console.log('Stream done, adding sources');
                  if (sources) {
                    addSourcesToMessage(messageDiv, sources);
                  }
                }
              } else {
                // 文字列の場合（テキストチャンク）
                console.log('Adding text chunk:', parsed);
                fullContent += parsed;
                contentDiv.textContent = fullContent;
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
