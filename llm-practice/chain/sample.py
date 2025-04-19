from langchain import PromptTemplate, LLMChain
from langchain_openai import ChatOpenAI

template = "{input}"
prompt_temp = PromptTemplate.from_template(template=template)
model = ChatOpenAI(model="gpt-3.5-turbo", temperature=0)
# 调用LLMChain，返回结果



