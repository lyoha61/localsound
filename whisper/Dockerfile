FROM nvidia/cuda:12.6.3-runtime-ubuntu24.04

RUN apt-get update && \
    apt-get install -y python3 python3-pip ffmpeg && \
    rm -rf /var/lib/apt/lists/*

RUN pip3 install whisperx fastapi uvicorn python-multipart

COPY app.py /app/app.py

WORKDIR /app

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]