from sentence_transformers import SentenceTransformer, util
from fastapi import FastAPI

# グローバル変数でモデルを管理
model = None
model_name = "paraphrase-multilingual-MiniLM-L12-v2"

# FastAPIのインスタンスを作成
app = FastAPI()

@app.on_event("startup")
async def startup_event():
    """アプリケーション起動時にモデルを初期化"""
    global model
    print(f"Loading model: {model_name}")
    model = SentenceTransformer(model_name)
    print("Model loaded successfully!")

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
