from transformers import AutoTokenizer, AutoModel
from optimum.onnxruntime import ORTModelForFeatureExtraction
import torch
import os

# モデルとトークナイザーの初期化
model = ORTModelForFeatureExtraction.from_pretrained(
    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2", 
    export=True
)
tokenizer = AutoTokenizer.from_pretrained(
    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"
)

# ONNXモデルを保存
onnx_model_path = "./onnx_model"
model.save_pretrained(onnx_model_path)
tokenizer.save_pretrained(onnx_model_path)
print(f"ONNX model saved to: {onnx_model_path}")

# 保存されたファイルを確認
if os.path.exists(onnx_model_path):
    files = os.listdir(onnx_model_path)
    print(f"Saved files: {files}")

# .onnx 生成目的であれば削除できそう
def vectorize_text(text):
    """テキストをベクトル化する関数"""
    # トークン化
    inputs = tokenizer(text, return_tensors="pt", padding=True, truncation=True, max_length=512)
    
    # 推論実行
    with torch.no_grad():
        outputs = model(**inputs)
    
    # 平均プーリング
    embeddings = outputs.last_hidden_state
    attention_mask = inputs['attention_mask']
    
    # マスクされた平均を計算
    masked_embeddings = embeddings * attention_mask.unsqueeze(-1)
    summed = torch.sum(masked_embeddings, 1)
    counts = torch.clamp(attention_mask.sum(1), min=1e-9)
    mean_pooled = summed / counts.unsqueeze(-1)
    
    return mean_pooled.numpy()

# .onnx 生成目的であれば削除できそう
def main():
    print("Model loaded successfully!")
    
    # テスト
    test_text = "こんにちは、世界"
    vector = vectorize_text(test_text)
    print(f"Vector shape: {vector.shape}")
    print(f"First 5 values: {vector[0][:5]}")

# .onnx 生成目的であれば削除できそう
if __name__ == "__main__":
    main()
