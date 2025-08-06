from sentence_transformers import SentenceTransformer, util
from fastapi import FastAPI

# グローバル変数でモデルを管理
model = None
model_name = "paraphrase-multilingual-MiniLM-L12-v2"

# FastAPIのインスタンスを作成
app = FastAPI()

# =====================================================
# ルーティング
# =====================================================
@app.on_event("startup")
async def startup_event():
    """アプリケーション起動時にモデルを初期化"""
    global model
    print(f"モデルを読み込み中: {model_name}")
    model = SentenceTransformer(model_name)
    print("モデルの読み込みが完了しました！")

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/convert/{input_text}")
def convert_text(input_text: str):
    vector = convert_to_vector(input_text)
    return {
        "input_text": input_text,
        "vector": vector.tolist(),  # numpy配列をリストに変換
        "dimensions": len(vector)
    }

# =====================================================
# ベクトル化関数
# =====================================================
def convert_to_vector(input_text: str, is_query: bool = True):
    """
    テキストをベクトルに変換する関数
    
    Args:
        input_text (str): ベクトル化するテキスト
        is_query (bool): クエリかどうか（True: クエリ, False: ドキュメント）

    Returns:
        numpy.ndarray: 正規化されたベクトル（384次元）
    """

    if is_query:
        prefix = "query: "
    else:
        prefix = "passage: "

    vector = model.encode(prefix + input_text, normalize_embeddings=True)

    return vector
