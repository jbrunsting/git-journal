#!/usr/bin/env python3

import json
import os
import subprocess
import sys
import tempfile
import time
import uuid

TIME_FORMAT = "%x %X %Z"
JOURNAL_MD = "journal.md"
JOURNAL_JSON = "journal.json"

command = sys.argv[1]
directory = sys.argv[2]
os.chdir(directory)
path_md = "{}/{}".format(directory, JOURNAL_MD)
path_json = "{}/{}".format(directory, JOURNAL_JSON)


def files_exist():
    return os.path.exists(path_md) and os.path.exists(path_json)


def commit_entry(entry_id, entry_timestamp, entry_body):
    with open(path_json, "r+") as json_file:
        entries = json.loads(json_file.read())
        entries[entry_id] = {
            "timestamp": entry_timestamp,
            "body": entry_body
        }
        json_file.seek(0)
        json_file.write(json.dumps(entries))
        json_file.truncate()

    with open(path_md, "w") as md_file:
        sorted_entries = sorted(
            entries.values(),
            key=lambda e: time.strptime(e["timestamp"], TIME_FORMAT),
            reverse=True)
        formatted = "\n".join([
            "## {}\n\n{}\n".format(entry["timestamp"], entry["body"])
            for entry in sorted_entries
        ])
        md_file.write(formatted)

    subprocess.run(["git", "add", JOURNAL_JSON])
    subprocess.run(["git", "add", JOURNAL_MD])
    title = entry_body.split("\n")[0]
    subprocess.run(
        ["git", "commit", "-m", title[:75] + (title[75:] and '...')])


def update_from_origin():
    with open(path_json, "r+") as json_file:
        local_json = json.loads(json_file.read())

    subprocess.run(["git", "fetch"])
    subprocess.run(["git", "reset", "--hard", "origin/master"])
    with open(path_json, "r+") as json_file:
        origin_json = json.loads(json_file.read())
        if origin_json == local_json:
            return False

    new_entries = {k: v for k, v in local_json.items() if k not in origin_json}
    for entry_id, entry in new_entries.items():
        commit_entry(entry_id, entry["timestamp"], entry["body"])

    return True


def init():
    remote = sys.argv[3]

    if files_exist():
        print("Journal files already exist")
        sys.exit(1)

    subprocess.run(["git", "init"])
    subprocess.run(["git", "remote", "add", "origin", remote])
    subprocess.run(["git", "pull", "origin", "master"])
    subprocess.run(["git", "push", "--set-upstream", "origin", "master"])

    if not files_exist():
        with open(path_json, "w") as json_file:
            json_file.write("{}")

        with open(path_md, "w") as md_file:
            md_file.write("")

        subprocess.run(["git", "add", JOURNAL_JSON])
        subprocess.run(["git", "add", JOURNAL_MD])
        subprocess.run(["git", "commit", "-m", "Initial commit"])
        subprocess.run(["git", "push", "-u", "origin", "master"])


def add_entry():
    if not files_exist():
        print("Journal files do not exist")
        sys.exit(1)

    editor = os.environ.get("EDITOR", "vim")

    with tempfile.NamedTemporaryFile(suffix=".tmp") as tf:
        tf.flush()
        subprocess.run([editor, tf.name])
        tf.seek(0)
        entry_body = tf.read().decode("utf-8")

    commit_entry(str(uuid.uuid4()), time.strftime(TIME_FORMAT), entry_body)


def sync():
    if not files_exist():
        print("Journal files do not exist")
        sys.exit(1)

    subprocess.run(["git", "checkout", "master"])
    subprocess.run(["git", "branch", "-D", "journal-backup"])
    subprocess.run(["git", "branch", "journal-backup"])

    if update_from_origin():
        subprocess.run(["git", "push"])
    else:
        print("Already up to date")


if command == "init":
    init()
elif command == "entry":
    add_entry()
elif command == "sync":
    sync()
elif command == "pentry":
    sync()
    add_entry()
    subprocess.run(["git", "push"])
