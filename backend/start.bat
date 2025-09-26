@echo off
echo 启动智慧交通实时数据平台...

echo 启动后端服务...
start "后端服务" cmd /k "cd backend && go run main.go"

timeout /t 3 /nobreak > nul

echo 启动前端服务...
start "前端服务" cmd /k "cd frontend && npm run dev"

echo 服务启动完成！
echo 后端服务: http://localhost:8080
echo 前端服务: http://localhost:5173
echo 按任意键退出...
pause > nul
