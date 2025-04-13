from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_openai import ChatOpenAI

# curl https://api.openai.com/v1/chat/completions
# -H "Content-Type: application/json"
# -H "Authorization: Bearer $OPENAI_API_KEY"
# -d '{
#         "model": "gpt-4o",
#         "messages": [
#             {"role": "user", "content": "写一首关于AI的诗"}
#         ]
#     }'

prompt_template = ChatPromptTemplate.from_messages(
    [
        ("system", "Translate the following from English into Chinese:"),
        ("user", "{text}")
    ]
)

model = ChatOpenAI(model="gpt-3.5-turbo")
parser = StrOutputParser()

chain = prompt_template | model | parser
result = chain.invoke({"text":"Welcome to LLM application development!"})
print(result)