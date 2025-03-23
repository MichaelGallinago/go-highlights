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

stub = None


async def init_grpc_stub():
    global stub
    channel = grpc.aio.insecure_channel(GRPC_SERVER)
    stub = api_pb2_grpc.RequesterServiceStub(channel)


def format_with_title(title: str, content: str) -> str:
    return f"{title}\n\n{content or '⚠️ Нет данных.'}"


def parse_month(month: str):
    months = {
        "январь": 1, "февраль": 2, "март": 3, "апрель": 4, "май": 5, "июнь": 6,
        "июль": 7, "август": 8, "сентябрь": 9, "октябрь": 10, "ноябрь": 11, "декабрь": 12
    }
    return months.get(month.lower())


async def handle_highlight_month(message: types.Message, month_name: str):
    month_num = parse_month(month_name)
    if not month_num:
        await message.answer("❌ Неверное название месяца! Введите, например: январь, февраль, март...")
        return

    try:
        response = await stub.GetMemesByMonth(api_pb2.MonthHighlightRequest(month=month_num))
        title = f"📅 *Хайлайты за месяц:* _{month_name.capitalize()}_"
        await message.answer(format_with_title(title, response.text), parse_mode="Markdown")
    except grpc.aio.AioRpcError as e:
        logging.error(f"gRPC error: {e.code()} - {e.details()}")
        await message.answer("❌ Ошибка при вызове gRPC.")


async def handle_highlight_word(message: types.Message, search_phrase: str):
    try:
        response = await stub.SearchMemesBySubstring(api_pb2.SearchHighlightRequest(query=search_phrase))
        title = f"🔍 *Поиск по фразе:* _{search_phrase}_"
        await message.answer(format_with_title(title, response.text), parse_mode="Markdown")
    except grpc.aio.AioRpcError as e:
        logging.error(f"gRPC error: {e.code()} - {e.details()}")
        await message.answer("❌ Ошибка при вызове gRPC.")


async def handle_highlight_random(message: types.Message):
    try:
        response = await stub.GetRandomMeme(api_pb2.EmptyHighlightRequest())
        title = "🎲 *Случайный хайлайт*"
        await message.answer(format_with_title(title, response.text), parse_mode="Markdown")
    except grpc.aio.AioRpcError as e:
        logging.error(f"gRPC error: {e.code()} - {e.details()}")
        await message.answer("❌ Ошибка при вызове gRPC.")


async def handle_highlight_top5(message: types.Message):
    try:
        response = await stub.GetTopLongMemes(api_pb2.TopLongMemesHighlightRequest(limit=5))
        title = "🏆 *ТОП-5 самых длинных хайлайтов*"
        await message.answer(format_with_title(title, response.text), parse_mode="Markdown")
    except grpc.aio.AioRpcError as e:
        logging.error(f"gRPC error: {e.code()} - {e.details()}")
        await message.answer("❌ Ошибка при вызове gRPC.")


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
            await message.answer("❌ Введите слово или фразу! Пример: /highlight_word 1с")
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
    logging.basicConfig(level=logging.INFO)
    await init_grpc_stub()
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
