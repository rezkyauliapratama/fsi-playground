from locust import HttpUser, between, task
import datetime

class LoadTestTrx(HttpUser):


    @task
    def t(self):
        dt = datetime.datetime.now()
        req = {"user_id": "03598cb5-5f84-4c30-b524-5f9bb68873e9","description" : "topup","amount" : 10,"timestamp" : dt.strftime("%Y-%m-%d %H:%M:%S")}
        self.client.post("/transactions/credit", json=req)