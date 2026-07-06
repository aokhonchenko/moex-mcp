"""Subprocess runners used by transactional sessions."""

from __future__ import annotations

import os
import queue
import subprocess
import threading
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Sequence


@dataclass(frozen=True)
class CommandResult:
    returncode: int
    stdout: str = ""
    stderr: str = ""


class CommandExecutionError(RuntimeError):
    """Raised when a command runner cannot execute a subprocess."""


HEARTBEAT_SECONDS = 30.0
SUBPROCESS_ENCODING = "utf-8"
SUBPROCESS_ERRORS = "replace"


def subprocess_env() -> dict[str, str]:
    env = os.environ.copy()
    env.pop("VIRTUAL_ENV", None)
    return env


def default_runner(args: Sequence[str], cwd: Path) -> CommandResult:
    completed = subprocess.run(
        list(args),
        cwd=cwd,
        env=subprocess_env(),
        text=True,
        encoding=SUBPROCESS_ENCODING,
        errors=SUBPROCESS_ERRORS,
        capture_output=True,
    )
    return CommandResult(completed.returncode, completed.stdout, completed.stderr)


def run_streaming_command(
    args: Sequence[str],
    cwd: Path,
    heartbeat_seconds: float = HEARTBEAT_SECONDS,
) -> CommandResult:
    command = subprocess.list2cmdline(list(args))
    print(f"[cmd] {cwd}> {command}", flush=True)
    process = subprocess.Popen(
        list(args),
        cwd=cwd,
        env=subprocess_env(),
        text=True,
        encoding=SUBPROCESS_ENCODING,
        errors=SUBPROCESS_ERRORS,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        bufsize=1,
    )
    if process.stdout is None:
        raise CommandExecutionError("streaming command did not expose stdout")

    output_queue: queue.Queue[str] = queue.Queue()

    def read_output() -> None:
        try:
            for line in iter(process.stdout.readline, ""):
                output_queue.put(line)
        finally:
            process.stdout.close()

    reader = threading.Thread(target=read_output, daemon=True)
    reader.start()

    output: list[str] = []
    last_activity = time.monotonic()

    def drain_output() -> None:
        while True:
            try:
                line = output_queue.get_nowait()
            except queue.Empty:
                return
            output.append(line)
            print(line, end="", flush=True)

    while True:
        try:
            line = output_queue.get(timeout=0.2)
        except queue.Empty:
            if process.poll() is not None:
                reader.join(timeout=1)
                drain_output()
                break
            now = time.monotonic()
            if heartbeat_seconds > 0 and now - last_activity >= heartbeat_seconds:
                print(f"[wait] command is still running: {command}", flush=True)
                last_activity = now
            continue

        output.append(line)
        print(line, end="", flush=True)
        last_activity = time.monotonic()

    return CommandResult(process.wait(), "".join(output), "")


def streaming_runner(args: Sequence[str], cwd: Path) -> CommandResult:
    return run_streaming_command(args, cwd)