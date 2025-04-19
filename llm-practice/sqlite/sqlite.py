# 导入langchain的实用工具和相关的模块
from langchain.utilities import SQLDatabase
from langchain_openai import OpenAI
from langchain.agents import create_sql_agent
from langchain.agents.agent_toolkits import SQLDatabaseToolkit
from langchain.agents.agent_types import AgentType

# 连接到FlowerShop数据库（之前我们使用的是Chinook.db）
db = SQLDatabase.from_uri("mysql://viu_root:myviu_4359@10.63.4.122:23306/d_cgboss_info")

# 创建OpenAI的低级语言模型（LLM）实例，这里我们设置温度为0，意味着模型输出会更加确定性
llm = OpenAI(model="gpt-3.5-turbo", temperature=0, verbose=True)

# 创建SQL Agent
agent_executor = create_sql_agent(
    llm=llm,
    toolkit=SQLDatabaseToolkit(db=db, llm=llm),
    verbose=True,
    agent_type=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
)

questions = [
    "总共有多少个项目",
    ]

for question in questions:
    response = agent_executor.run(question)
    print(response)
