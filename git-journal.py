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

if command == "init":
    remote = sys.argv[3]

    if os.path.exists(path_md) or os.path.exists(path_json):
        print("Journal files already exist")
        sys.exit(1)

    with open(path_json, "w+") as json_file:
        json_file.write("{}")

    with open(path_md, "w+") as md_file:
        md_file.write("")

    subprocess.run(["git", "init"])
    subprocess.run(["git", "add", JOURNAL_JSON])
    subprocess.run(["git", "add", JOURNAL_MD])
    subprocess.run(["git", "commit", "-m", "Initial commit"])
    subprocess.run(["git", "remote", "add", "origin", remote])
    subprocess.run(["git", "push", "-u", "origin", "master"])

elif command == "entry":
    subprocess.run(["git", "checkout", "master"])
    subprocess.run(["git", "branch", "-D", "journal-backup"])
    subprocess.run(["git", "branch", "journal-backup"])

    editor = os.environ.get("EDITOR", "vim")

    with tempfile.NamedTemporaryFile(suffix=".tmp") as tf:
        tf.flush()
        subprocess.run([editor, tf.name])
        tf.seek(0)

        entry_id = str(uuid.uuid4())
        entry_timestamp = time.strftime(TIME_FORMAT)
        entry_body = tf.read().decode("utf-8")

    with open(path_json, "r+") as json_file:
        local_json = json.loads(json_file.read())
        local_json[entry_id] = {
            "timestamp": entry_timestamp,
            "body": entry_body
        }

    subprocess.run(["git", "fetch"])
    subprocess.run(["git", "reset", "--hard", "origin/master"])
    with open(path_json, "r+") as json_file:
        combined_json = json.loads(json_file.read())
        combined_json.update(local_json)

        json_file.seek(0)
        json_file.write(json.dumps(combined_json))
        json_file.truncate()

    with open(path_md, "w") as md_file:
        entries = sorted(
            combined_json.values(),
            key=lambda e: time.strptime(e["timestamp"], TIME_FORMAT),
            reverse=True)
        formatted = "\n".join([
            "## {}\n\n{}\n".format(entry["timestamp"], entry["body"])
            for entry in entries
        ])
        md_file.write(formatted)

    subprocess.run(["git", "add", JOURNAL_JSON])
    subprocess.run(["git", "add", JOURNAL_MD])
    title = entry_body.split("\n")[0]
    subprocess.run(
        ["git", "commit", "-m", title[:75] + (title[75:] and '...')])
    subprocess.run(["git", "push"])
