from fastapi import FastAPI
from fastapi import APIRouter
import uvicorn

# 创建 APIRouter 实例
router = APIRouter()

# 定义 /health 路由处理函数
@router.get("/health")
async def health_check():
    print("health check")
    return {"status": "ok"}

# 创建 FastAPI 应用实例
app = FastAPI()

# 将 router 挂载到 FastAPI 应用上
app.include_router(router)


if __name__ == "__main__":
    uvicorn.run(app)