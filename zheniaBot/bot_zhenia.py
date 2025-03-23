import asyncio
import logging
import os
import grpc
from aiogram import Bot, Dispatcher, types

import api_pb2
import api_pb2_grpc

# Указываем токен Telegram-бота
TOKEN = os.getenv("TELEGRAM_ZHENIA_BOT_TOKEN")
GRPC_SERVER = os.getenv("ZHENIA_BOT_GRPC_SERVER")

bot = Bot(token=TOKEN)
dp = Dispatcher()


# Функция вызова gRPC-методов с обработкой ошибок
async def call_grpc_method(method_name: str, request):
    try:
        async with grpc.aio.insecure_channel(GRPC_SERVER) as channel:
            stub = api_pb2_grpc.RequesterServiceStub(channel)
            method = getattr(stub, method_name, None)
            if method:
                response = await method(request)
                return response
            else:
                logging.error(f"Метод {method_name} не найден в gRPC-сервисе.")
                return None
    except grpc.aio.AioRpcError as e:
        logging.error(f"Ошибка gRPC: {e.code()} - {e.details()}")
        return None


# Проверка месяца и конвертация в число
def parse_month(month: str):
    months = {
        "январь": 1, "февраль": 2, "март": 3, "апрель": 4, "май": 5, "июнь": 6,
        "июль": 7, "август": 8, "сентябрь": 9, "октябрь": 10, "ноябрь": 11, "декабрь": 12
    }
    return months.get(month.lower(), None)


# Обработчик команд
@dp.message()
async def handle_commands(message: types.Message):
    args = message.text.split(maxsplit=1)
    command = args[0]

    if command == "/highlight_month":
        if len(args) < 2:
            await message.answer("❌ Введите месяц! Пример: /highlight_month январь")
            return

        month_num = parse_month(args[1])
        if not month_num:
            await message.answer("❌ Неверное название месяца! Введите, например: январь, февраль, март...")
            return

        request = api_pb2.MonthHighlightRequest(month=month_num)
        response = await call_grpc_method("GetMemesByMonth", request)
        if response:
            await message.answer(response.text or "⚠️ Нет данных.")
        else:
            await message.answer("❌ Ошибка: сервис недоступен или произошла ошибка при обработке запроса.")

    elif command == "/highlight_word":
        if len(args) < 2:
            await message.answer("❌ Введите слово или фразу! Пример: /highlight_word важная новость")
            return

        search_phrase = args[1]
        logging.info(f"🔍 Поиск хайлайтов по фразе: {search_phrase}")

        request = api_pb2.SearchHighlightRequest(query=search_phrase)
        response = await call_grpc_method("SearchMemesBySubstring", request)

        if response:
            await message.answer(response.text or "⚠️ Ничего не найдено.")
        else:
            await message.answer("❌ Ошибка: сервис недоступен или произошла ошибка при обработке запроса.")

    elif command == "/highlight_random":
        response = await call_grpc_method("GetRandomMeme", api_pb2.EmptyHighlightRequest())
        if response:
            await message.answer(response.text or "⚠️ Нет данных.")
        else:
            await message.answer("❌ Ошибка: сервис недоступен или произошла ошибка при обработке запроса.")

    elif command == "/highlight_top5":
        request = api_pb2.TopLongMemesHighlightRequest(limit=5)
        response = await call_grpc_method("GetTopLongMemes", request)
        if response:
            await message.answer(response.text or "⚠️ Нет данных.")
        else:
            await message.answer("❌ Ошибка: сервис недоступен или произошла ошибка при обработке запроса.")

    else:
        await message.answer("⚠️ Неизвестная команда. Используйте:\n"
                             "/highlight_month <месяц>\n"
                             "/highlight_word <фраза>\n"
                             "/highlight_random\n"
                             "/highlight_top5")


# Запуск бота
async def main():
    logging.basicConfig(level=logging.INFO)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())