import trainer.config as config
import requests, time
import numpy as np
import pandas as pd

from sklearn.linear_model import LinearRegression
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import cross_validate
from sklearn.metrics import r2_score, root_mean_squared_error


def wait_for_api():
    while True:
        try:
            r = requests.get(config.API, params={"limit": 1})
            exp = int(r.headers["X-Expected-Count"])
            loaded = int(r.headers["X-Loaded-Count"])
            print(f"Expected: {exp}, Loaded: {loaded}")
            if loaded == exp:
                print("API is ready!")
                break
        except Exception as e:
            print(f"Error connecting to API: {e}")
        print("Waiting for API to be ready, sleeping 15 seconds...")
        time.sleep(15)


def prepare(df: pd.DataFrame):
    X = df.drop(columns=["ID", "CreatedAt", "UpdatedAt", "DeletedAt", "id", "split"])
    X = X.apply(pd.to_numeric)
    y = X.pop("price")
    return X, y


def main():
    wait_for_api()

    splits = ["train", "val", "test"]
    dfs = []
    for split in splits:
        print(f"Loading {split} data...")
        r = requests.get(config.API, params={"split": split})
        r.raise_for_status()
        data = r.json()
        print(f"  → Loaded {len(data)} records for {split}")
        dfs.append(pd.DataFrame(data))

    (X_train, y_train), (X_val, y_val), (X_test, y_test) = [prepare(df) for df in dfs]
    print("Data prepared: ", X_train.shape, X_val.shape, X_test.shape, y_train.shape, y_val.shape, y_test.shape)

    models = {
        "LinearRegression": LinearRegression(),
        "RandomForest": RandomForestRegressor(n_estimators=100, random_state=42),
    }

    # 2) Cross-validate on TRAIN spli
    cv_results = {}
    for name, model in models.items():
        print(f"Running 5-fold CV for {name}...")
        res = cross_validate(
            model,
            X_train,
            y_train,
            cv=5,
            scoring=["r2", "neg_root_mean_squared_error"],
            return_train_score=True,
        )
        cv_results[name] = {
            "train_r2": np.mean(res["train_r2"]),
            "valid_r2": np.mean(res["test_r2"]),
            "train_rmse": -np.mean(res["train_neg_root_mean_squared_error"]),
            "valid_rmse": -np.mean(res["test_neg_root_mean_squared_error"]),
        }
        print(
            f"  train squared R={cv_results[name]['train_r2']:.3f}, "
            f"valid squared R={cv_results[name]['valid_r2']:.3f}"
        )
        print(
            f"  train RMSE={cv_results[name]['train_rmse']:.0f}, "
            f"valid RMSE={cv_results[name]['valid_rmse']:.0f}"
        )

    best_name = max(cv_results, key=lambda n: cv_results[n]["valid_r2"])
    best_model = models[best_name]
    print(f"\nBest model by CV: {best_name}")

    best_model.fit(X_train, y_train)

    for split_name, (X, y) in [
        ("Validation", (X_val, y_val)),
        ("Test", (X_test, y_test)),
    ]:
        preds = best_model.predict(X)
        r2 = r2_score(y, preds)
        rmse = root_mean_squared_error(y, preds)
        print(f"{split_name} performance ({best_name}):")
        print(f"  squared R   = {r2:.3f}")
        print(f"  RMSE = {rmse:.0f}")


if __name__ == "__main__":
    main()
