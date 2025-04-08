import getpass
import os


os.environ["OPENAI_API_BASE"] = "http://10.86.3.248:3000/v1"
if not os.environ.get("OPENAI_API_KEY"):
  os.environ["OPENAI_API_KEY"] = getpass.getpass("Enter API key for OpenAI: ")

from langchain.chat_models import init_chat_model

model = init_chat_model("gpt-3.5-turbo", model_provider="openai")

from langchain.prompts import ChatPromptTemplate

system_template = "Translate the following from English into {language}"
prompt_template = ChatPromptTemplate.from_messages([
    ("system", system_template),
    ("user", "{text}"),
])

prompt = prompt_template.invoke({"language": "Italian", "text": "hi!"})
print(prompt.to_messages())


response = model.invoke(prompt)
print(response.content)