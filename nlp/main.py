from sentence_transformers import SentenceTransformer, util

def main():

    # モデルのロード（ローカルにダウンロードされ、キャッシュされる）
    # AI 曰く、日本語対応かつ軽量なモデル（少し古いモデルらしい）
    # 確かに GPU なしのハイエンドノート PC 程度でも動いた。
    model_name = "paraphrase-multilingual-MiniLM-L12-v2"
    model = SentenceTransformer(model_name)

    # サンプルのドキュメント群（チャンキング後の文章を想定）
    documents = [
        "犬の健康管理について詳しく書かれたページです。",
        "ドッグランの利用規約はこちらをご確認ください。",
        "天気予報とペットの散歩に関する情報を掲載しています。",
        "カフェの営業時間とおすすめメニューについて"
    ]

    # ユーザーの検索語句（例）
    query = "ドッグランのルールを知りたい"

    # ベクトル化
    # query や documents はプレフィックスといい、これをつけることで多くの現代的なモデルにおいてより良い結果が得られる。
    # 理由はモデルが文脈をより正確に理解できること、最近のモデルでは学習時にプレフィックスを使用していることが多いから。
    query_embedding = model.encode("query: " + query, normalize_embeddings=True)
    # print("🔍 クエリのベクトル化:", query_embedding)  # 結果は 384 次元のベクトル
    doc_embeddings = model.encode(["passage: " + d for d in documents], normalize_embeddings=True)
    # print("📄 ドキュメントのベクトル化:", doc_embeddings)

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


main()