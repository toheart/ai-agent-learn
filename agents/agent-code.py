from langchain.embeddings.openai import OpenAIEmbeddings
from langchain_core.vectorstores import InMemoryVectorStore
from langchain.document_loaders import TextLoader
from langchain.text_splitter import CharacterTextSplitter, PythonCodeTextSplitter
from langchain.chat_models import init_chat_model
from langchain.chains import ConversationalRetrievalChain



sources = []
loader = TextLoader('./agent-code.py', encoding='utf8')
sources.extend(loader.load_and_split())

splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=0)
docs = splitter.split_documents(sources)

embeddings = OpenAIEmbeddings(model="text-embedding-3-large")

db = InMemoryVectorStore(embeddings)

db.add_documents(docs)

retriever = db.as_retriever()
retriever.search_kwargs['distance_metric'] = 'cos'
retriever.search_kwargs['fetch_k'] = 100
retriever.search_kwargs['maximal_marginal_relevance'] = True
retriever.search_kwargs['k'] = 10

model = init_chat_model(model='gpt-3.5-turbo', model_provider='openai', temperature=0) # 'gpt-3.5-turbo',

qa = ConversationalRetrievalChain.from_llm(model,retriever=retriever)


questions = [
    '解释这段代码',
    'generate方法是做什么的?'
]
chat_history = []

for question in questions:  
    result = qa({"question": question, "chat_history": chat_history})
    chat_history.append((question, result['answer']))
    print(f"-> **Question**: {question} \n")
    print(f"**Answer**: {result['answer']} \n")