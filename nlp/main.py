from transformers import AutoTokenizer
from optimum.onnxruntime import ORTModelForFeatureExtraction
import torch
import os

onnx_model_path = "./onnx_model"
model_id = "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"
commit_hash='86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d'

if os.path.exists(onnx_model_path + "/model.onnx"):
    # 既存のONNXモデルを使用
    print("Loading existing ONNX model...")
    model = ORTModelForFeatureExtraction.from_pretrained(onnx_model_path)
    tokenizer = AutoTokenizer.from_pretrained(onnx_model_path)
else:
    # Hugging Faceからモデルをダウンロードし、同時にONNX形式に変換
    print("Converting model to ONNX...")
    model = ORTModelForFeatureExtraction.from_pretrained(
        model_id,
        revision=commit_hash,
        export=True
    )
    tokenizer = AutoTokenizer.from_pretrained(model_id)
    
    # メモリ上のモデルをディスクに保存
    model.save_pretrained(onnx_model_path)
    tokenizer.save_pretrained(onnx_model_path)

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
