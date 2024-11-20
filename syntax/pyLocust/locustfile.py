from locust import HttpUser, task

class MyUser(HttpUser):
    @task
    def my_task(self):
        self.client.get("/")

# 命令行中运行 locust -f locustfile.py
# 打开浏览器并访问 http://localhost:8089（默认端口），您将看到 Locust 的 Web 界面。