"""Small OpenAI-compatible chat completions client."""

from __future__ import annotations

import json
import os
import urllib.error
import urllib.request
from typing import Any


class LlmClientError(RuntimeError):
    """Raised when an OpenAI-compatible request fails."""


DEFAULT_TIMEOUT_SECONDS = 300


def optional_env(name: str) -> str:
    return os.environ.get(name, "").strip()


def require_env(name: str) -> str:
    value = optional_env(name)
    if not value:
        raise LlmClientError(f"Environment variable {name} is required")
    return value


def chat_completions_url(base_url: str) -> str:
    normalized = base_url.rstrip("/")
    if normalized.endswith("/chat/completions"):
        return normalized
    return f"{normalized}/chat/completions"


def request_headers(api_key: str) -> dict[str, str]:
    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    return headers


def post_chat_completion(
    *,
    base_url: str,
    api_key: str,
    model: str,
    messages: list[dict[str, Any]],
    tools: list[dict[str, Any]],
    temperature: float = 0.2,
    timeout: int = DEFAULT_TIMEOUT_SECONDS,
) -> dict[str, Any]:
    body = {
        "model": model,
        "messages": messages,
        "tools": tools,
        "tool_choice": "auto",
        "temperature": temperature,
    }
    payload = json.dumps(body).encode("utf-8")
    request = urllib.request.Request(
        chat_completions_url(base_url),
        data=payload,
        headers=request_headers(api_key),
        method="POST",
    )
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            raw = response.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="replace")
        raise LlmClientError(f"chat completion HTTP {exc.code}: {details}") from exc
    except urllib.error.URLError as exc:
        raise LlmClientError(f"chat completion request failed: {exc}") from exc

    try:
        return json.loads(raw)
    except json.JSONDecodeError as exc:
        raise LlmClientError(f"chat completion returned invalid JSON: {raw[:500]}") from exc


def extract_message(response: dict[str, Any]) -> dict[str, Any]:
    choices = response.get("choices")
    if not isinstance(choices, list) or not choices:
        raise LlmClientError("chat completion response does not contain choices")
    message = choices[0].get("message")
    if not isinstance(message, dict):
        raise LlmClientError("chat completion choice does not contain a message")
    return message


class OpenAICompatibleClient:
    def __init__(self, *, base_url: str, api_key: str, model: str, timeout: int, temperature: float):
        self.base_url = base_url
        self.api_key = api_key
        self.model = model
        self.timeout = timeout
        self.temperature = temperature

    @classmethod
    def from_environment(cls, *, timeout: int, temperature: float) -> "OpenAICompatibleClient":
        return cls(
            base_url=require_env("AI_BASE_URL"),
            api_key=optional_env("AI_API_KEY"),
            model=require_env("AI_MODEL"),
            timeout=timeout,
            temperature=temperature,
        )

    def complete(self, messages: list[dict[str, Any]], tools: list[dict[str, Any]]) -> dict[str, Any]:
        response = post_chat_completion(
            base_url=self.base_url,
            api_key=self.api_key,
            model=self.model,
            messages=messages,
            tools=tools,
            temperature=self.temperature,
            timeout=self.timeout,
        )
        return extract_message(response)