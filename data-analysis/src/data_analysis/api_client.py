import requests
from typing import Optional


class APIClient:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")

    def search_uris(self, query: str, limit: int = 10) -> list[str]:
        url = f"{self.base_url}/api/search"
        params = {"query": query, "limit": limit}
        resp = requests.get(url, params=params, timeout=30)
        resp.raise_for_status()
        return resp.json()

    def analyze_by_uri(self, uri: str, model: str) -> dict:
        url = f"{self.base_url}/api/analysis/by-uri"
        params = {"uri": uri, "model": model}
        resp = requests.get(url, params=params, timeout=60)
        resp.raise_for_status()
        return resp.json()

    def health_check(self) -> bool:
        url = f"{self.base_url}/health"
        resp = requests.get(url, timeout=10)
        return resp.status_code == 200
