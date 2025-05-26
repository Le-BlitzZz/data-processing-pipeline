import pandas as pd


def transform(rows):
    df = pd.DataFrame(rows)

    for col in df.columns:
        try:
            df[col] = pd.to_numeric(df[col])
        except ValueError:
            pass

    records = df.astype(str).to_dict(orient="records")

    return records
