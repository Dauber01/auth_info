"""Auth API contract cases; see docs/tasks/20260912-docs-task-tests/04-test-cases.md."""

import os
import unittest

from tests.api.client import APIClient


class AuthAPITest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        base_url = os.environ.get("API_BASE_URL")
        if not base_url:
            raise RuntimeError("API_BASE_URL is required; run python3 -B tests/run_api.py")
        cls.client = APIClient(base_url)

    def assert_reply(self, response, status):
        self.assertEqual(response.status, status)
        self.assertIn("application/json", response.content_type)
        body = response.json()
        self.assertEqual(body["code"], status)
        self.assertIsInstance(body["message"], str)
        self.assertTrue(body["message"])
        return body

    def register(self, username):
        return self.client.request("POST", "/api/v1/auth/register",
                                   {"username": username, "password": "fixture-password"})

    def test_register_login_and_access_protected_hello(self):
        self.assert_reply(self.register("api_success"), 200)
        response = self.client.request("POST", "/api/v1/auth/login",
                                       {"username": "api_success", "password": "fixture-password"})
        token = self.assert_reply(response, 200)["data"]["token"]
        self.assertIsInstance(token, str)
        self.assertTrue(token)
        hello = self.client.request("GET", "/api/v1/hello?name=API", token=token)
        self.assertEqual(self.assert_reply(hello, 200)["data"]["message"], "Hello, API!")

    def test_duplicate_username_is_conflict(self):
        self.assert_reply(self.register("api_duplicate"), 200)
        self.assert_reply(self.register("api_duplicate"), 409)

    def test_invalid_registration_is_rejected(self):
        for payload in ({}, {"username": "ab", "password": "123456"},
                        {"username": "api_shortpass", "password": "12345"}):
            with self.subTest(payload=payload):
                self.assert_reply(self.client.request("POST", "/api/v1/auth/register", payload), 400)

    def test_missing_login_fields_are_rejected(self):
        self.assert_reply(self.client.request("POST", "/api/v1/auth/login", {}), 400)

    def test_wrong_password_and_unknown_user_are_unauthenticated(self):
        self.assert_reply(self.register("api_wrong_password"), 200)
        for username in ("api_wrong_password", "api_unknown_user"):
            with self.subTest(username=username):
                response = self.client.request("POST", "/api/v1/auth/login",
                                               {"username": username, "password": "incorrect-password"})
                self.assert_reply(response, 401)
                self.assertNotIn("token", response.json().get("data", {}))

    def test_protected_route_requires_token(self):
        self.assert_reply(self.client.request("GET", "/api/v1/hello"), 401)

    def test_protected_route_rejects_invalid_token(self):
        self.assert_reply(self.client.request("GET", "/api/v1/hello", token="invalid-fixture-token"), 401)
