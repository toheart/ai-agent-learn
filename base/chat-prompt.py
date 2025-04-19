from langchain.chat_models import init_chat_model
from langchain.prompts import ChatPromptTemplate
model = init_chat_model("gpt-3.5-turbo", model_provider="openai")
system_template = "Translate the following from English into {language}"
prompt_template = ChatPromptTemplate.from_messages([
    ("system", system_template),
    ("user", "{text}"),
])
prompt = prompt_template.invoke({"language": "Italian", "text": "hi!"})
response = model.invoke(prompt)
print(response.content)

prompt = prompt_template.invoke({"language": "Chinese", "text": "hi!"})
response = model.invoke(prompt)
print(response.content)