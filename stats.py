#!/usr/bin/env python3

import json
import os
import re
from pathlib import Path

EXTENSIONS = {".go", ".js", ".jsx", ".css", ".sql"}
FOLDERS = ["frontend", "backend"]
EXCLUDED_DIRS = {"node_modules", ".git", "dist", "build", "vendor", "__pycache__", ".next"}

# Rough "function-like construct" patterns per extension.
FUNC_PATTERNS = {
    ".go": re.compile(r"^\s*func\b", re.MULTILINE),
    ".js": re.compile(r"\bfunction\b|=>"),
    ".jsx": re.compile(r"\bfunction\b|=>"),
}

TODO_PATTERN = re.compile(r"todo", re.IGNORECASE)
USE_CLIENT_PATTERN = re.compile(r"""^\s*['"]use client['"]""", re.MULTILINE)


def read_file(filepath):
    try:
        with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
            return f.read()
    except OSError:
        return ""


def count_lines(content):
    if not content:
        return 0
    return content.count("\n") + (1 if content and not content.endswith("\n") else 0)


def count_functions(content, ext):
    pattern = FUNC_PATTERNS.get(ext)
    if not pattern:
        return 0
    return len(pattern.findall(content))


def count_todos(content):
    return len(TODO_PATTERN.findall(content))


def count_package_json_deps(folder):
    path = os.path.join(folder, "package.json")
    if not os.path.isfile(path):
        return 0
    try:
        with open(path, "r", encoding="utf-8", errors="ignore") as f:
            data = json.load(f)
    except (OSError, json.JSONDecodeError):
        return 0
    deps = len(data.get("dependencies", {}))
    dev_deps = len(data.get("devDependencies", {}))
    return deps + dev_deps


def count_go_mod_deps(folder):
    path = os.path.join(folder, "go.mod")
    if not os.path.isfile(path):
        return 0
    content = read_file(path)
    # Matches lines like: "	github.com/foo/bar v1.2.3"
    entries = re.findall(r"^\s+[^\s]+\s+v\d+\.\d+\.\d+", content, re.MULTILINE)
    return len(entries)


def count_pages(folder):
    app_dir = os.path.join(folder, "app")
    if not os.path.isdir(app_dir):
        return 0
    count = 0
    for root, dirs, files in os.walk(app_dir):
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
        for name in files:
            stem = Path(name).stem
            ext = Path(name).suffix
            if stem == "page" and ext in (".js", ".jsx"):
                count += 1
    return count


def count_client_server_components(folder):
    client_count = 0
    server_count = 0
    for root, dirs, files in os.walk(folder):
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
        for name in files:
            ext = Path(name).suffix
            if ext in (".js", ".jsx"):
                content = read_file(os.path.join(root, name))
                if USE_CLIENT_PATTERN.search(content):
                    client_count += 1
                else:
                    server_count += 1
    return client_count, server_count


def analyze_folder(folder):
    stats_by_ext = {ext: {"files": 0, "lines": 0} for ext in EXTENSIONS}
    total_files = 0
    total_lines = 0
    total_functions = 0
    total_todos = 0

    if not os.path.isdir(folder):
        print(f"Folder '{folder}' not found, skipping.\n")
        return {
            "files": 0, "lines": 0, "functions": 0, "todos": 0,
            "package_deps": 0,
        }

    for root, dirs, files in os.walk(folder):
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
        for name in files:
            ext = Path(name).suffix
            if ext in EXTENSIONS:
                filepath = os.path.join(root, name)
                content = read_file(filepath)
                lines = count_lines(content)
                functions = count_functions(content, ext)
                todos = count_todos(content)

                stats_by_ext[ext]["files"] += 1
                stats_by_ext[ext]["lines"] += lines
                total_files += 1
                total_lines += lines
                total_functions += functions
                total_todos += todos

    print(f"{folder}/")
    print("   By extension:")
    for ext in sorted(EXTENSIONS):
        f_count = stats_by_ext[ext]["files"]
        l_count = stats_by_ext[ext]["lines"]
        if f_count > 0:
            print(f"     {ext:6} -> {f_count:4} files, {l_count:6} lines")
    print(f"[total] -> {total_files} files, {total_lines} lines")

    package_deps = count_package_json_deps(folder)

    # Frontend-specific stats
    if folder == "frontend":
        pages = count_pages(folder)
        client_count, server_count = count_client_server_components(folder)
        print(f"   Pages/routes: {pages}")
        print(f"   Client components: {client_count}")
        print(f"   Server components: {server_count}")

    # Backend-specific stats
    if folder == "backend":
        go_deps = count_go_mod_deps(folder)
        print(f"   go.mod dependencies: {go_deps}")

    print()

    return {
        "files": total_files,
        "lines": total_lines,
        "functions": total_functions,
        "todos": total_todos,
        "package_deps": package_deps,
    }


def main():
    print("=== Code Stats ===\n")
    grand = {"files": 0, "lines": 0, "functions": 0, "todos": 0, "package_deps": 0}

    for folder in FOLDERS:
        result = analyze_folder(folder)
        for key in grand:
            grand[key] += result[key]

    print("=== Grand Total ===")
    print(f"   Total files: {grand['files']}")
    print(f"   Total lines: {grand['lines']}")
    print(f"   Total functions: {grand['functions']}")
    print(f"   TODO count: {grand['todos']}")
    print(f"   package.json dependencies: {grand['package_deps']}")


if __name__ == "__main__":
    main()