from sentence_transformers import SentenceTransformer, util

def main():

    # 1. モデルのロード（ローカルにダウンロードされる）
    # AI 曰く、日本語対応かつ軽量なモデル
    model_name = "paraphrase-multilingual-MiniLM-L12-v2"
    model = SentenceTransformer(model_name)

    # 2. サンプルのドキュメント群（チャンキング後の文章を想定）
    documents = [
        "犬の健康管理について詳しく書かれたページです。",
        "ドッグランの利用規約はこちらをご確認ください。",
        "天気予報とペットの散歩に関する情報を掲載しています。",
        "カフェの営業時間とおすすめメニューについて"
    ]

    # 3. ユーザーの検索語句（例）
    query = "ドッグランのルールを知りたい"

    # 4. e5モデル向けに prefix をつける（重要！）
    query_embedding = model.encode("query: " + query, normalize_embeddings=True)
    doc_embeddings = model.encode(["passage: " + d for d in documents], normalize_embeddings=True)

    # 5. コサイン類似度の計算
    cos_scores = util.cos_sim(query_embedding, doc_embeddings)[0]

    # 6. スコアが高い順にソート
    top_k = min(3, len(documents))
    top_results = zip(range(len(documents)), cos_scores)
    top_results = sorted(top_results, key=lambda x: x[1], reverse=True)[:top_k]

    # 7. 結果表示
    print("🔍 検索クエリ:", query)
    print("\n📄 類似ドキュメント:")
    for idx, score in top_results:
        print(f"  [score={score:.4f}] {documents[idx]}")


main()