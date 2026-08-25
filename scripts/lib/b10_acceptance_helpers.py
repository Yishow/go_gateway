"""Transport and process helpers for the B10 executable acceptance harness."""

from __future__ import annotations

import hashlib
import json
import os
import shutil
import signal
import socket
import subprocess
import threading
import time
import urllib.error
import urllib.request
from urllib.parse import urlencode
from pathlib import Path
from typing import Any


def free_port() -> int:
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        return int(sock.getsockname()[1])
class SourceModbus:
    """Minimal real TCP Modbus source used by the collector, not an in-process fake."""

    def __init__(self, port: int, value: int = 0x1234):
        self.port = port
        self.value = value
        self.stop_event = threading.Event()
        self.listener: socket.socket | None = None
        self.thread = threading.Thread(target=self._run, daemon=True)

    def start(self) -> None:
        self.thread.start()
        deadline = time.time() + 5
        while time.time() < deadline:
            try:
                with socket.create_connection(("127.0.0.1", self.port), 0.2):
                    return
            except OSError:
                time.sleep(0.02)
        raise RuntimeError("source Modbus server did not start")

    def stop(self) -> None:
        self.stop_event.set()
        if self.listener is not None:
            self.listener.close()
        if self.thread.is_alive():
            self.thread.join(timeout=2)

    def _run(self) -> None:
        with socket.socket() as listener:
            self.listener = listener
            listener.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            listener.bind(("127.0.0.1", self.port))
            listener.listen(8)
            listener.settimeout(0.2)
            while not self.stop_event.is_set():
                try:
                    conn, _ = listener.accept()
                except (TimeoutError, socket.timeout):
                    continue
                except OSError:
                    break
                threading.Thread(target=self._connection, args=(conn,), daemon=True).start()

    def _connection(self, conn: socket.socket) -> None:
        with conn:
            conn.settimeout(2)
            while not self.stop_event.is_set():
                try:
                    header = recv_exact(conn, 7)
                    length = int.from_bytes(header[4:6], "big")
                    pdu = recv_exact(conn, length - 1)
                except (OSError, EOFError):
                    return
                if not pdu or pdu[0] != 3 or len(pdu) < 5:
                    return
                quantity = int.from_bytes(pdu[3:5], "big")
                if quantity < 1 or quantity > 125:
                    return
                data = b"".join((self.value & 0xFFFF).to_bytes(2, "big") for _ in range(quantity))
                response_pdu = bytes([3, len(data)]) + data
                response = header[:2] + b"\x00\x00" + (len(response_pdu) + 1).to_bytes(2, "big") + header[6:7] + response_pdu
                conn.sendall(response)
def recv_exact(conn: socket.socket, count: int) -> bytes:
    chunks: list[bytes] = []
    while count:
        chunk = conn.recv(count)
        if not chunk:
            raise EOFError
        chunks.append(chunk)
        count -= len(chunk)
    return b"".join(chunks)


def read_holding(host: str, port: int, unit: int = 1, address: int = 0, quantity: int = 1) -> list[int]:
    transaction = int(time.time() * 1000) & 0xFFFF
    pdu = bytes([3]) + address.to_bytes(2, "big") + quantity.to_bytes(2, "big")
    request = transaction.to_bytes(2, "big") + b"\x00\x00" + (len(pdu) + 1).to_bytes(2, "big") + bytes([unit]) + pdu
    with socket.create_connection((host, port), 2) as conn:
        conn.sendall(request)
        header = recv_exact(conn, 7)
        response = recv_exact(conn, int.from_bytes(header[4:6], "big") - 1)
    if response[0] != 3 or response[1] != quantity * 2:
        raise RuntimeError(f"unexpected Modbus response: {response.hex()}")
    return [int.from_bytes(response[i : i + 2], "big") for i in range(2, len(response), 2)]
