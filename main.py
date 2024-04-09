from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class User(BaseModel):
    username: str
    email: str

@app.get("/user")
async def read_user():
    return {"username": "test_user", "email": "test@example.com"}

@app.post("/register")
async def register_user(user: User):
    return {"username": user.username, "email": user.email}
