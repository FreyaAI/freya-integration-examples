import base64
import datetime

import httpx
import uvicorn
from Cryptodome.Hash import SHA256
from Cryptodome.PublicKey import RSA
from Cryptodome.Signature import PKCS1_v1_5
from fastapi import FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


def read_pem_file(file_path: str) -> str:
    with open(file_path, "r") as file:
        return file.read()


def sign(message: str, private_key_str: str) -> str:
    priv_key = RSA.importKey(private_key_str, passphrase=None)
    h = SHA256.new(message.encode("utf-8"))
    signature = PKCS1_v1_5.new(priv_key).sign(h)
    result = base64.b64encode(signature).decode()
    return result


async def post_authentication_details(
    company_code: str, user_id: str, signature: str, timestamp: str
) -> str:
    url = "https://api.freyafashion.ai/api/v1/authenticate"
    body = {
        "company_code": company_code,
        "user_id": user_id,
        "signature": signature,
        "timestamp": timestamp,
    }
    timeout = 30.0  # Set your desired timeout in seconds here

    async with httpx.AsyncClient() as client:
        # Pass the timeout parameter to the post method
        response = await client.post(url, json=body, timeout=timeout)
        response.raise_for_status()  # This will raise an exception for 4XX/5XX responses
        data = response.json()  # Parse the JSON response
        return data["token"]


@app.post("/demo/v1/authenticate")
async def authenticate(request: Request):
    body = await request.json()
    user_id = body.get("user_id")
    company_code = body.get("company_code")

    if not user_id or not company_code:
        raise HTTPException(status_code=400, detail="Missing user_id or company_code")

    private_key_str = read_pem_file(f"{company_code}.cer")
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    signed_message = sign(timestamp, private_key_str)

    # Send authentication details to external service and get token
    token = await post_authentication_details(
        company_code, user_id, signed_message, timestamp
    )

    return {"token": token}


if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8000)
