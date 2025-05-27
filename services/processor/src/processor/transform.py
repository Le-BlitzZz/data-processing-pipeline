import sys
import pandas as pd
import joblib
import processor.config as config
from transformation.pipeline import PoiAndBoolPreprocessor

sys.modules["__main__"].PoiAndBoolPreprocessor = PoiAndBoolPreprocessor


def transform(rows):
    preprocessor = joblib.load(config.PREPROCESSOR_PATH)
    df = pd.DataFrame(rows)
    uuid = df.pop("id")
    price = df.pop("price")
    split = df.pop("split")

    for col in df.columns:
        try:
            df[col] = pd.to_numeric(df[col])
        except ValueError:
            pass

    df = pd.DataFrame(
        preprocessor.fit_transform(df),
        columns=preprocessor.named_steps["cols"].get_feature_names_out(),
    )

    df.insert(0, "id", uuid)
    df["price"] = price
    df["split"] = split

    records = df.astype(str).to_dict(orient="records")

    return records
