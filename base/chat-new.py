# from langchain_openai import ChatOpenAI
# # 配置OpenAI模型
# gpt_model = ChatOpenAI(
#     model_name="gpt-3.5-turbo", 
#     temperature=0.5,
# )

# response = gpt_model.invoke("你好")
# print(response)

# from langchain.chat_models import init_chat_model

# # 自动推断厂商（需安装对应SDK）
# gpt_model = init_chat_model("gpt-3.5-turbo", temperature=0)  # 推断为OpenAI
# response = gpt_model.invoke("你好")
# print(response)
import os
apikey = os.getenv("OPENAI_API_KEY")
base_url = os.getenv("OPENAI_API_BASE")

from openai import OpenAI
# 原生OpenAI调用
client = OpenAI(api_key=apikey, base_url=base_url)
response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "你好"}]
)
print(response)
# 转换为LangChain消息格式
from langchain_core.messages import AIMessage
message = AIMessage(content=response.choices[0].message.content)
print(message)