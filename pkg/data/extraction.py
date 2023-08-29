import pandas as pd
from zipfile import ZipFile
import os

ROOT = "/home/ubuntu/sla-operator" # change to match your workdir
CURRENT = os.path.join(ROOT, "pkg/data")
fp = os.path.join(CURRENT, "borg_traces_data.csv.zip")

with ZipFile(fp, 'r') as zip:
    zip.printdir()
    zip.extractall(CURRENT)

print(pd.read_csv(f'{CURRENT}/borg_traces_data.csv').head())