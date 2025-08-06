from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer, util
import neologdn

# グローバル変数でモデルを管理
model = None
model_name = "paraphrase-multilingual-MiniLM-L12-v2"

# FastAPIのインスタンスを作成
app = FastAPI()

# リクエストボディ用のPydanticモデル
class ConvertRequest(BaseModel):
    text: str
    is_query: bool = True

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

@app.post("/convert")
def convert(request: ConvertRequest):
    """ テキストをベクトルに変換するエンドポイント """
    normalized_text = normalize_text(request.text)  # テキストを正規化
    vector = convert_to_vector(normalized_text, request.is_query)  # ベクトルに変換
    return {
        "input_text": request.text,
        "normalized_text": normalized_text,
        "is_query": request.is_query,
        "model_name": model_name,
        "dimensions": len(vector),
        "vector": vector.tolist(),  # numpy配列をリストに変換
    }

# =====================================================
# 処理関数
# =====================================================
def normalize_text(text: str) -> str:
    """
    テキストを正規化する関数
    
    Args:
        text (str): 正規化するテキスト

    Returns:
        str: 正規化されたテキスト
    """
    return neologdn.normalize(text)

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
