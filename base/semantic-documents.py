# 使用文档加载器, 嵌入和矢量存储抽象
from langchain_core.documents import Document

# 创建文档
documents = [
    Document(
        page_content="Dogs are great companions, known for their loyalty and friendliness.",
        metadata={"source": "mammal-pets-doc"},
    ),
    Document(
        page_content="Cats are independent pets that often enjoy their own space.",
        metadata={"source": "mammal-pets-doc"},
    ),
]

# 使用PyPDFLoader加载PDF文件
from langchain_community.document_loaders import PyPDFLoader

file_path = "D:\\geektime\\short_url.pdf"

loader = PyPDFLoader(file_path)
documents = loader.load()

print(f"{documents[0].page_content[:200]}\n")
print(documents[0].metadata)

# 防止文档内容过长, SplitDocuments
from langchain_text_splitters import RecursiveCharacterTextSplitter

text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200, add_start_index=True)
all_splits = text_splitter.split_documents(documents)

print(f"Splitted into {len(all_splits)} chunks")

from langchain_openai import OpenAIEmbeddings

# 这些模型指定了如何将文本转换为数字向量, 由于AI模型是基于向量空间来计算相似度的, 所以需要将文本转换为向量
embeddings = OpenAIEmbeddings(model="text-embedding-3-small")

vector_1 = embeddings.embed_query(all_splits[0].page_content)
vector_2 = embeddings.embed_query(all_splits[1].page_content)

assert len(vector_1) == len(vector_2)
print(f"Generated vectors of length {len(vector_1)}\n")
print(vector_1[:10])

from langchain_core.vectorstores import InMemoryVectorStore

# 使用InMemoryVectorStore存储向量
vector_store = InMemoryVectorStore(embeddings)

# 将向量存储到InMemoryVectorStore中
vector_store.add_documents(all_splits)


results = vector_store.similarity_search("如何生成短链数据库")

print(results[0])