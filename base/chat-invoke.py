import getpass
import os


os.environ["OPENAI_API_BASE"] = "http://10.86.3.248:3000/v1"
if not os.environ.get("OPENAI_API_KEY"):
  os.environ["OPENAI_API_KEY"] = getpass.getpass("Enter API key for OpenAI: ")

from langchain.chat_models import init_chat_model
from langchain_core.messages import HumanMessage, SystemMessage

model = init_chat_model("gpt-3.5-turbo", model_provider="openai")


# System Message定义AI的角色    
# Human Message定义用户的问题
# AI会根据System Message和Human Message生成响应
messages = [
    SystemMessage("Translate the following from English into Italian"),
    HumanMessage("hi!"),
]

response = model.invoke(messages)

print(response)
