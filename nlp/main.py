import onnxruntime as ort
from tokenizers import Tokenizer
import numpy as np
import os

# ONNXモデルとトークナイザーのパス
download_dir = os.getenv("DOWNLOAD_DIR")
onnx_model_path = os.path.join(download_dir, os.getenv("ONNX_MODEL_PATH"))
onnx_model_path = os.path.dirname(onnx_model_path) + "/model.onnx" # 最後の / 以降を削除して model.onnx とする
tokenizer_path = os.path.join(download_dir, os.getenv("TOKENIZER_PATH"))

# モデルとトークナイザーの読み込み
tokenizer = Tokenizer.from_file(tokenizer_path)
session = ort.InferenceSession(onnx_model_path, providers=["CPUExecutionProvider"])

# モデルの入力名を確認
input_names = [input.name for input in session.get_inputs()]
print(f"Required inputs: {input_names}")

# テキストをトークナイズ
text = "これは日本語の文章です。"
encoded = tokenizer.encode(text)
input_ids = np.array([encoded.ids], dtype=np.int64)
attention_mask = np.array([encoded.attention_mask], dtype=np.int64)

# token_type_idsを追加（全て0で初期化）
seq_length = len(encoded.ids)
token_type_ids = np.zeros((1, seq_length), dtype=np.int64)

# 推論
outputs = session.run(
    None,
    {
        "input_ids": input_ids,
        "attention_mask": attention_mask,
        "token_type_ids": token_type_ids  # 追加
    }
)

# 出力（[batch, seq_len, hidden_size]）
last_hidden_state = outputs[0]

# プーリング（ここではmean pooling）
attention_mask_expanded = np.expand_dims(attention_mask, axis=-1)
sum_embeddings = np.sum(last_hidden_state * attention_mask_expanded, axis=1)
sum_mask = np.clip(attention_mask_expanded.sum(axis=1), a_min=1e-9, a_max=None)
sentence_embedding = sum_embeddings / sum_mask

print("ベクトル次元:", sentence_embedding.shape)  # (1, 384)
print("先頭5次元:", sentence_embedding[0][:5]) # [-0.13136206  0.2220016  -0.03401006  0.17481032  0.1125979 ]
