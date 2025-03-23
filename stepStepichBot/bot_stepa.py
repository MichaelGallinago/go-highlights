import asyncio
import logging
import os
import pika
import json
from aiogram import Bot, Dispatcher, F
from aiogram.types import Message

TOKEN = os.getenv("TELEGRAM_STEPA_BOT_TOKEN")
RABBITMQ_HOST = os.getenv("RABBITMQ_HOST", "rabbit")
QUEUE_NAME = os.getenv("RABBITMQ_QUEUE", "messages")

bot = Bot(token=TOKEN)
dp = Dispatcher()


def send_to_rabbitmq(message_text: str, timestamp: str):
    try:
        connection = pika.BlockingConnection(pika.ConnectionParameters(host=RABBITMQ_HOST))
        channel = connection.channel()
        channel.queue_declare(queue=QUEUE_NAME, durable=True)

        message = {
            "timestamp": timestamp,
            "text": message_text
        }

        message_json = json.dumps(message)

        channel.basic_publish(
            exchange='',
            routing_key=QUEUE_NAME,
            body=message_json
        )

        connection.close()
        logging.info(f"✅ Сообщение отправлено в RabbitMQ: {message}")
    except Exception as e:
        logging.error(f"❌ Ошибка отправки в RabbitMQ: {e}")


@dp.message(F.text.contains("#хайлайты"))
async def handle_highlight_messages(message: Message):
    timestamp = message.date.isoformat()
    send_to_rabbitmq(message.text, timestamp)
    await message.answer("📨 Сообщение с #хайлайты отправлено в очередь!")


async def main():
    logging.basicConfig(level=logging.INFO)
    await bot.delete_webhook(drop_pending_updates=True)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
