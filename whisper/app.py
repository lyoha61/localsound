import os
import tempfile

import whisperx
from fastapi import FastAPI, UploadFile, File

app = FastAPI()

device = "cuda"
compute_type = "float16"

model = whisperx.load_model(
    "large-v3",
    device=device,
    compute_type=compute_type
)


@app.post("/transcribe")
async def transcribe(file: UploadFile = File(...)):
    with tempfile.NamedTemporaryFile(
        suffix=os.path.splitext(file.filename)[1],
        delete=False
    ) as temp:
        temp.write(await file.read())
        audio_path = temp.name

    try:
        audio = whisperx.load_audio(audio_path)

        result = model.transcribe(audio)

        model_a, metadata = whisperx.load_align_model(
            language_code=result["language"],
            device=device
        )

        result = whisperx.align(
            result["segments"],
            model_a,
            metadata,
            audio,
            device
        )

        return result

    finally:
        os.remove(audio_path)