from locust import HttpUser, task, between
import random
import uuid
from datetime import datetime
import time


class LoadTestTrx(HttpUser):
    wait_time = between(1, 5)  # Simulate user wait time between 1 to 5 seconds

    # List of random descriptions
    descriptions = [
        "Payment for order #1234",
        "Subscription renewal",
        "Gift card purchase",
        "Service fee payment",
        "Bonus credit",
        "Monthly subscription",
        "Charity donation",
        "Online course payment"
    ]

    @task
    def send_post_request(self):
        # Generate random data for the request
        request_data = {
            "user_id": "7cc06ed7-75ee-4ae5-b74a-8dc5cd382553",  # Static credit account number
            "description": random.choice(self.descriptions),  # Randomly select a description
            "credit_account": "ec34e120-75e0-11ef-b1a4-0242ac130002",  # Static credit account number
            "amount": round(random.uniform(10.0, 1000.0), 2),  # Random amount between 10.0 and 1000.0
            "timestamp": datetime.utcnow().isoformat()  # Current timestamp in ISO 8601 format
        }

        # Sleep for 1 second before sending the request
        time.sleep(1)

        # Send a POST request to the endpoint
        self.client.post("/transactions/debit", json=request_data)
