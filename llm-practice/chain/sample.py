from langchain import PromptTemplate

template = "{flower}的花语是?"
prompt_temp = PromptTemplate.from_template(template=template)
prompt = prompt_temp.format(flower="玫瑰")
print(prompt)

from  langchain_community.llms import OpenAI    

model = OpenAI(model="gpt-3.5-turbo", temperature=0)

result = model.invoke(prompt)
print(result)




