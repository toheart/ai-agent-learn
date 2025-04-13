from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
from openai import OpenAI

client = OpenAI(base_url="http://10.86.3.248:3000/v1")

# 步骤 1：生成图像描述
prompt_template = ChatPromptTemplate.from_template(
    "将以下中文描述翻译为英文的 DALL-E 提示词：{input}"
)
llm = ChatOpenAI(model="gpt-3.5-turbo")
description_chain = prompt_template | llm | StrOutputParser()

# 步骤 2：调用 DALL-E 生成图像
def generate_image(description: str) -> str:
    response = client.images.generate(
        model="dall-e-3",
        prompt=description,
        size="1024x1024",
        quality="standard",
        n=1,
    )
    return response.data[0].url

# 组合 Chain
from langchain_core.runnables import RunnableLambda

full_chain = description_chain | RunnableLambda(generate_image)

# 执行完整流程
image_url = full_chain.invoke({"input": "赛博朋克风格的未来城市"})
print(f"Generated Image URL: {image_url}")