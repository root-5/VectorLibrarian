from sentence_transformers import SentenceTransformer, util

# グローバルでモデルを初期化（関数呼び出しのたびに読み込まないため）
model_name = "paraphrase-multilingual-MiniLM-L12-v2"
model = SentenceTransformer(model_name)

def convertToVector(inputText, isQuery=True):
    """
    テキストをベクトルに変換する関数
    
    Args:
        inputText (str): ベクトル化するテキスト
        isQuery (bool): クエリかどうか（True: クエリ, False: ドキュメント）
    
    Returns:
        numpy.ndarray: 正規化されたベクトル（384次元）
    """
    if isQuery:
        prefix = "query: "
    else:
        prefix = "passage: "
    
    vector = model.encode(prefix + inputText, normalize_embeddings=True)
    return vector

def main():
    # サンプルのドキュメント群（チャンキング後の文章を想定）
    documents = [
        "犬の健康管理について詳しく書かれたページです。",
        "ドッグランの利用規約はこちらをご確認ください。",
        "天気予報とペットの散歩に関する情報を掲載しています。",
        "カフェの営業時間とおすすめメニューについて"
    ]

    # ユーザーの検索語句（例）
    query = "ドッグランのルールを知りたい"

    # ベクトル化（関数を使用）
    query_embedding = convertToVector(query, isQuery=True)
    doc_embeddings = [convertToVector(doc, isQuery=False) for doc in documents]

    # コサイン類似度の計算
    cos_scores = util.cos_sim(query_embedding, doc_embeddings)[0]

    # スコアが高い順にソート
    top_k = min(3, len(documents))
    top_results = zip(range(len(documents)), cos_scores)
    top_results = sorted(top_results, key=lambda x: x[1], reverse=True)[:top_k]

    # 結果表示
    print("🔍 検索クエリ:", query)
    print("\n📄 類似ドキュメント:")
    for idx, score in top_results:
        print(f"  [score={score:.4f}] {documents[idx]}")

# 直接実行された場合のみmainを呼び出す
if __name__ == "__main__":
    main()