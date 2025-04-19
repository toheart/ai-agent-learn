from langchain_openai import ChatOpenAI
# 配置OpenAI模型
gpt_model = ChatOpenAI(
    model_name="gpt-3.5-turbo", 
    temperature=0.5,
)

response = gpt_model.invoke("你好, 我叫小唐")
print(response.content)

response = gpt_model.invoke("请问我叫什么名字？")
print(response.content)

# python .\chat-memory.py
# 你好，小唐！很高兴认识你。有什么事情我可以帮助你的吗？
# 您还没有告诉我您的名字。您可以告诉我您的名字吗？或者，如果您是想在对话中使用一个特定的名字，请直接告诉我，我很乐意称呼您。