# ============================================================
# utils/validators.py
# ============================================================
import re

def is_valid_email(email):
    pattern = r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
    return re.match(pattern, email)

def is_valid_age(age):
    return age > 0 and age < 150


# ============================================================
# service/app.py
# ============================================================
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
from typing import Optional
from validators import is_valid_email, is_valid_age

app = FastAPI(title="User Service")

users_db = {}
next_id = 1

class UserCreate(BaseModel):
    name: str
    email: str
    age: int

class UserUpdate(BaseModel):
    name: Optional[str] = None
    email: Optional[str] = None
    age: Optional[int] = None

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/users", status_code=201)
def create_user(user: UserCreate):
    global next_id
    if not is_valid_email(user.email):
        raise HTTPException(status_code=400, detail="Invalid email")
    if not is_valid_age(user.age):
        raise HTTPException(status_code=400, detail="Invalid age")
    for u in users_db.values():
        if u["email"] == user.email:
            raise HTTPException(status_code=400, detail="Email already exists")
    user_id = next_id
    users_db[user_id] = {"id": user_id, "name": user.name, "email": user.email, "age": user.age}
    next_id += 1
    return users_db[user_id]

@app.get("/users")
def list_users(min_age: Optional[int] = Query(None)):
    result = list(users_db.values())
    if min_age:
        result = [u for u in result if u["age"] >= min_age]
    return result

@app.get("/users/{user_id}")
def get_user(user_id: int):
    return users_db[user_id]

@app.put("/users/{user_id}")
def update_user(user_id: int, user: UserUpdate):
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    stored = users_db[user_id]
    if user.name is not None:
        stored["name"] = user.name
    if user.email is not None:
        if not is_valid_email(user.email):
            raise HTTPException(status_code=400, detail="Invalid email")
        stored["email"] = user.email
    if user.age is not None:
        if not is_valid_age(user.age):
            raise HTTPException(status_code=400, detail="Invalid age")
        stored["age"] = user.age
    return stored

@app.delete("/users/{user_id}", status_code=204)
def delete_user(user_id: int):
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    del users_db[user_id]
    return None