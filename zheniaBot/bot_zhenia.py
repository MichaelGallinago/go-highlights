import asyncio
import logging
import os
import grpc
from aiogram import Bot, Dispatcher, types

import api_pb2
import api_pb2_grpc

TOKEN = os.getenv("TELEGRAM_ZHENIA_BOT_TOKEN")
GRPC_SERVER = os.getenv("ZHENIA_BOT_GRPC_SERVER")

bot = Bot(token=TOKEN)
dp = Dispatcher()

stub: api_pb2_grpc.RequesterServiceStub = None

def parse_month(month: str):
    months = {
        "январь": 1, "февраль": 2, "март": 3, "апрель": 4, "май": 5, "июнь": 6,
        "июль": 7, "август": 8, "сентябрь": 9, "октябрь": 10, "ноябрь": 11, "декабрь": 12
    }
    return months.get(month.lower(), None)


def format_with_title(title: str, content: str) -> str:
    content = content.strip() if content else ""
    return f"{title}\n\n{content or '⚠️ Нет данных.'}"


def split_text(text: str, max_length: int = 4096) -> list[str]:
    lines = text.split('\n')
    chunks = []
    current = ""
    for line in lines:
        if len(current) + len(line) + 1 <= max_length:
            current += line + '\n'
        else:
            chunks.append(current.strip())
            current = line + '\n'
    if current:
        chunks.append(current.strip())
    return chunks


async def safe_send_message(message: types.Message, text: str, parse_mode="Markdown"):
    for chunk in split_text(text, max_length=4096):
        await message.answer(chunk, parse_mode=parse_mode)


async def handle_highlight_month(message: types.Message, arg: str):
    month_num = parse_month(arg)
    if not month_num:
        await message.answer("❌ Неверное название месяца! Введите, например: январь, февраль, март...")
        return

    request = api_pb2.MonthHighlightRequest(month=month_num)
    response = await stub.GetMemesByMonth(request)

    title = f"📅 *Мемы за месяц:* _{arg.capitalize()}_"
    await safe_send_message(message, format_with_title(title, response.text))


async def handle_highlight_word(message: types.Message, query: str):
    request = api_pb2.SearchHighlightRequest(query=query)
    response = await stub.SearchMemesBySubstring(request)

    title = f"🔍 *Поиск по фразе:* _{query}_"
    await safe_send_message(message, format_with_title(title, response.text))


async def handle_highlight_random(message: types.Message):
    response = await stub.GetRandomMeme(api_pb2.EmptyHighlightRequest())
    title = "🎲 *Случайный хайлайт*"
    await safe_send_message(message, format_with_title(title, response.text))


async def handle_highlight_top5(message: types.Message):
    request = api_pb2.TopLongMemesHighlightRequest(limit=5)
    response = await stub.GetTopLongMemes(request)
    title = "🏆 *ТОП-5 самых длинных мемов:*"
    await safe_send_message(message, format_with_title(title, response.text))


@dp.message()
async def handle_commands(message: types.Message):
    args = message.text.split(maxsplit=1)
    command = args[0]

    if command == "/highlight_month":
        if len(args) < 2:
            await message.answer("❌ Введите месяц! Пример: /highlight_month январь")
            return
        await handle_highlight_month(message, args[1])

    elif command == "/highlight_word":
        if len(args) < 2:
            await message.answer("❌ Введите слово или фразу! Пример: /highlight_word важная новость")
            return
        await handle_highlight_word(message, args[1])

    elif command == "/highlight_random":
        await handle_highlight_random(message)

    elif command == "/highlight_top5":
        await handle_highlight_top5(message)

    else:
        await message.answer("⚠️ Неизвестная команда. Используйте:\n"
                             "/highlight_month <месяц>\n"
                             "/highlight_word <фраза>\n"
                             "/highlight_random\n"
                             "/highlight_top5")


async def main():
    global stub
    logging.basicConfig(level=logging.INFO)
    channel = grpc.aio.insecure_channel(GRPC_SERVER)
    stub = api_pb2_grpc.RequesterServiceStub(channel)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
