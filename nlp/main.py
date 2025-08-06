from sentence_transformers import SentenceTransformer, util
from fastapi import FastAPI

# グローバルでモデルを初期化（関数呼び出しのたびに読み込まないため）
model_name = "paraphrase-multilingual-MiniLM-L12-v2"
model = SentenceTransformer(model_name)

# FastAPIのインスタンスを作成
app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/convert/{input_text}")
def convert_text(input_text: str):
    vector = convert_to_vector(input_text)
    return {"vector": vector}

def convert_to_vector(input_text: str, is_query: bool = True):
    """
    テキストをベクトルに変換する関数
    
    Args:
        input_text (str): ベクトル化するテキスト
        is_query (bool): クエリかどうか（True: クエリ, False: ドキュメント）

    Returns:
        numpy.ndarray: 正規化されたベクトル（384次元）
    """

    print(input_text)

    if is_query:
        prefix = "query: "
    else:
        prefix = "passage: "

    vector = model.encode(prefix + input_text, normalize_embeddings=True)

    print(vector)

    return vector

def main():
    # テスト用の入力
    input_text = "こんにちは"
    vector = convert_to_vector(input_text)
    print(f"Input: {input_text}\nVector: {vector}")

if __name__ == "__main__":
    main()
    # FastAPIのサーバーを起動するためのコマンドは、uvicornを使用して実行すること
    # 例: uvicorn main:app --reload --host