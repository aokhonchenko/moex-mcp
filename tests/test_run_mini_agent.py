import importlib
import os
from pathlib import Path

import pytest

from scripts import run_mini_agent


def test_read_settings_returns_defaults_for_missing_file(tmp_path):
    settings = run_mini_agent.read_settings(tmp_path / "missing.toml")

    assert settings["custom_llm_provider"] == "openai"
    assert settings["cost_tracking"] == "ignore_errors"


def test_read_settings_merges_project_values(tmp_path):
    path = tmp_path / "project.toml"
    path.write_text(
        '[mini_swe_agent]\nstep_limit = 3\n',
        encoding="utf-8",
    )

    settings = run_mini_agent.read_settings(path)

    assert settings["step_limit"] == 3
    assert settings["custom_llm_provider"] == "openai"


def test_require_env_rejects_missing_value(monkeypatch):
    monkeypatch.delenv("AI_MODEL", raising=False)

    with pytest.raises(RuntimeError, match="AI_MODEL"):
        run_mini_agent.require_env("AI_MODEL")


def test_build_model_kwargs_uses_ai_environment(monkeypatch):
    monkeypatch.setenv("AI_API_KEY", "secret")
    monkeypatch.setenv("AI_BASE_URL", "https://example.test/v1")
    monkeypatch.setenv("AI_MODEL", "openai/test-model")

    kwargs = run_mini_agent.build_model_kwargs(
        {"custom_llm_provider": "openai"}
    )

    assert kwargs == {
        "custom_llm_provider": "openai",
        "api_base": "https://example.test/v1",
        "api_key": "secret",
        "drop_params": True,
    }


def test_build_model_kwargs_uses_dummy_key_for_empty_ai_api_key(monkeypatch):
    monkeypatch.delenv("AI_API_KEY", raising=False)
    monkeypatch.setenv("AI_BASE_URL", "https://example.test/v1")

    kwargs = run_mini_agent.build_model_kwargs(
        {"custom_llm_provider": "openai"}
    )

    assert kwargs["api_key"] == "not-needed"


def test_read_task_requires_existing_prompt(tmp_path):
    with pytest.raises(RuntimeError, match="Prompt file does not exist"):
        run_mini_agent.read_task(tmp_path / "missing.md")

    prompt = tmp_path / "prompt.md"
    prompt.write_text("задача", encoding="utf-8")
    assert run_mini_agent.read_task(prompt) == "задача"


def test_main_returns_one_on_agent_error(tmp_path, monkeypatch, capsys):
    prompt = tmp_path / "prompt.md"
    prompt.write_text("task", encoding="utf-8")

    def fail_run_agent(root, prompt_file, settings_path):
        raise RuntimeError("configuration problem")

    monkeypatch.setattr(run_mini_agent, "run_agent", fail_run_agent)
    monkeypatch.setattr(
        run_mini_agent.sys,
        "argv",
        ["run_mini_agent.py", "--root", str(tmp_path), "--prompt-file", str(prompt)],
    )

    result = run_mini_agent.main()

    captured = capsys.readouterr()
    assert result == 1
    assert "mini-swe-agent session failed: configuration problem" in captured.err

def test_require_env_accepts_ai_model(monkeypatch):
    monkeypatch.setenv("AI_MODEL", "openai/test-model")

    assert run_mini_agent.require_env("AI_MODEL") == "openai/test-model"