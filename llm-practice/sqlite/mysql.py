from langchain_community.agent_toolkits import SQLDatabaseToolkit
from langchain.chat_models import init_chat_model
from langchain_community.utilities import SQLDatabase

db = SQLDatabase.from_uri("mysql://viu_root:myviu_4359@10.63.4.122:23306/d_cgboss_info")
llm = init_chat_model("gpt-3.5-turbo", model_provider="openai")
toolkit = SQLDatabaseToolkit(db=db, llm=llm)

tools = toolkit.get_tools()

from langchain import hub
prompt_template = hub.pull("langchain-ai/sql-agent-system-prompt")
system_message = prompt_template.format(dialect="MySQL", top_k=5)
from langchain_core.messages import HumanMessage
from langgraph.prebuilt import create_react_agent

agent_executor = create_react_agent(llm, tools, prompt=system_message)

question = "task_pitem表中总共有多少个任务"

for step in agent_executor.stream(
    {"messages": [{"role": "user", "content": question}]},
    stream_mode="values",
):
    step["messages"][-1].pretty_print()