from langchain_core.messages import HumanMessage, SystemMessage
from langchain.chat_models import init_chat_model

model = init_chat_model(
    model="gpt-3.5-turbo",
)

response = model.invoke([HumanMessage(content="Hi! I'm Bob")])

print(response.content)

response = model.invoke([HumanMessage(content="What's my name?")])

print(response.content)


# 缺少记忆功能, 每个对话都是独立的

# 解决将整个对话记录作为上下文传递给模型
from langchain_core.messages import AIMessage
response = model.invoke([
    HumanMessage(content="Hi! I'm Bob"),
    AIMessage(content="Hi! I'm Alice"),
    HumanMessage(content="What's my name?"),
])

print(response.content)


# 问题: 如果对话记录很长, 如何处理?
# 使用Langgraph来管理对话记录

from langgraph.checkpoint.memory import MemorySaver
from langgraph.graph import START, MessagesState, StateGraph


# 定义状态
workflow = StateGraph(state_schema=MessagesState)


# 定义模型调用函数
def call_model(state: MessagesState):
    response = model.invoke(state["messages"])
    return {"messages": response}


# 定义图的节点
workflow.add_edge(START, "model")
workflow.add_node("model", call_model)

# 添加记忆
memory = MemorySaver()
app = workflow.compile(checkpointer=memory)

config = {"configurable": {"thread_id": "abc123"}}
query = "Hi! I'm Bob."

input_messages = [HumanMessage(query)]
output = app.invoke({"messages": input_messages}, config)
output["messages"][-1].pretty_print()  # output contains all messages in state

query = "What's my name?"
input_messages = [HumanMessage(query)]
output = app.invoke({"messages": input_messages}, config)
output["messages"][-1].pretty_print()

# 使用不同的线程ID, 进行验证
config = {"configurable": {"thread_id": "abc234"}}
input_messages = [HumanMessage(query)]
output = app.invoke({"messages": input_messages}, config)
output["messages"][-1].pretty_print()

