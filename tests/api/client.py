"""Small standard-library HTTP client that preserves API error responses."""

import json
from dataclasses import dataclass
from urllib.error import HTTPError
from urllib.parse import urlsplit
from urllib.request import HTTPRedirectHandler, ProxyHandler, Request, build_opener


class NoRedirects(HTTPRedirectHandler):
    def redirect_request(self, request, fp, code, message, headers, new_url):
        return None


@dataclass(frozen=True)
class Response:
    status: int
    content_type: str
    body: bytes

    def json(self):
        return json.loads(self.body)


class APIClient:
    def __init__(self, base_url):
        url = urlsplit(base_url)
        if (url.scheme not in {"http", "https"} or not url.netloc or url.username or
                url.password or url.query or url.fragment):
            raise ValueError("API_BASE_URL must be an HTTP(S) service address without credentials or query")
        self.base_url = base_url.rstrip("/")
        self.opener = build_opener(ProxyHandler({}), NoRedirects())

    def request(self, method, path, body=None, token=None):
        if not path.startswith("/") or path.startswith("//"):
            raise ValueError("API request path must be service-relative")
        headers = {"Accept": "application/json"}
        if body is not None:
            headers["Content-Type"] = "application/json"
        if token is not None:
            headers["Authorization"] = "Bearer " + token
        data = json.dumps(body).encode("utf-8") if body is not None else None
        request = Request(self.base_url + path, data=data, headers=headers, method=method)
        try:
            response = self.opener.open(request, timeout=10)
        except HTTPError as error:
            response = error
        with response:
            return Response(response.code, response.headers.get("Content-Type", ""), response.read())
