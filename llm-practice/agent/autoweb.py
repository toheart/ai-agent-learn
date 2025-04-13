from langchain.agents.agent_toolkits import PlayWrightBrowserToolkit
from langchain_community.tools.playwright.utils import create_async_playwright_browser

async_browser = create_async_playwright_browser()
toolkit = PlayWrightBrowserToolkit.from_browser(async_browser=async_browser)
tools = toolkit.get_tools()
print(tools)

from langchain.agents import initialize_agent, AgentType
from langchain_community.llms import OpenAI

# LLM不稳定，对于这个任务，可能要多跑几次才能得到正确结果
llm = OpenAI(model="gpt-3.5-turbo", temperature=0.5)  

agent_chain = initialize_agent(
    tools,
    llm,
    agent=AgentType.STRUCTURED_CHAT_ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True,
)

async def main():
    response = await agent_chain.arun("What are  python.langchain.com?")
    print(response)

import asyncio
loop = asyncio.get_event_loop()
loop.run_until_complete(main())