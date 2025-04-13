
#加载数据-分割数据-向量化-存储-检索-生成回答
from langchain.chat_models import init_chat_model
from langchain_openai import OpenAIEmbeddings
from langchain_core.vectorstores import InMemoryVectorStore

import bs4
from langchain import hub
from langchain_community.document_loaders import ConfluenceLoader
from langchain_core.documents import Document
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langgraph.graph import START, StateGraph
from typing_extensions import List, TypedDict
import logging

# 启用日志调试
logging.basicConfig(level=logging.DEBUG)
cql_query = 'creator="levi.tang"'  
#创建模型
llm = init_chat_model("gpt-3.5-turbo", model_provider="openai")
#创建向量模型, 用来将文本转换为向量
embeddings = OpenAIEmbeddings(model="text-embedding-3-large")
#创建向量存储
vector_store = InMemoryVectorStore(embeddings)
# Load and chunk contents of the blog
loader = ConfluenceLoader(
    space_key="YYX", include_attachments=False, limit=50, max_pages=100, cql=cql_query
)
docs = loader.load()

text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200)
all_splits = text_splitter.split_documents(docs)

# Index chunks
_ = vector_store.add_documents(documents=all_splits)


# Define prompt for question-answering
prompt = hub.pull("rlm/rag-prompt")

# Define state for application
class State(TypedDict):
    question: str  # 问题
    context: List[Document]  # 上下文
    answer: str  # 回答


# Define application steps
def retrieve(state: State):
    retrieved_docs = vector_store.similarity_search(state["question"])
    return {"context": retrieved_docs}

def generate(state: State):
    docs_content = "\n\n".join(doc.page_content for doc in state["context"])
    messages = prompt.invoke({"question": state["question"], "context": docs_content})
    response = llm.invoke(messages)
    return {"answer": response.content}

# Compile application and test
# 创建状态图
graph_builder = StateGraph(State).add_sequence([retrieve, generate])
# 添加边
graph_builder.add_edge(START, "retrieve")
# 编译图
graph = graph_builder.compile()

# 执行图
response = graph.invoke({"question": "什么是云游戏？"})
print(response["answer"])