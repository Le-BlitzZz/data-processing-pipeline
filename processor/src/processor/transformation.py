import pandas as pd
import json


def transform(rows):
    df = pd.DataFrame(map(json.loads, rows))

    for col in df.columns:
        try:
            df[col] = pd.to_numeric(df[col])
        except ValueError:
            pass

    records = df.astype(str).to_dict(orient="records")

    return [json.dumps(record) for record in records]