class HTTP:
    def __init__(self, base: str):
        self.base = base.rstrip("/")
        self.commands: list[str] = []

    def call(self, method: str, path: str, body: dict[str, Any] | None = None, expected: set[int] | None = None) -> tuple[int, dict[str, Any]]:
        url = self.base + path
        data = None if body is None else json.dumps(body, separators=(",", ":")).encode()
        self.commands.append(f"curl -sS -X {method} {url}" + (" -d '<json>'" if data else ""))
        request = urllib.request.Request(url, data=data, method=method, headers={"Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(request, timeout=8) as response:
                status = response.status
                raw = response.read()
        except urllib.error.HTTPError as error:
            status = error.code
            raw = error.read()
        try:
            payload = json.loads(raw.decode() or "{}")
        except json.JSONDecodeError:
            payload = {"raw": raw.decode(errors="replace")}
        if expected is not None and status not in expected:
            raise RuntimeError(f"{method} {path} returned {status}: {payload}")
        return status, payload


def data(payload: dict[str, Any]) -> dict[str, Any]:
    value = payload.get("data")
    if not isinstance(value, dict):
        raise RuntimeError(f"response has no object data: {payload}")
    return value


def list_data(payload: dict[str, Any]) -> list[dict[str, Any]]:
    value = payload.get("data")
    if not isinstance(value, list):
        raise RuntimeError(f"response has no list data: {payload}")
    if any(not isinstance(item, dict) for item in value):
        raise RuntimeError(f"response list data contains a non-object mapping: {payload}")
    return value


def candidate_scope_path(path: str, workspace_id: str, workspace_revision: str, revision_id: str) -> str:
    """Build a candidate endpoint path with the production workspace scope contract."""
    return f"{path}?{urlencode({'workspace_id': workspace_id, 'expected_workspace_revision': workspace_revision, 'revision_id': revision_id})}"


def redact_witness(value: Any) -> Any:
    if isinstance(value, dict):
        return {
            key: redact_witness(item)
            for key, item in value.items()
            if not any(secret in key.lower() for secret in ("token", "password", "secret", "credential", "dsn"))
        }
    if isinstance(value, list):
        return [redact_witness(item) for item in value]
    return value


def wait_gateway(http: HTTP) -> None:
    deadline = time.time() + 20
    while time.time() < deadline:
        try:
            status, _ = http.call("GET", "/api/v1/datalink/studio-v2/workspace", expected={200})
            if status == 200:
                return
        except (OSError, urllib.error.URLError, RuntimeError):
            pass
        time.sleep(0.1)
    raise RuntimeError("gateway HTTP server did not become ready")


def stop_gateway(proc: subprocess.Popen[bytes], log_path: Path) -> bool:
    process_code = proc.poll()
    if process_code is not None:
        _close_process_log(proc)
        raise RuntimeError(f"gateway exited early with code {process_code}; see {log_path}")
    try:
        proc.send_signal(signal.SIGTERM)
    except ProcessLookupError:
        _close_process_log(proc)
        return True
    try:
        proc.wait(timeout=10)
    except subprocess.TimeoutExpired:
        try:
            proc.kill()
        except ProcessLookupError:
            _close_process_log(proc)
            return True
        try:
            proc.wait(timeout=3)
        finally:
            _close_process_log(proc)
        return False
    finally:
        _close_process_log(proc)
    return True


def _close_process_log(proc: subprocess.Popen[bytes]) -> None:
    log = getattr(proc, "_b10_log", None)
    if log is None:
        return
    try:
        log.close()
    except Exception:
        pass


def _safe_finalize_process(proc: subprocess.Popen[bytes], log_path: Path) -> list[str]:
    errors: list[str] = []
    try:
        if proc.poll() is None:
            try:
                stop_gateway(proc, log_path)
            except Exception as error:
                errors.append(f"stop_gateway: {error}")
                try:
                    if proc.poll() is None:
                        try:
                            proc.kill()
                        except (ProcessLookupError, OSError) as kill_error:
                            errors.append(f"kill: {kill_error}")
                        try:
                            proc.wait(timeout=3)
                        except (subprocess.TimeoutExpired, ProcessLookupError, OSError) as wait_error:
                            errors.append(f"wait: {wait_error}")
                except Exception as fallback_error:
                    errors.append(f"process fallback: {fallback_error}")
    except Exception as error:
        errors.append(f"process poll: {error}")
    finally:
        _close_process_log(proc)
    return errors


def start_gateway(binary: Path, db_path: Path, http_port: int, log_path: Path, root: Path) -> subprocess.Popen[bytes]:
    env = os.environ.copy()
    env.update({
        "GATEWAY_DB_PATH": str(db_path),
        "HOST": "127.0.0.1",
        "PORT": str(http_port),
        "AUTO_OPEN_BROWSER": "false",
        "GIN_MODE": "release",
    })
    log = log_path.open("ab")
    proc = subprocess.Popen([str(binary)], cwd=root, env=env, stdout=log, stderr=subprocess.STDOUT)
    proc._b10_log = log  # type: ignore[attr-defined]
    return proc


def finalize_b10_witness(
    proc: subprocess.Popen[bytes] | None,
    log_path: Path,
    source: SourceModbus,
    result: dict[str, Any],
    command_log: list[str],
    current_http: HTTP | None,
    change: Path,
    temp_root: Path,
) -> None:
    cleanup_errors: list[str] = []
    if proc is not None:
        cleanup_errors.extend(_safe_finalize_process(proc, log_path))
    try:
        source.stop()
    except Exception as error:
        cleanup_errors.append(f"source.stop: {error}")
    if log_path.exists():
        result["log_sha256"] = hashlib.sha256(log_path.read_bytes()).hexdigest()
    result["log_path"] = str(log_path)
    if isinstance(current_http, HTTP):
        command_log.extend(current_http.commands)
    result["commands"] = list(dict.fromkeys(command_log))
    result["finished_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    if result.get("status") != "pass":
        result.setdefault("status", "fail")
        if log_path.exists():
            result["failure_log_tail"] = log_path.read_text(encoding="utf-8", errors="replace")[-4000:]
    if cleanup_errors:
        result["cleanup_errors"] = cleanup_errors
    witness_dir = change / "witness"
    witness_dir.mkdir(parents=True, exist_ok=True)
    witness = witness_dir / f"b10-exe-acceptance-{time.strftime('%Y%m%dT%H%M%SZ', time.gmtime())}.json"
    safe_result = redact_witness(result)
    safe_result["witness"] = str(witness)
    witness.write_text(json.dumps(safe_result, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    witnesses = sorted(witness_dir.glob("b10-exe-acceptance-*.json"), key=lambda path: path.stat().st_mtime, reverse=True)
    for superseded in witnesses[2:]:
        superseded.unlink()
    result.clear()
    result.update(safe_result)
    shutil.rmtree(temp_root, ignore_errors=True)
