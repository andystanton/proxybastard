#!/usr/bin/env python3

from datetime import datetime
import json
import os.path
import subprocess
import argparse

goxc_file = ".goxc.json"
goxc_local_file = ".goxc.local.json"

def validate_goxc():
    if os.path.isfile(goxc_file):
        goxc_contents = open(goxc_file).read()
        goxc_json = json.loads(goxc_contents)

        if "TaskSettings" not in goxc_json:
            print("Missing TaskSettings in {0}".format(goxc_file))
            exit(1)
        else:
            task_settings_json = goxc_json["TaskSettings"]
            if "bintray" not in task_settings_json:
                print("Missing bintray in {0}".format(goxc_file))
                exit(1)
            else:
                bintray_json = task_settings_json["bintray"]
                if not all(k in bintray_json for k in ("downloadspage", "package", "repository", "subject", "user")):
                    print("Invalid bintray configuration in {0}".format(goxc_file))
                    print(json.dumps(bintray_json, indent=4, separators=(',', ': ')))
                    exit(1)
    else :
        print("Missing {0}".format(goxc_file))
        exit(1)

def validate_goxc_local():
    if os.path.isfile(goxc_local_file):
        goxc_local_contents = open(goxc_local_file).read()
        goxc_local_json = json.loads(goxc_local_contents)

        if "TaskSettings" not in goxc_local_json:
            print("Missing TaskSettings in {0}".format(goxc_local_file))
            exit(1)
        else:
            task_settings_json = goxc_local_json["TaskSettings"]
            if "bintray" not in task_settings_json:
                print("Missing bintray in {0}".format(goxc_local_file))
                exit(1)
            else:
                bintray_json = task_settings_json["bintray"]
                if "apikey" not in bintray_json:
                    print("Invalid bintray configuration in {0}".format(goxc_local_file))
                    print(json.dumps(bintray_json, indent=4, separators=(',', ': ')))
                    exit(1)
    else :
        print("Missing {0}".format(goxc_local_file))
        exit(1)

def set_snapshot_build_info():
    goxc_local_contents = open(goxc_local_file).read()
    goxc_local_json = json.loads(goxc_local_contents)
    goxc_local_json["BuildName"] = datetime.now().strftime("%Y%m%d%H%M%S")
    goxc_local_json["PrereleaseInfo"] = "snapshot"
    open(goxc_local_file, 'w').write(json.dumps(goxc_local_json, indent=4, separators=(',', ': ')) + "\n")

def remove_snapshot_build_info():
    goxc_local_contents = open(goxc_local_file).read()
    goxc_local_json = json.loads(goxc_local_contents)
    goxc_local_json.pop("PrereleaseInfo", None)
    goxc_local_json.pop("BuildName", None)
    open(goxc_local_file, 'w').write(json.dumps(goxc_local_json, indent=4, separators=(',', ': ')) + "\n")

def do_release_snapshot(args):
    set_snapshot_build_info()

    subprocess.call(["goxc"])
    if args.push:
        subprocess.call(["goxc", "bintray"])

def do_release(args):
    try:
        remove_snapshot_build_info()

        subprocess.call(["goxc"])
        if args.push:
            subprocess.call(["goxc", "bintray"])
            subprocess.call(["goxc", "bump"])
            subprocess.call(["git", "add", ".goxc.json"])
            subprocess.call(["git", "commit", "-m", "'Bumping version'"])
    finally:
        set_snapshot_build_info()

parser = argparse.ArgumentParser()
parser.add_argument('mode', choices=['snapshot', 'release'])
parser.add_argument('-p', '--push', help="Push build to bintray", action="store_true")
args = parser.parse_args()

validate_goxc()
validate_goxc_local()

if args.mode == "snapshot":
    do_release_snapshot(args)
elif args.mode == "release":
    do_release(args)
