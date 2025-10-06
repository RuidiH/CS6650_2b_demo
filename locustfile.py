# locustfile.py
import random
from locust import HttpUser, FastHttpUser, task, between

PRODUCT_ID_MIN = 1
PRODUCT_ID_MAX = 1000

def make_payload(pid):
    return {
        "product_id": pid,
        "sku": f"SKU-{pid:04d}",
        "manufacturer": "Acme",
        "category_id": 10,
        "weight": 100 + (pid % 50),
        "some_other_id": 9
    }

# -------------------------
# Regular HTTP client
# -------------------------
class StoreUser(HttpUser):
    wait_time = between(0.05, 0.2)

    @task(5)
    def get_product(self):
        pid = random.randint(PRODUCT_ID_MIN, PRODUCT_ID_MAX)
        self.client.get(f"/products/{pid}", name="/products/:id")

    @task(1)
    def post_details(self):
        pid = random.randint(PRODUCT_ID_MIN, PRODUCT_ID_MAX)
        payload = make_payload(pid)
        self.client.post(f"/products/{pid}/details",
                         json=payload, name="/products/:id/details")

# -------------------------
# FastHTTP client
# -------------------------
class StoreUserFast(FastHttpUser):
    wait_time = between(0.05, 0.2)

    @task(5)
    def get_product(self):
        pid = random.randint(PRODUCT_ID_MIN, PRODUCT_ID_MAX)
        self.client.get(f"/products/{pid}", name="/products/:id")

    @task(1)
    def post_details(self):
        pid = random.randint(PRODUCT_ID_MIN, PRODUCT_ID_MAX)
        payload = make_payload(pid)
        self.client.post(f"/products/{pid}/details",
                         json=payload, name="/products/:id/details")
