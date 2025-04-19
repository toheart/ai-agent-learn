from langchain_text_splitters import RecursiveCharacterTextSplitter, Language
from langchain_community.vectorstores import FAISS
from langchain_openai import OpenAIEmbeddings


def read_code_from_file(file_path: str) -> str:
    with open(file_path, 'r', encoding='utf-8') as file:
        return file.read()


def split_code(code: str) -> list[str]:
    text_splitter = RecursiveCharacterTextSplitter.from_language(
        language=Language.GO, 
        chunk_size=1000, 
        chunk_overlap=0,
        )
    return text_splitter.create_documents([code])


if __name__ == "__main__":
    code = read_code_from_file("D:\\code\\rock-stack\\apps\\cgboss\\internal\\biz\\game\\dos\\cloudify.go")
    splits = split_code(code)
    embeddings = OpenAIEmbeddings(model="text-embedding-3-large")
    vectorstore = FAISS.from_documents(splits, embeddings)
    retriever = vectorstore.as_retriever()
    print(retriever.get_relevant_documents("CloudifyConfigDo"))

