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

# 使用stream方法，可以实时输出AI的响应, 提高用户体验
for token in model.stream(messages):
    print(token.content, end="|")