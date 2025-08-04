from transformers import AutoTokenizer
from optimum.onnxruntime import ORTModelForFeatureExtraction

# モデルとトークナイザーの初期化
model = ORTModelForFeatureExtraction.from_pretrained(
    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2", 
    export=True
)
tokenizer = AutoTokenizer.from_pretrained(
    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"
)

# ONNXモデル、トークナイザーの保存
onnx_model_path = "./onnx_model"
model.save_pretrained(onnx_model_path)
tokenizer.save_pretrained(onnx_model_path)
print(f"ONNX model saved to: {onnx_model_path}")
