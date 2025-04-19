
#加载数据-分割数据-向量化-存储-检索-生成回答
from langchain.chat_models import init_chat_model
from langchain_openai import OpenAIEmbeddings
from langchain_core.vectorstores import InMemoryVectorStore

import bs4
from langchain import hub
from langchain_community.document_loaders import WebBaseLoader
from langchain_core.documents import Document
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langgraph.graph import START, StateGraph
from typing_extensions import List, TypedDict
# =====================创建模型=======================
llm = init_chat_model("gpt-3.5-turbo", model_provider="openai")
# =====================创建向量模型=======================
embeddings = OpenAIEmbeddings(model="text-embedding-3-large")
# =====================创建向量存储=======================
vector_store = InMemoryVectorStore(embeddings)
# =====================加载数据=======================
loader = WebBaseLoader(
    web_paths=("https://lilianweng.github.io/posts/2023-06-23-agent/",),
    bs_kwargs=dict(
        parse_only=bs4.SoupStrainer(
            class_=("post-content", "post-title", "post-header")
        )
    ),
)
docs = loader.load()

text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200)
all_splits = text_splitter.split_documents(docs)
# =====================存储向量=======================
_ = vector_store.add_documents(documents=all_splits)

# =====================定义prompt=======================
prompt = hub.pull("rlm/rag-prompt")
# =====================定义状态=======================
class State(TypedDict):
    question: str  # 问题
    context: List[Document]  # 上下文
    answer: str  # 回答

# =====================定义检索函数=======================
def retrieve(state: State):
    retrieved_docs = vector_store.similarity_search(state["question"])
    return {"context": retrieved_docs}

# =====================定义生成函数=======================
def generate(state: State):
    docs_content = "\n\n".join(doc.page_content for doc in state["context"])
    messages = prompt.invoke({"question": state["question"], "context": docs_content})
    response = llm.invoke(messages)
    return {"answer": response.content}

# =====================创建状态图=======================
graph_builder = StateGraph(State).add_sequence([retrieve, generate])
# =====================添加边=======================
graph_builder.add_edge(START, "retrieve")
# =====================编译图=======================
graph = graph_builder.compile()

# =====================执行=======================
response = graph.invoke({"question": "What is Task Decomposition?"})
print(response["answer"])