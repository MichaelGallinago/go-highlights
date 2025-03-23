FROM python:3.10

RUN apt-get update && apt-get install -y \
    build-essential \
    && apt-get clean

WORKDIR /app

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "bot_zhenia.py"]