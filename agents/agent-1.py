from langchain_openai import ChatOpenAI
from langchain_community.tools.tavily_search import TavilySearchResults
from langchain_core.messages import HumanMessage
from langgraph.checkpoint.memory import MemorySaver
from langgraph.prebuilt import create_react_agent


# 创建记忆
memory = MemorySaver()

# 创建模型
model = ChatOpenAI(model_name="gpt-3.5-turbo")

# 创建搜索工具
search = TavilySearchResults(max_results=2)

# 创建工具列表
tools = [search]

# 使用代理
config = {"configurable": {"thread_id": "abc123"}}
# 创建代理
agent_executor = create_react_agent(model, tools)



for step in agent_executor.stream(
    {"messages": [HumanMessage(content="hi im bob! and i live in sf")]},
    config,
    stream_mode="values",
):
    step["messages"][-1].pretty_print()
# 体验不好, 需要优化
# response = agent_executor.invoke(
#     {"messages": [HumanMessage(content="whats the weather in sf?")]}, config
# )
# response["messages"]

# 使用信息流, 信息发生改变时, 信息流会更新
# for step in agent_executor.stream(
#     {"messages": [HumanMessage(content="whats the weather where I live?")]},
#     config,
#     stream_mode="values",
# ):
#     step["messages"][-1].pretty_print()

for step, metadata in agent_executor.stream(
    {"messages": [HumanMessage(content="whats the weather in sf?")]},
    stream_mode="messages",
):
    if metadata["langgraph_node"] == "agent" and (text := step.text()):
        print(text, end="|")