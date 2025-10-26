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

  // メッセージを追加する関数
  function addMessage(role, content, sources = null) {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message');
    messageDiv.classList.add(role === 'user' ? 'user-message' : 'assistant-message');

    const roleDiv = document.createElement('div');
    roleDiv.classList.add('message-role');
    roleDiv.textContent = role === 'user' ? 'あなた' : 'アシスタント';

    const contentDiv = document.createElement('div');
    contentDiv.classList.add('message-content');
    contentDiv.textContent = content;

    messageDiv.appendChild(roleDiv);
    messageDiv.appendChild(contentDiv);

    // ソース情報を追加
    if (sources && sources.length > 0) {
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

    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  // ローディング表示を追加
  function addLoadingMessage() {
    const loadingDiv = document.createElement('div');
    loadingDiv.classList.add('loading');
    loadingDiv.id = 'loading-message';
    loadingDiv.textContent = '回答を生成中...';
    chatMessages.appendChild(loadingDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  // ローディング表示を削除
  function removeLoadingMessage() {
    const loadingDiv = document.getElementById('loading-message');
    if (loadingDiv) {
      loadingDiv.remove();
    }
  }

  // エラーメッセージを表示
  function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.classList.add('error-message');
    errorDiv.textContent = `エラー: ${message}`;
    chatMessages.appendChild(errorDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  // メッセージを送信する関数
  async function sendMessage() {
    const query = chatInput.value.trim();
    if (query === '') return;

    // ユーザーメッセージを表示
    addMessage('user', query);
    chatInput.value = '';

    // ボタンを無効化
    sendButton.disabled = true;
    chatInput.disabled = true;

    // ローディング表示
    addLoadingMessage();

    try {
      const response = await fetch(`/rag_search?q=${encodeURIComponent(query)}`);

      removeLoadingMessage();

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      // アシスタントの回答を表示
      addMessage('assistant', data.answer, data.sources);
    } catch (error) {
      removeLoadingMessage();
      console.error('エラーが発生しました:', error);
      showError('回答の生成中にエラーが発生しました。もう一度お試しください。');
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
